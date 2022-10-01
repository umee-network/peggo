package oracle

import (
	"context"
	"fmt"
<<<<<<< HEAD
	"sync"
	"time"

	sdkmath "cosmossdk.io/math"
=======
	"strings"
	"sync"
	"time"

>>>>>>> dbba311d3ef1e6ec73aa7b4d5366620ef63ad4e0
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

<<<<<<< HEAD
	pforacle "github.com/umee-network/umee/price-feeder/oracle"
	pfprovider "github.com/umee-network/umee/price-feeder/oracle/provider"
	pftypes "github.com/umee-network/umee/price-feeder/oracle/types"
	pfsync "github.com/umee-network/umee/price-feeder/pkg/sync"
	umeeparams "github.com/umee-network/umee/v3/app/params"
)

const (
	// tickerTimeout is the minimum timeout between each oracle loop.
	tickerTimeout = 1000 * time.Millisecond
	// availablePairsReload is the amount of time to reload the providers available pairs.
	availablePairsReload = 24 * time.Hour
	// SymbolETH refers to the ethereum symbol.
	SymbolETH = "ETH"
=======
	ummedpforacle "github.com/umee-network/umee/price-feeder/oracle"
	umeedpfprovider "github.com/umee-network/umee/price-feeder/oracle/provider"
	umeedpftypes "github.com/umee-network/umee/price-feeder/oracle/types"
	ummedpfsync "github.com/umee-network/umee/price-feeder/pkg/sync"
)

// We define tickerTimeout as the minimum timeout between each oracle loop.
const (
	tickerTimeout        = 1000 * time.Millisecond
	availablePairsReload = 24 * time.Hour
	BaseSymbolETH        = "ETH"
)

var (
	// deviationThreshold defines how many 𝜎 a provider can be away from the mean
	// without being considered faulty.
	deviationThreshold = sdk.MustNewDecFromStr("2")
>>>>>>> dbba311d3ef1e6ec73aa7b4d5366620ef63ad4e0
)

// Oracle implements the core component responsible for fetching exchange rates
// for a given set of currency pairs and determining the correct exchange rates.
type Oracle struct {
	logger zerolog.Logger
<<<<<<< HEAD
	closer *pfsync.Closer

	mtx                   sync.RWMutex
	providers             map[pfprovider.Name]*Provider // providerName => Provider
	prices                map[string]sdk.Dec            // baseSymbol => price ex.: UMEE, ETH => sdk.Dec
	subscribedBaseSymbols map[string]struct{}           // baseSymbol => nothing
	// this field could be calculated each time by looping providers.subscribedPairs
	// but the time to process is not worth the amount of memory
	providerSubscribedPairs map[pfprovider.Name][]pftypes.CurrencyPair // providerName => []CurrencyPair
=======
	closer *ummedpfsync.Closer

	mtx                   sync.RWMutex
	providers             map[string]*Provider // providerName => Provider
	prices                map[string]sdk.Dec   // baseSymbol => price ex.: UMEE, ETH => sdk.Dec
	subscribedBaseSymbols map[string]struct{}  // baseSymbol => nothing
>>>>>>> dbba311d3ef1e6ec73aa7b4d5366620ef63ad4e0
}

// Provider wraps the umee provider interface.
type Provider struct {
<<<<<<< HEAD
	pfprovider.Provider
	availablePairs  map[string]struct{}             // Symbol => nothing
	subscribedPairs map[string]pftypes.CurrencyPair // Symbol => currencyPair
}

