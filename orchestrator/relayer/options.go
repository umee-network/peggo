package relayer

import "github.com/umee-network/peggo/orchestrator/coingecko"

// SetMinBatchFee sets the (optional) minimum batch fee denominated in USD.
func SetMinBatchFee(minFee float64) func(PeggyRelayer) {
	return func(s PeggyRelayer) { s.SetMinBatchFee(minFee) }
}

func (s *peggyRelayer) SetMinBatchFee(minFee float64) {
	s.minBatchFeeUSD = minFee
}

// SetPriceFeeder sets the (optional) price feeder used when performing profitable
// batch calculations. Note, this should be supplied only when the min batch
// fee is non-zero.
func SetPriceFeeder(pf *coingecko.PriceFeed) func(PeggyRelayer) {
	return func(s PeggyRelayer) { s.SetPriceFeeder(pf) }
}

func (s *peggyRelayer) SetPriceFeeder(pf *coingecko.PriceFeed) {
	s.priceFeeder = pf
}
