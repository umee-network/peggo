package coingecko

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

var logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel).With().Timestamp().Logger()

func TestQueryUSDPriceByCoinID(t *testing.T) {
	t.Run("ok", func(t *testing.T) {

		expected := `{"ethereum": {"usd": 4271.57}}`
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, expected)
		}))
		defer svr.Close()
		coingeckoFeed := NewCoingeckoPriceFeed(logger, &Config{BaseURL: svr.URL})

		price, err := coingeckoFeed.QueryUSDPriceByCoinID("ethereum")
		assert.NoError(t, err)
		assert.Equal(t, 4271.57, price)
	})

	t.Run("failed to parse response body", func(t *testing.T) {
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer svr.Close()
		coingeckoFeed := NewCoingeckoPriceFeed(logger, &Config{BaseURL: svr.URL})

		_, err := coingeckoFeed.QueryUSDPriceByCoinID("ethereum")
		assert.EqualError(t, err, "failed to parse response body from "+svr.URL+"/simple/price?ids=ethereum&vs_currencies=usd: EOF")
	})

	t.Run("price is zero", func(t *testing.T) {

		expected := `{"ethereum": {"usd": 0.0}}`
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, expected)
		}))
		defer svr.Close()
		coingeckoFeed := NewCoingeckoPriceFeed(logger, &Config{BaseURL: svr.URL})

		price, err := coingeckoFeed.QueryUSDPriceByCoinID("ethereum")
		assert.EqualError(t, err, "failed to get price for ethereum")
		assert.Equal(t, 0.0, price)
	})
}

func TestQueryTokenUSDPrice(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		expected := `{"0xdac17f958d2ee523a2206206994597c13d831ec7":{"usd":0.998233}}`
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, expected)
		}))
		defer svr.Close()
		coingeckoFeed := NewCoingeckoPriceFeed(logger, &Config{BaseURL: svr.URL})

		price, err := coingeckoFeed.QueryTokenUSDPrice(ethcmn.HexToAddress("0xdac17f958d2ee523a2206206994597c13d831ec7"))
		assert.NoError(t, err)
		assert.Equal(t, 0.998233, price)
	})

	t.Run("failed to parse response body", func(t *testing.T) {
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer svr.Close()
		coingeckoFeed := NewCoingeckoPriceFeed(logger, &Config{BaseURL: svr.URL})

		_, err := coingeckoFeed.QueryTokenUSDPrice(ethcmn.HexToAddress("0xdac17f958d2ee523a2206206994597c13d831ec7"))
		assert.EqualError(t, err, "failed to parse response body from "+svr.URL+"/simple/token_price/ethereum?contract_addresses=0xdac17f958d2ee523a2206206994597c13d831ec7&vs_currencies=usd: EOF")
	})

	t.Run("price is zero", func(t *testing.T) {
		expected := `{"0xdac17f958d2ee523a2206206994597c13d831ec7":{"usd":0.0}}`
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, expected)
		}))
		defer svr.Close()
		coingeckoFeed := NewCoingeckoPriceFeed(logger, &Config{BaseURL: svr.URL})

		price, err := coingeckoFeed.QueryTokenUSDPrice(ethcmn.HexToAddress("0xdac17f958d2ee523a2206206994597c13d831ec7"))
		assert.EqualError(t, err, "failed to get price for token 0xdAC17F958D2ee523a2206206994597C13D831ec7")
		assert.Equal(t, 0.0, price)
	})
}

func TestGetTokenSymbol(t *testing.T) {
	coinGecko := NewCoingeckoPriceFeed(logger, nil)
	symbol := "UMEE"
	umeeContractAddr := ethcmn.HexToAddress("0xc0a4Df35568F116C370E6a6A6022Ceb908eedDaC")
	coinGecko.setCoinSymbol(umeeContractAddr, symbol)

	requestedSymbol, err := coinGecko.GetTokenSymbol(umeeContractAddr)
	assert.Nil(t, err)
	assert.Equal(t, symbol, requestedSymbol)
}

func TestRequestCoinSymbol(t *testing.T) {
	t.Run("get umee symbol from contract", func(t *testing.T) {
		expected := `{"symbol": "umee"}`
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, expected)
		}))
		defer svr.Close()

		coinGecko := NewCoingeckoPriceFeed(logger, &Config{BaseURL: svr.URL})
		symbol, err := coinGecko.requestCoinSymbol(ethcmn.HexToAddress("0xc0a4Df35568F116C370E6a6A6022Ceb908eedDaC"))
		assert.Nil(t, err)
		assert.Equal(t, symbol, "UMEE")
	})

	t.Run("symbol not found", func(t *testing.T) {
		expected := `{"error": "Could not find coin with the given id"}`
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, expected)
		}))
		defer svr.Close()

		coinGecko := NewCoingeckoPriceFeed(logger, &Config{BaseURL: svr.URL})
		symbol, err := coinGecko.requestCoinSymbol(ethcmn.HexToAddress("----"))
		assert.EqualError(t, err, "Could not find coin with the given id")
		assert.Equal(t, symbol, "")
	})
}

func TestGetRequestCoinSymbolURL(t *testing.T) {
	coinGecko := NewCoingeckoPriceFeed(logger, nil)
	url, err := coinGecko.getRequestCoinSymbolURL(ethcmn.HexToAddress("0xc0a4Df35568F116C370E6a6A6022Ceb908eedDaC"))
	assert.Nil(t, err)
	assert.Equal(t, url.String(), "https://api.coingecko.com/api/v3/coins/ethereum/contract/0xc0a4Df35568F116C370E6a6A6022Ceb908eedDaC")
}

func TestCheckCoingeckoConfig(t *testing.T) {
	assert.NotNil(t, checkCoingeckoConfig(nil))
	assert.NotNil(t, checkCoingeckoConfig(&Config{BaseURL: ""}))
}
