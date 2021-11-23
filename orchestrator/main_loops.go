package orchestrator

import (
	"context"
	"errors"
	"time"

	"github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum/common"
	"github.com/umee-network/peggo/orchestrator/cosmos"
	"github.com/umee-network/peggo/orchestrator/loops"
	"github.com/umee-network/umee/x/peggy/types"
)

const (
	// We define loops durations based on multipliers of average block times from Eth and Cosmos.
	//
	// Ref: https://github.com/umee-network/peggo/issues/55

	// Run every approximately 5 Ethereum blocks to allow time to receive new blocks.
	// If we run this faster we wouldn't be getting new blocks, which is not efficient.
	ethOracleLoopMultiplier = 5

	// Run every approximately 3 Cosmos blocks; so we sign batches and valset updates ASAP but not run these requests
	// too often that we make too many requests to Cosmos.
	ethSignerLoopMultiplier = 3
)

// Start combines the all major roles required to make
// up the Orchestrator, all of these are async loops.
func (p *peggyOrchestrator) Start(ctx context.Context) error {
	var pg loops.ParanoidGroup

	pg.Go(func() error {
		return p.EthOracleMainLoop(ctx)
	})
	pg.Go(func() error {
		return p.BatchRequesterLoop(ctx)
	})
	pg.Go(func() error {
		return p.EthSignerMainLoop(ctx)
	})
	pg.Go(func() error {
		return p.RelayerMainLoop(ctx)
	})

	return pg.Wait()
}

// EthOracleMainLoop is responsible for making sure that Ethereum events are retrieved from the Ethereum blockchain
// and ferried over to Cosmos where they will be used to issue tokens or process batches.
func (p *peggyOrchestrator) EthOracleMainLoop(ctx context.Context) (err error) {
	logger := p.logger.With().Str("loop", "EthOracleMainLoop").Logger()
	lastResync := time.Now()

	var (
		lastCheckedBlock uint64
		peggyParams      *types.Params
	)

	if err := retry.Do(func() (err error) {
		peggyParams, err = p.cosmosQueryClient.PeggyParams(ctx)
		if err != nil {
			logger.Fatal().Err(err).Msg("failed to query peggy params, is umeed running?")
		}

		return err
	}, retry.Context(ctx), retry.OnRetry(func(n uint, err error) {
		logger.Err(err).Uint("retry", n).Msg("failed to get Peggy params; retrying...")
	})); err != nil {
		logger.Err(err).Msg("got error, loop exits")
		return err
	}

	if err := retry.Do(func() (err error) {
		lastCheckedBlock, err = p.GetLastCheckedBlock(ctx)
		if lastCheckedBlock == 0 {
			lastCheckedBlock = peggyParams.BridgeContractStartHeight
		}

		return err
	}, retry.Context(ctx), retry.OnRetry(func(n uint, err error) {
		logger.Err(err).Uint("retry", n).Msg("failed to get last checked block; retrying...")
	})); err != nil {
		logger.Err(err).Msg("got error, loop exits")
		return err
	}

	logger.Info().Uint64("last_checked_block", lastCheckedBlock).Msg("start scanning for events")

	return loops.RunLoop(ctx, p.logger, p.ethereumBlockTime*ethOracleLoopMultiplier, func() error {
		// Relays events from Ethereum -> Cosmos
		var currentBlock uint64
		if err := retry.Do(func() (err error) {
			currentBlock, err = p.CheckForEvents(ctx, lastCheckedBlock, getEthBlockDelay(peggyParams.BridgeChainId))
			return err
		}, retry.Context(ctx), retry.OnRetry(func(n uint, err error) {
			logger.Err(err).Uint("retry", n).Msg("error during Eth event checking; retrying...")
		})); err != nil {
			logger.Err(err).Msg("got error, loop exits")
			return err
		}

		lastCheckedBlock = currentBlock

		// Auto re-sync to catch up the nonce. Reasons why event nonce fall behind.
		//	1. It takes some time for events to be indexed on Ethereum. So if peggo queried events immediately as
		//	   block produced, there is a chance the event is missed. We need to re-scan this block to ensure events
		//	   are not missed due to indexing delay.
		//	2. If validator was in UnBonding state, the claims broadcasted in last iteration are failed.
		//	3. If the ETH call failed while filtering events, the peggo missed to broadcast claim events occurred in
		//	   last iteration.
		if time.Since(lastResync) >= 48*time.Hour {
			if err := retry.Do(func() (err error) {
				lastCheckedBlock, err = p.GetLastCheckedBlock(ctx)
				return err
			}, retry.Context(ctx), retry.OnRetry(func(n uint, err error) {
				logger.Err(err).Uint("retry", n).Msg("failed to get last checked block; retrying...")
			})); err != nil {
				logger.Err(err).Msg("got error, loop exits")
				return err
			}

			lastResync = time.Now()
			logger.Info().
				Time("last_resync", lastResync).
				Uint64("last_checked_block", lastCheckedBlock).
				Msg("auto resync")
		}

		return nil
	})
}

