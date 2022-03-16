package coingecko

import (
	"encoding/json"
	"strings"

	"net/http"
	"net/url"
	"path"
	"time"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const (
	maxRespTime        = 15 * time.Second
	maxRespHeadersTime = 15 * time.Second
	EthereumCoinID     = "ethereum"
)

var zeroPrice = float64(0)

type CoinGecko struct {
	client *http.Client
	config *Config

	coinsSymbol map[ethcmn.Address]string // contract addr => token symbol

	logger zerolog.Logger
}

// Config wraps config variable to get information at CoinGecko
type Config struct {
	BaseURL string
}

// CoinInfo wraps the coin information receveid from an contract address.
// https://api.coingecko.com/api/v3/coins/ethereum/contract/${CONTRACT_ADDR}
type CoinInfo struct {
	Symbol string `json:"symbol"`
	Error  string `json:"error"`
}

type priceResponse map[string]struct {
	USD float64 `json:"usd"`
}

// NewCoingeckoPriceFeed returns price puller for given symbol. The price will be pulled
// from endpoint and divided by scaleFactor. Symbol name (if reported by endpoint) must match.
func NewCoingeckoPriceFeed(logger zerolog.Logger, endpointConfig *Config) *CoinGecko {
	return &CoinGecko{
		client: &http.Client{
			Transport: &http.Transport{
				ResponseHeaderTimeout: maxRespHeadersTime,
			},
			Timeout: maxRespTime,
		},
		config:      checkCoingeckoConfig(endpointConfig),
		coinsSymbol: make(map[ethcmn.Address]string),
		logger:      logger.With().Str("oracle", "coingecko").Logger(),
	}
}

func urlJoin(baseURL string, segments ...string) (*url.URL, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(append([]string{u.Path}, segments...)...)
	return u, nil
}

// GetTokenSymbol returns the token symbol checked by CoinGecko API.
func (cp *CoinGecko) GetTokenSymbol(erc20Contract ethcmn.Address) (string, error) {
	symbol, ok := cp.coinsSymbol[erc20Contract]
	if !ok {
		symbol, err := cp.requestCoinSymbol(erc20Contract)
		if err != nil {
			return "", err
		}
		cp.setCoinSymbol(erc20Contract, symbol)

		return symbol, nil
	}

	return symbol, nil
}

func (cp *CoinGecko) setCoinSymbol(erc20Contract ethcmn.Address, symbol string) {
	cp.coinsSymbol[erc20Contract] = symbol
}

func (cp *CoinGecko) getRequestCoinSymbolURL(erc20Contract ethcmn.Address) (*url.URL, error) {
	return urlJoin(cp.config.BaseURL, "coins", EthereumCoinID, "contract", erc20Contract.Hex())
}

func (cp *CoinGecko) requestCoinSymbol(erc20Contract ethcmn.Address) (string, error) {
	u, err := cp.getRequestCoinSymbolURL(erc20Contract)
	if err != nil {
		cp.logger.Fatal().Err(err).Msg("failed to parse coin info URL")
	}

	reqURL := u.String()
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		cp.logger.Fatal().Err(err).Msg("failed to create HTTP request of coin info")
	}

	resp, err := cp.client.Do(req)
	if err != nil {
		return "", errors.Wrapf(err, "failed to fetch coin info from %s", reqURL)
	}
	defer resp.Body.Close()

	var coinInfo CoinInfo
	if err := json.NewDecoder(resp.Body).Decode(&coinInfo); err != nil {
		return "", errors.Wrapf(err, "failed to parse response body from %s", reqURL)
	}

	if len(coinInfo.Error) > 0 {
		return "", errors.New(coinInfo.Error)
	}

	if len(coinInfo.Symbol) == 0 {
		return "", errors.New("Fail to get coin info for contract: " + erc20Contract.Hex())
	}

	return strings.ToUpper(coinInfo.Symbol), nil
}

func (cp *CoinGecko) QueryUSDPriceByCoinID(coinID string) (float64, error) {
	u, err := urlJoin(cp.config.BaseURL, "simple", "price")
	if err != nil {
		cp.logger.Fatal().Err(err).Msg("failed to parse URL")
	}

	q := make(url.Values)

	q.Set("ids", coinID)
	q.Set("vs_currencies", "usd")
	u.RawQuery = q.Encode()

	reqURL := u.String()
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		cp.logger.Fatal().Err(err).Msg("failed to create HTTP request")
	}

	resp, err := cp.client.Do(req)
	if err != nil {
		err = errors.Wrapf(err, "failed to fetch price from %s", reqURL)
		return zeroPrice, err
	}

	defer resp.Body.Close()

	var respBody priceResponse

	err = json.NewDecoder(resp.Body).Decode(&respBody)

	if err != nil {
		return zeroPrice, errors.Wrapf(err, "failed to parse response body from %s", reqURL)
	}

	price := respBody[coinID].USD

	if price == zeroPrice {
		return zeroPrice, errors.Errorf("failed to get price for %s", coinID)
	}

	return price, nil
}

func (cp *CoinGecko) QueryTokenUSDPrice(erc20Contract ethcmn.Address) (float64, error) {
	// If the token is one of the deployed by the Gravity contract, use the
	// stored coin ID to look up the price.
	if coinID, ok := bridgeTokensCoinIDs[erc20Contract.Hex()]; ok {
		return cp.QueryUSDPriceByCoinID(coinID)
	}

	u, err := urlJoin(cp.config.BaseURL, "simple", "token_price", EthereumCoinID)
	if err != nil {
		cp.logger.Fatal().Err(err).Msg("failed to parse URL")
	}

	q := make(url.Values)

	q.Set("contract_addresses", strings.ToLower(erc20Contract.String()))
	q.Set("vs_currencies", "usd")
	u.RawQuery = q.Encode()

	reqURL := u.String()
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		cp.logger.Fatal().Err(err).Msg("failed to create HTTP request")
	}

	resp, err := cp.client.Do(req)
	if err != nil {
		err = errors.Wrapf(err, "failed to fetch price from %s", reqURL)
		return zeroPrice, err
	}

	defer resp.Body.Close()

	var respBody priceResponse

	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return zeroPrice, errors.Wrapf(err, "failed to parse response body from %s", reqURL)
	}

	price := respBody[strings.ToLower(erc20Contract.String())].USD

	if price == zeroPrice {
		return zeroPrice, errors.Errorf("failed to get price for token %s", erc20Contract.Hex())
	}

	return price, nil
}

func checkCoingeckoConfig(cfg *Config) *Config {
	if cfg == nil {
		cfg = &Config{}
	}

	if len(cfg.BaseURL) == 0 {
		cfg.BaseURL = "https://api.coingecko.com/api/v3"
	}

	return cfg
}