func New(ctx context.Context, logger zerolog.Logger, providersName []pfprovider.Name) (*Oracle, error) {
	providers := map[pfprovider.Name]*Provider{}

	for _, providerName := range providersName {
		provider, err := pforacle.NewProvider(
			ctx,
			providerName,
			logger,
			pfprovider.Endpoint{},
			pftypes.CurrencyPair{},
		)
=======
	umeedpfprovider.Provider
	availablePairs  map[string]struct{}                  // Symbol => nothing
	subscribedPairs map[string]umeedpftypes.CurrencyPair // Symbol => currencyPair
}

func New(ctx context.Context, logger zerolog.Logger, providersName []string) (*Oracle, error) {
	providers := map[string]*Provider{}

	for _, providerName := range providersName {
		provider, err := ummedpforacle.NewProvider(ctx, providerName, logger, umeedpftypes.CurrencyPair{})
>>>>>>> dbba311d3ef1e6ec73aa7b4d5366620ef63ad4e0
		if err != nil {
			return nil, err
		}

		providers[providerName] = &Provider{
			Provider:        provider,
			availablePairs:  map[string]struct{}{},
<<<<<<< HEAD
			subscribedPairs: map[string]pftypes.CurrencyPair{},
		}
	}

	o := &Oracle{
		logger:                  logger.With().Str("module", "oracle").Logger(),
		closer:                  pfsync.NewCloser(),
		providers:               providers,
		subscribedBaseSymbols:   map[string]struct{}{},
		providerSubscribedPairs: map[pfprovider.Name][]pftypes.CurrencyPair{},
	}
	o.loadAvailablePairs()
	o.mtx.Lock()
	defer o.mtx.Unlock()
	if err := o.subscribeProviders([]pftypes.CurrencyPair{
		{Base: symbolUSDT, Quote: symbolUSD},
		{Base: symbolDAI, Quote: symbolUSD},
	}); err != nil {
		return nil, err
	}
	go o.start(ctx)

	return o, nil
=======
			subscribedPairs: map[string]umeedpftypes.CurrencyPair{},
		}
	}

	oracle := &Oracle{
		logger:                logger.With().Str("module", "oracle").Logger(),
		closer:                ummedpfsync.NewCloser(),
		providers:             providers,
		subscribedBaseSymbols: map[string]struct{}{},
	}
	oracle.loadAvailablePairs()
	go oracle.start(ctx)

	return oracle, nil
>>>>>>> dbba311d3ef1e6ec73aa7b4d5366620ef63ad4e0
}

// GetPrices returns the price for the provided base symbols.
func (o *Oracle) GetPrices(baseSymbols ...string) (map[string]sdk.Dec, error) {
	o.mtx.RLock()
	defer o.mtx.RUnlock()

	// Creates a new array for the prices in the oracle.
	prices := make(map[string]sdk.Dec, len(baseSymbols))

	for _, baseSymbol := range baseSymbols {
		price, ok := o.prices[baseSymbol]
		if !ok {
			return nil, fmt.Errorf("error getting price for %s", baseSymbol)
		}
		prices[baseSymbol] = price
	}

	return prices, nil
}

// GetPrice returns the price based on the symbol ex.: UMEE, ETH.
func (o *Oracle) GetPrice(baseSymbol string) (sdk.Dec, error) {
	o.mtx.RLock()
	defer o.mtx.RUnlock()

	price, ok := o.prices[baseSymbol]
	if !ok {
		return sdk.Dec{}, fmt.Errorf("error getting price for %s", baseSymbol)
	}

	return price, nil
}

// SubscribeSymbols attempts to subscribe the symbols in all the providers.
// baseSymbols is the base to be subscribed ex.: ["UMEE", "ATOM"].
func (o *Oracle) SubscribeSymbols(baseSymbols ...string) error {
	o.mtx.Lock()
	defer o.mtx.Unlock()

	for _, baseSymbol := range baseSymbols {
		_, ok := o.subscribedBaseSymbols[baseSymbol]
		if ok {
			// pair already subscribed
			continue
		}

		currencyPairs := GetStablecoinsCurrencyPair(baseSymbol)
		if err := o.subscribeProviders(currencyPairs); err != nil {
			return err
		}

		o.logger.Debug().
			Str("token_symbol", baseSymbol).
			Msg("New symbol subscribed")

		o.subscribedBaseSymbols[baseSymbol] = struct{}{}
	}

	return nil
}