// EthSignerMainLoop simply signs off on any batches or validator sets provided by the validator
// since these are provided directly by a trusted Cosmsos node they can simply be assumed to be
// valid and signed off on.
func (p *peggyOrchestrator) EthSignerMainLoop(ctx context.Context) (err error) {
	logger := p.logger.With().Str("loop", "EthSignerMainLoop").Logger()

	var peggyID common.Hash
	if err := retry.Do(func() (err error) {
		peggyID, err = p.peggyContract.GetPeggyID(ctx, p.peggyContract.FromAddress())
		return err
	}, retry.Context(ctx), retry.OnRetry(func(n uint, err error) {
		logger.Err(err).Uint("retry", n).Msg("failed to get PeggyID from Ethereum contract; retrying...")
	})); err != nil {
		logger.Err(err).Msg("got error, loop exits")
		return err
	}

	logger.Debug().Hex("peggyID", peggyID[:]).Msg("received peggyID")

	return loops.RunLoop(ctx, p.logger, p.cosmosBlockTime*ethSignerLoopMultiplier, func() error {
		var oldestUnsignedValsets []*types.Valset
		if err := retry.Do(func() error {
			oldestValsets, err := p.cosmosQueryClient.OldestUnsignedValsets(ctx, p.peggyBroadcastClient.AccFromAddress())
			if err != nil {
				if err == cosmos.ErrNotFound || oldestValsets == nil {
					logger.Debug().Msg("no Valset waiting to be signed")
					return nil
				}

				return err
			}

			oldestUnsignedValsets = oldestValsets
			return nil
		}, retry.Context(ctx), retry.OnRetry(func(n uint, err error) {
			logger.Err(err).Uint("retry", n).Msg("failed to get unsigned Valset for signing; retrying...")
		})); err != nil {
			logger.Err(err).Msg("got error, loop exits")
			return err
		}

		for _, oldestValset := range oldestUnsignedValsets {
			logger.Info().Uint64("oldest_valset_nonce", oldestValset.Nonce).Msg("sending Valset confirm for nonce")
			valset := oldestValset

			if err := retry.Do(func() error {
				return p.peggyBroadcastClient.SendValsetConfirm(ctx, p.ethFrom, peggyID, valset)
			}, retry.Context(ctx), retry.OnRetry(func(n uint, err error) {
				logger.Err(err).
					Uint("retry", n).
					Msg("failed to sign and send Valset confirmation to Cosmos; retrying...")
			})); err != nil {
				logger.Err(err).Msg("got error, loop exits")
				return err
			}
		}

		var oldestUnsignedTransactionBatch *types.OutgoingTxBatch
		if err := retry.Do(func() error {
			// sign the last unsigned batch, TODO check if we already have signed this
			txBatch, err := p.cosmosQueryClient.OldestUnsignedTransactionBatch(ctx, p.peggyBroadcastClient.AccFromAddress())
			if err != nil {
				if err == cosmos.ErrNotFound || txBatch == nil {
					logger.Debug().Msg("no TransactionBatch waiting to be signed")
					return nil
				}

				return err
			}

			oldestUnsignedTransactionBatch = txBatch
			return nil
		}, retry.Context(ctx), retry.OnRetry(func(n uint, err error) {
			logger.Err(err).
				Uint("retry", n).
				Msg("failed to get unsigned TransactionBatch for signing; retrying...")
		})); err != nil {
			logger.Err(err).Msg("got error, loop exits")
			return err
		}

		if oldestUnsignedTransactionBatch != nil {
			logger.Info().
				Uint64("batch_nonce", oldestUnsignedTransactionBatch.BatchNonce).
				Msg("sending TransactionBatch confirm for BatchNonce")
			if err := retry.Do(func() error {
				return p.peggyBroadcastClient.SendBatchConfirm(ctx, p.ethFrom, peggyID, oldestUnsignedTransactionBatch)
			}, retry.Context(ctx), retry.OnRetry(func(n uint, err error) {
				logger.Err(err).
					Uint("retry", n).
					Msg("failed to sign and send TransactionBatch confirmation to Cosmos; retrying...")
			})); err != nil {
				logger.Err(err).Msg("got error, loop exits")
				return err
			}
		}

		return nil
	})
}

