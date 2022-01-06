package orchestrator

import (
	"context"
	"sync"
	"time"

	gravitytypes "github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravity/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	sidechain "github.com/umee-network/peggo/orchestrator/cosmos"
	peggy "github.com/umee-network/peggo/orchestrator/ethereum/gravity"
	"github.com/umee-network/peggo/orchestrator/ethereum/keystore"
	"github.com/umee-network/peggo/orchestrator/ethereum/provider"
	"github.com/umee-network/peggo/orchestrator/relayer"
)

type PeggyOrchestrator interface {
	Start(ctx context.Context) error
	CheckForEvents(ctx context.Context, startingBlock, ethBlockConfirmationDelay uint64) (currentBlock uint64, err error)
	GetLastCheckedBlock(ctx context.Context, ethBlockConfirmationDelay uint64) (uint64, error)
	EthOracleMainLoop(ctx context.Context) error
	EthSignerMainLoop(ctx context.Context) error
	BatchRequesterLoop(ctx context.Context) error
	RelayerMainLoop(ctx context.Context) error
}

type peggyOrchestrator struct {
	logger                     zerolog.Logger
	cosmosQueryClient          gravitytypes.QueryClient
	gravityBroadcastClient     sidechain.PeggyBroadcastClient
	gravityContract            peggy.Contract
	ethProvider                provider.EVMProvider
	ethFrom                    ethcmn.Address
	ethSignerFn                keystore.SignerFn
	ethPersonalSignFn          keystore.PersonalSignFn
	relayer                    relayer.PeggyRelayer
	cosmosBlockTime            time.Duration
	ethereumBlockTime          time.Duration
	batchRequesterLoopDuration time.Duration
	startingEthBlock           uint64
	ethBlocksPerLoop           uint64

	mtx             sync.Mutex
	erc20DenomCache map[string]string
}

func NewPeggyOrchestrator(
	logger zerolog.Logger,
	cosmosQueryClient gravitytypes.QueryClient,
	gravityBroadcastClient sidechain.PeggyBroadcastClient,
	gravityContract peggy.Contract,
	ethFrom ethcmn.Address,
	ethSignerFn keystore.SignerFn,
	ethPersonalSignFn keystore.PersonalSignFn,
	relayer relayer.PeggyRelayer,
	cosmosBlockTime time.Duration,
	ethereumBlockTime time.Duration,
	batchRequesterLoopDuration time.Duration,
	ethBlocksPerLoop int64,
	options ...func(PeggyOrchestrator),
) PeggyOrchestrator {

	orch := &peggyOrchestrator{
		logger:                     logger.With().Str("module", "orchestrator").Logger(),
		cosmosQueryClient:          cosmosQueryClient,
		gravityBroadcastClient:     gravityBroadcastClient,
		gravityContract:            gravityContract,
		ethProvider:                gravityContract.Provider(),
		ethFrom:                    ethFrom,
		ethSignerFn:                ethSignerFn,
		ethPersonalSignFn:          ethPersonalSignFn,
		relayer:                    relayer,
		cosmosBlockTime:            cosmosBlockTime,
		ethereumBlockTime:          ethereumBlockTime,
		batchRequesterLoopDuration: batchRequesterLoopDuration,
		ethBlocksPerLoop:           uint64(ethBlocksPerLoop),
		startingEthBlock:           uint64(6149808),
	}

	for _, option := range options {
		option(orch)
	}

	return orch
}
