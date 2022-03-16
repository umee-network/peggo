package relayer

import (
	"github.com/umee-network/peggo/orchestrator/coingecko"
	"github.com/umee-network/peggo/orchestrator/oracle"
)

func SetPriceFeeder(coinGecko *coingecko.CoinGecko) func(GravityRelayer) {
	return func(s GravityRelayer) { s.SetPriceFeeder(coinGecko) }
}

func (s *gravityRelayer) SetPriceFeeder(coinGecko *coingecko.CoinGecko) {
	s.coinGecko = coinGecko
}

// SetOracle sets a new oracle to the Gravity Relayer.
func SetOracle(o oracle.PriceFeeder) func(GravityRelayer) {
	return func(s GravityRelayer) { s.SetOracle(o) }
}

// SetOracle sets a new oracle to the Gravity Relayer.
func (s *gravityRelayer) SetOracle(o oracle.PriceFeeder) {
	s.oracle = o
}