func (p *peggyOrchestrator) BatchRequesterLoop(ctx context.Context) (err error) {
	logger := p.logger.With().Str("loop", "BatchRequesterLoop").Logger()

	// Run every approximately 60 Cosmos blocks (around 5m) to allow time to receive new transactions.
	// Running this faster will cause a lot of small batches and lots of messages going around the network.
	// We need to remember that this call is going to be made by all the validators.
	// This loop is configurable so it can be adjusted for E2E tests.
	loopDuration := time.Duration(
		float64(p.cosmosBlockTime.Milliseconds())*
			p.batchRequesterLoopMultiplier) *
		time.Millisecond

	return loops.RunLoop(ctx, p.logger, loopDuration, func() error {
		// Each loop performs the following:
		//
		// - get All the denominations
		// - broadcast Request batch
		var pg loops.ParanoidGroup

		pg.Go(func() error {
			var unbatchedTokensWithFees []*types.BatchFees

			if err := retry.Do(func() (err error) {
				unbatchedTokensWithFees, err = p.cosmosQueryClient.UnbatchedTokensWithFees(ctx)
				return
			}, retry.Context(ctx), retry.OnRetry(func(n uint, err error) {
				logger.Err(err).Uint("retry", n).Msg("failed to get UnbatchedTokensWithFees; retrying...")
			})); err != nil {
				// non-fatal, just alert
				logger.Warn().Msg("unable to get UnbatchedTokensWithFees for the token")
				return nil
			}

			for _, unbatchedToken := range unbatchedTokensWithFees {
				unbatchedToken := unbatchedToken
				tokenAddr := common.HexToAddress(unbatchedToken.Token)

				denom, err := p.ERC20ToDenom(ctx, tokenAddr)
				if err != nil {
					// do not return error, just continue with the next unbatched tx
					logger.Err(err).Str("token_contract", tokenAddr.String()).Msg("failed to get denom; will not request a batch")
					return nil
				}

				logger.Info().Str("token_contract", tokenAddr.String()).Str("denom", denom).Msg("sending batch request")

				if err := p.peggyBroadcastClient.SendRequestBatch(ctx, denom); err != nil {
					logger.Err(err).Msg("failed to send batch request")
				}
			}

			return nil
		})

		return pg.Wait()
	})
}

func (p *peggyOrchestrator) RelayerMainLoop(ctx context.Context) (err error) {
	if p.relayer != nil {
		return p.relayer.Start(ctx)
	}

	return errors.New("relayer is nil")
}

// ERC20ToDenom attempts to return the denomination that maps to an ERC20 token
// contract on the Cosmos chain. First, we check the cache. If the token address
// does not exist in the cache, we query the Cosmos chain and cache the result.
func (p *peggyOrchestrator) ERC20ToDenom(ctx context.Context, tokenAddr common.Address) (string, error) {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	tokenAddrStr := tokenAddr.String()
	denom, ok := p.erc20DenomCache[tokenAddrStr]
	if ok {
		return denom, nil
	}

	resp, err := p.cosmosQueryClient.ERC20ToDenom(ctx, tokenAddr)
	if err != nil {
		return "", err
	}

	p.erc20DenomCache[tokenAddrStr] = resp.Denom
	return resp.Denom, nil
}

// getEthBlockDelay returns the right amount of Ethereum blocks to wait until we
// consider a block final. This depends on the chain we are talking to.
// Copying from https://github.com/althea-net/cosmos-gravity-bridge/blob/main/orchestrator/orchestrator/src/ethereum_event_watcher.rs#L222
func getEthBlockDelay(chainID uint64) uint64 {
	switch chainID {
	// Mainline Ethereum, Ethereum classic, or the Ropsten, Kotti, Mordor testnets
	// all POW Chains
	case 1, 3, 6, 7:
		return 6

	// Dev, our own Gravity Ethereum testnet, and Hardhat respectively
	// all single signer chains with no chance of any reorgs
	case 2018, 15, 31337:
		return 0

	// Rinkeby and Goerli use Clique (POA) Consensus, finality takes
	// up to num validators blocks. Number is higher than Ethereum based
	// on experience with operational issues
	case 4, 5:
		return 10

	// assume the safe option (POW) where we don't know
	default:
		return 6
	}
}
