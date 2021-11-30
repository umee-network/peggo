package relayer

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/umee-network/peggo/orchestrator/cosmos"
	"github.com/umee-network/peggo/orchestrator/ethereum/peggy"
	"github.com/umee-network/peggo/orchestrator/ethereum/provider"
	"github.com/umee-network/umee/x/peggy/types"
)

type PeggyRelayer interface {
	Start(ctx context.Context) error

	FindLatestValset(ctx context.Context) (*types.Valset, error)
	RelayBatches(ctx context.Context) error
	RelayValsets(ctx context.Context) error
}

type peggyRelayer struct {
	logger             zerolog.Logger
	cosmosQueryClient  cosmos.PeggyQueryClient
	peggyContract      peggy.Contract
	ethProvider        provider.EVMProvider
	valsetRelayEnabled bool
	batchRelayEnabled  bool
	loopDuration       time.Duration
<<<<<<< HEAD
=======
	priceFeeder        *coingecko.PriceFeed
	pendingTxWait      time.Duration

	// store locally the last tx this validator made to avoid sending duplicates
	// or invalid txs
	lastSentBatchNonce uint64
>>>>>>> 5c56265 (feat: tweak loops durations (#60))
}

func NewPeggyRelayer(
	logger zerolog.Logger,
	cosmosQueryClient cosmos.PeggyQueryClient,
	peggyContract peggy.Contract,
	valsetRelayEnabled bool,
	batchRelayEnabled bool,
	loopDuration time.Duration,
<<<<<<< HEAD
=======
	pendingTxWait time.Duration,
	options ...func(PeggyRelayer),
>>>>>>> 5c56265 (feat: tweak loops durations (#60))
) PeggyRelayer {
	return &peggyRelayer{
		logger:             logger.With().Str("module", "peggy_relayer").Logger(),
		cosmosQueryClient:  cosmosQueryClient,
		peggyContract:      peggyContract,
		ethProvider:        peggyContract.Provider(),
		valsetRelayEnabled: valsetRelayEnabled,
		batchRelayEnabled:  batchRelayEnabled,
		loopDuration:       loopDuration,
<<<<<<< HEAD
=======
		pendingTxWait:      pendingTxWait,
>>>>>>> 5c56265 (feat: tweak loops durations (#60))
	}
}