<<<<<<< HEAD
func (o *Oracle) subscribeProviders(currencyPairs []pftypes.CurrencyPair) error {
	for providerName, provider := range o.providers {
		var pairsToSubscribe []pftypes.CurrencyPair
=======
func (o *Oracle) subscribeProviders(currencyPairs []umeedpftypes.CurrencyPair) error {
	for providerName, provider := range o.providers {
		var pairsToSubscribe []umeedpftypes.CurrencyPair
>>>>>>> dbba311d3ef1e6ec73aa7b4d5366620ef63ad4e0

		for _, currencyPair := range currencyPairs {
			symbol := currencyPair.String()

			_, ok := provider.subscribedPairs[symbol]
			if ok {
				// currency pair already subscribed
				continue
			}

			_, availablePair := provider.availablePairs[symbol]
			if !availablePair {
<<<<<<< HEAD
				o.logger.Debug().Str("provider_name", string(providerName)).Str("symbol", symbol).Msg("symbol is not available")
=======
				o.logger.Debug().Str("provider_name", providerName).Str("symbol", symbol).Msg("symbol is not available")
>>>>>>> dbba311d3ef1e6ec73aa7b4d5366620ef63ad4e0
				continue
			}

			pairsToSubscribe = append(pairsToSubscribe, currencyPair)
		}

		if len(pairsToSubscribe) == 0 {
<<<<<<< HEAD
			o.logger.Debug().Str("provider_name", string(providerName)).
=======
			o.logger.Debug().Str("provider_name", providerName).
>>>>>>> dbba311d3ef1e6ec73aa7b4d5366620ef63ad4e0
				Msgf("No pairs to subscribe, received pairs to try: %+v", currencyPairs)
			continue
		}

		if err := provider.SubscribeCurrencyPairs(pairsToSubscribe...); err != nil {
<<<<<<< HEAD
			o.logger.Err(err).Str("provider_name", string(providerName)).Msg("subscribing to new currency pairs")
=======
			o.logger.Err(err).Str("provider_name", providerName).Msg("subscribing to new currency pairs")
>>>>>>> dbba311d3ef1e6ec73aa7b4d5366620ef63ad4e0
			return err
		}

		for _, pair := range pairsToSubscribe {
			provider.subscribedPairs[pair.String()] = pair
<<<<<<< HEAD
			o.providerSubscribedPairs[providerName] = append(o.providerSubscribedPairs[providerName], pair)

			o.logger.Debug().Str("provider_name", string(providerName)).
=======

			o.logger.Debug().Str("provider_name", providerName).
>>>>>>> dbba311d3ef1e6ec73aa7b4d5366620ef63ad4e0
				Str("pair_symbol", pair.String()).
				Msg("Subscribed new pair")
		}

<<<<<<< HEAD
		o.logger.Info().Str("provider_name", string(providerName)).
=======
		o.logger.Info().Str("provider_name", providerName).
>>>>>>> dbba311d3ef1e6ec73aa7b4d5366620ef63ad4e0
			Int("currency_pairs_length", len(pairsToSubscribe)).
			Msgf("Subscribed pairs %+v", pairsToSubscribe)
	}

	return nil
}

// Stop stops the oracle process and waits for it to gracefully exit.
func (o *Oracle) Stop() {
	o.closer.Close()
	<-o.closer.Done()
}

// start starts the oracle process in a blocking fashion.
func (o *Oracle) start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			o.closer.Close()

		case <-time.After(tickerTimeout):
			if err := o.tick(); err != nil {
				o.logger.Err(err).Msg("oracle tick failed")
			}

		case <-time.After(availablePairsReload):
			o.loadAvailablePairs()
		}
	}
}

// loadAvailablePairs loads all the available pairs from providers.
func (o *Oracle) loadAvailablePairs() {
	o.mtx.Lock()
	defer o.mtx.Unlock()

	for providerName, provider := range o.providers {
		availablePairs, err := provider.GetAvailablePairs()
		if err != nil {
<<<<<<< HEAD
			o.logger.Debug().Err(err).Str(
				"provider_name",
				string(providerName),
			).Msg("Error getting available pairs for provider")
=======
			o.logger.Debug().Err(err).Str("provider_name", providerName).Msg("Error getting available pairs for provider")
>>>>>>> dbba311d3ef1e6ec73aa7b4d5366620ef63ad4e0
			continue
		}
		if len(availablePairs) == 0 {
			continue
		}
		provider.availablePairs = availablePairs
	}
}

