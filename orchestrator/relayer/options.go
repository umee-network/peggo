package relayer

import (
	"github.com/umee-network/peggo/orchestrator/coingecko"
	"github.com/umee-network/peggo/orchestrator/oracle"
)

func SetCoinGecko(coinGecko *coingecko.CoinGecko) func(GravityRelayer) {
	return func(s GravityRelayer) { s.SetCoinGecko(coinGecko) }
}

func (s *gravityRelayer) SetCoinGecko(coinGecko *coingecko.CoinGecko) {
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
