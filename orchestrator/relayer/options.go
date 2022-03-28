package relayer

import (
	"github.com/umee-network/peggo/orchestrator/oracle"
)

func SetSymbolRetriever(coinGecko SymbolRetriever) func(GravityRelayer) {
	return func(s GravityRelayer) { s.SetSymbolRetriever(coinGecko) }
}

func (s *gravityRelayer) SetSymbolRetriever(symbolRetriever SymbolRetriever) {
	s.symbolRetriever = symbolRetriever
}

// SetOracle sets a new oracle to the Gravity Relayer.
func SetOracle(o oracle.PriceFeeder) func(GravityRelayer) {
	return func(s GravityRelayer) { s.SetOracle(o) }
}

// SetOracle sets a new oracle to the Gravity Relayer.
func (s *gravityRelayer) SetOracle(o oracle.PriceFeeder) {
	s.oracle = o
}