// setPrices retrieves all the prices and candles from our set of providers as
// determined in the config. If candles are available, uses TVWAP in order
// to determine prices. If candles are not available, uses the most recent prices
// with VWAP. Warns the the user of any missing prices, and filters out any faulty
// providers which do not report prices or candles within 2𝜎 of the others.
// code originally from https://github.com/umee-network/umee/blob/2a69b56ae1c6098cb2d23ef8384f5acf28f76d35/price-feeder/oracle/oracle.go#L166-L167
func (o *Oracle) setPrices() error {
	g := new(errgroup.Group)
	mtx := new(sync.Mutex)
<<<<<<< HEAD
	providerPrices := make(pfprovider.AggregatedProviderPrices)
	providerCandles := make(pfprovider.AggregatedProviderCandles)
=======
	providerPrices := make(umeedpfprovider.AggregatedProviderPrices)
	providerCandles := make(umeedpfprovider.AggregatedProviderCandles)
>>>>>>> dbba311d3ef1e6ec73aa7b4d5366620ef63ad4e0

	for providerName, provider := range o.providers {
		providerName := providerName
		provider := provider
<<<<<<< HEAD
		subscribedPrices := o.providerSubscribedPairs[providerName]
=======
		subscribedPrices := umeedpftypes.MapPairsToSlice(provider.subscribedPairs)
>>>>>>> dbba311d3ef1e6ec73aa7b4d5366620ef63ad4e0

		g.Go(func() error {
			var (
				tickerErr error
				candleErr error
			)

			prices, tickerErr := provider.GetTickerPrices(subscribedPrices...)
			candles, candleErr := provider.GetCandlePrices(subscribedPrices...)

			if tickerErr != nil && candleErr != nil {
				// only generates error if ticker and candle generate errors
<<<<<<< HEAD
				o.logger.Debug().Msgf("provider: %s ticker error: %+v\ncandle error: %+v", providerName, tickerErr, candleErr)
				return nil
=======
				return fmt.Errorf("ticker error: %+v\ncandle error: %+v", tickerErr, candleErr)
>>>>>>> dbba311d3ef1e6ec73aa7b4d5366620ef63ad4e0
			}

			// flatten and collect prices based on the base currency per provider
			//
			// e.g.: {ProviderKraken: {"ATOM": <price, volume>, ...}}
			mtx.Lock()
			for _, pair := range subscribedPrices {
<<<<<<< HEAD
				pforacle.SetProviderTickerPricesAndCandles(
					providerName,
					providerPrices,
					providerCandles,
					prices,
					candles,
					pair,
				)
=======
				if _, ok := providerPrices[providerName]; !ok {
					providerPrices[providerName] = make(map[string]umeedpfprovider.TickerPrice)
				}
				if _, ok := providerCandles[providerName]; !ok {
					providerCandles[providerName] = make(map[string][]umeedpfprovider.CandlePrice)
				}

				tp, pricesOk := prices[pair.String()]
				if pricesOk {
					providerPrices[providerName][pair.Base] = tp
				}

				cp, candlesOk := candles[pair.String()]
				if candlesOk {
					providerCandles[providerName][pair.Base] = cp
				}
>>>>>>> dbba311d3ef1e6ec73aa7b4d5366620ef63ad4e0
			}

			mtx.Unlock()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		o.logger.Debug().Err(err).Msg("failed to get ticker prices from provider")
	}

<<<<<<< HEAD
	deviationTreshold := sdk.NewDecFromIntWithPrec(sdkmath.NewInt(15), 1)

	computedPrices, err := pforacle.GetComputedPrices(
		o.logger,
		providerCandles,
		providerPrices,
		o.providerSubscribedPairs,
		map[string]sdk.Dec{
			SymbolETH:            deviationTreshold,
			umeeparams.BondDenom: deviationTreshold,
		}, // uses default deviation for other bases
	)
=======
	filteredCandles, err := o.filterCandleDeviations(providerCandles)
>>>>>>> dbba311d3ef1e6ec73aa7b4d5366620ef63ad4e0
	if err != nil {
		return err
	}

<<<<<<< HEAD
	o.prices = computedPrices
	return nil
}

=======
	// attempt to use candles for TVWAP calculations
	tvwapPrices, err := ummedpforacle.ComputeTVWAP(filteredCandles)
	if err != nil {
		return err
	}

	// If TVWAP candles are not available or were filtered out due to staleness,
	// use most recent prices & VWAP instead.
	if len(tvwapPrices) == 0 {
		filteredProviderPrices, err := o.filterTickerDeviations(providerPrices)
		if err != nil {
			return err
		}

		vwapPrices, err := ummedpforacle.ComputeVWAP(filteredProviderPrices)
		if err != nil {
			return err
		}

		// warn the user of any missing prices
		reportedPrices := make(map[string]struct{})
		for _, providers := range filteredProviderPrices {
			for base := range providers {
				if _, ok := reportedPrices[base]; !ok {
					reportedPrices[base] = struct{}{}
				}
			}
		}

		o.prices = vwapPrices
	} else {
		// warn the user of any missing candles
		reportedCandles := make(map[string]struct{})
		for _, providers := range filteredCandles {
			for base := range providers {
				if _, ok := reportedCandles[base]; !ok {
					reportedCandles[base] = struct{}{}
				}
			}
		}

		o.prices = tvwapPrices
	}

	return nil
}

// filterCandleDeviations finds the standard deviations of the tvwaps of
// all assets, and filters out any providers that are not within 2𝜎 of the mean.
// code originally from https://github.com/umee-network/umee/blob/2a69b56ae1c6098cb2d23ef8384f5acf28f76d35/price-feeder/oracle/oracle.go#L458-L459
func (o *Oracle) filterCandleDeviations(
	candles umeedpfprovider.AggregatedProviderCandles,
) (umeedpfprovider.AggregatedProviderCandles, error) {
	var (
		filteredCandles = make(umeedpfprovider.AggregatedProviderCandles)
		tvwaps          = make(map[string]map[string]sdk.Dec)
	)

	for providerName, priceCandles := range candles {
		candlePrices := make(umeedpfprovider.AggregatedProviderCandles)

		for base, cp := range priceCandles {
			if _, ok := candlePrices[providerName]; !ok {
				candlePrices[providerName] = make(map[string][]umeedpfprovider.CandlePrice)
			}

			candlePrices[providerName][base] = cp
		}

		tvwap, err := ummedpforacle.ComputeTVWAP(candlePrices)
		if err != nil {
			return nil, err
		}

		for base, asset := range tvwap {
			if _, ok := tvwaps[providerName]; !ok {
				tvwaps[providerName] = make(map[string]sdk.Dec)
			}

			tvwaps[providerName][base] = asset
		}
	}

	deviations, means, err := ummedpforacle.StandardDeviation(tvwaps)
	if err != nil {
		return nil, err
	}

	// accept any tvwaps that are within 2𝜎, or for which we couldn't get 𝜎
	for providerName, priceMap := range tvwaps {
		for base, price := range priceMap {
			if _, ok := deviations[base]; !ok ||
				(price.GTE(means[base].Sub(deviations[base].Mul(deviationThreshold))) &&
					price.LTE(means[base].Add(deviations[base].Mul(deviationThreshold)))) {
				if _, ok := filteredCandles[providerName]; !ok {
					filteredCandles[providerName] = make(map[string][]umeedpfprovider.CandlePrice)
				}

				filteredCandles[providerName][base] = candles[providerName][base]
			} else {
				o.logger.Warn().
					Str("base", base).
					Str("provider", providerName).
					Str("price", price.String()).
					Msg("provider deviating from other candles")
			}
		}
	}

	return filteredCandles, nil
}

// filterTickerDeviations finds the standard deviations of the prices of
// all assets, and filters out any providers that are not within 2𝜎 of the mean.
// code originally from https://github.com/umee-network/umee/blob/2a69b56ae1c6098cb2d23ef8384f5acf28f76d35/price-feeder/oracle/oracle.go#L409-L410
func (o *Oracle) filterTickerDeviations(
	prices umeedpfprovider.AggregatedProviderPrices,
) (umeedpfprovider.AggregatedProviderPrices, error) {
	var (
		filteredPrices = make(umeedpfprovider.AggregatedProviderPrices)
		priceMap       = make(map[string]map[string]sdk.Dec)
	)

	for providerName, priceTickers := range prices {
		if _, ok := priceMap[providerName]; !ok {
			priceMap[providerName] = make(map[string]sdk.Dec)
		}
		for base, tp := range priceTickers {
			priceMap[providerName][base] = tp.Price
		}
	}

	deviations, means, err := ummedpforacle.StandardDeviation(priceMap)
	if err != nil {
		return nil, err
	}

	// accept any prices that are within 2𝜎, or for which we couldn't get 𝜎
	for providerName, priceTickers := range prices {
		for base, tp := range priceTickers {
			if _, ok := deviations[base]; !ok ||
				(tp.Price.GTE(means[base].Sub(deviations[base].Mul(deviationThreshold))) &&
					tp.Price.LTE(means[base].Add(deviations[base].Mul(deviationThreshold)))) {
				if _, ok := filteredPrices[providerName]; !ok {
					filteredPrices[providerName] = make(map[string]umeedpfprovider.TickerPrice)
				}

				filteredPrices[providerName][base] = tp
			} else {
				o.logger.Warn().
					Str("base", base).
					Str("provider", providerName).
					Str("price", tp.Price.String()).
					Msg("provider deviating from other prices")
			}
		}
	}

	return filteredPrices, nil
}

>>>>>>> dbba311d3ef1e6ec73aa7b4d5366620ef63ad4e0
func (o *Oracle) tick() error {
	if err := o.setPrices(); err != nil {
		return err
	}

	return nil
}
<<<<<<< HEAD
=======

// GetStablecoinsCurrencyPair return the currency pair of that symbol quoted by some
// stablecoins.
func GetStablecoinsCurrencyPair(baseSymbol string) []umeedpftypes.CurrencyPair {
	quotes := []string{"USD", "USDT", "UST"}
	currencyPairs := make([]umeedpftypes.CurrencyPair, len(quotes))

	for i, quote := range quotes {
		currencyPairs[i] = umeedpftypes.CurrencyPair{
			Base:  strings.ToUpper(baseSymbol),
			Quote: quote,
		}
	}

	return currencyPairs
}
>>>>>>> dbba311d3ef1e6ec73aa7b4d5366620ef63ad4e0
