package relayer

import (
	"context"
	"sort"
	"time"

	cosmtypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/umee-network/umee/x/peggy/types"
)

// RelayBatches gets the last batch of outgoing transactions and relays it to Ethereum in a TX.
// Validators will only relay a batch if they consider it profitable, that is, if the total fees are over its
// minimum-batch-fee parameter. If the batch isn't profitable for this validator, we'll give it some time (relayTimeout)
// to allow any other validator who considers it profitable to relay it. If after relayTimeout the batch hasn't being
// sent, we'll skip to the next one.
// relayTimeout is different from BatchTimeout. BatchTimeout is set 12hrs after the batch is created and relayTimeout is
// much shorter (2x relay loop) and it will only be used if the batch is not profitable and there is another batch
// in line waiting. This is to prevent the bridge from being halted for 12h if a validator erronously requested an
// unprofitable batch.
// Any transactions left in the unsent/skipped batch will be put back in the queue (this is handled by x/peggy).
func (s *peggyRelayer) RelayBatches(ctx context.Context) error {
	latestBatches, err := s.cosmosQueryClient.LatestTransactionBatches(ctx)
	if err != nil {
		return err
	}

	var selectedBatch *types.OutgoingTxBatch
	var selectedBatchSigs []*types.MsgConfirmBatch

	// order batches by nonce ASC. That means that the next batch is [0].
	sort.SliceStable(latestBatches, func(i, j int) bool {
		return latestBatches[i].BatchNonce > latestBatches[j].BatchNonce
	})

	for _, batch := range latestBatches {
		sigs, err := s.cosmosQueryClient.TransactionBatchSignatures(
			ctx,
			batch.BatchNonce,
			common.HexToAddress(batch.TokenContract),
		)
		if err != nil {
			return err
		}

		// if the batch is signed by the latest validator set, we can relay it
		if len(sigs) != 0 {

			// Check if the batch is profitable
			if !s.IsBatchProfitable(ctx, batch, s.minBatchFeeUSD) {
				// it's not profitable, check if we should wait for it (return)
				// or if it timeoutted, skip to the next batch

				// get Cosmos' latest block height and Peggy's average block time
				// to calculate how much time has passed since the batch was created
				latestBlockHeight, err := s.tmClient.GetLatestBlockHeight(ctx)
				if err != nil {
					s.logger.Err(err).Msg("failed to get latest block height")
					return err
				}

				peggyParams, err := s.cosmosQueryClient.PeggyParams(ctx)
				if err != nil {
					s.logger.Err(err).Msg("failed to query peggy params, is umeed running?")
				}

				elapsedBlocks := uint64(latestBlockHeight) - batch.Block
				elapsedTime := time.Duration(elapsedBlocks*peggyParams.AverageBlockTime) * time.Millisecond

				// if the batch has been created more than relayTimeout ago, skip it and try to send the next one
				if elapsedTime > (s.loopDuration * 2) {
					continue
				}

				// if the batch has been created less than relayTimeout ago, wait for it
				return nil
			}

			// batch is profitable and has signatures, we can relay it
			selectedBatch = batch
			selectedBatchSigs = sigs
			break
		}
	}

	if selectedBatch == nil {
		s.logger.Debug().Msg("could not find batch with signatures, nothing to relay")
		return nil
	}

	latestEthereumBatch, err := s.peggyContract.GetTxBatchNonce(
		ctx,
		common.HexToAddress(selectedBatch.TokenContract),
		s.peggyContract.FromAddress(),
	)
	if err != nil {
		return err
	}

	currentValset, err := s.FindLatestValset(ctx)
	if err != nil {
		return errors.New("failed to find latest valset")
	} else if currentValset == nil {
		return errors.New("latest valset not found")
	}

	s.logger.Debug().
		Uint64("oldest_batch_nonce", selectedBatch.BatchNonce).
		Uint64("latest_batch_nonce", latestEthereumBatch.Uint64()).
		Msg("found latest valsets")

	if selectedBatch.BatchNonce > latestEthereumBatch.Uint64() {

		latestEthereumBatch, err := s.peggyContract.GetTxBatchNonce(
			ctx,
			common.HexToAddress(selectedBatch.TokenContract),
			s.peggyContract.FromAddress(),
		)
		if err != nil {
			return err
		}
		// Check if oldestSignedBatch already submitted by other validators in mean time
		if selectedBatch.BatchNonce > latestEthereumBatch.Uint64() {
			s.logger.Info().
				Uint64("latest_batch", selectedBatch.BatchNonce).
				Uint64("latest_ethereum_batch", latestEthereumBatch.Uint64()).
				Msg("we have detected latest batch but Ethereum has a different one. Sending an update!")

			// Send SendTransactionBatch to Ethereum
			txHash, err := s.peggyContract.SendTransactionBatch(ctx, currentValset, selectedBatch, selectedBatchSigs)
			if err != nil {
				return err
			}
			s.logger.Info().Str("tx_hash", txHash.Hex()).Msg("sent Ethereum Tx (TransactionBatch)")
		}
	}

	return nil
}

func (s *peggyRelayer) IsBatchProfitable(
	ctx context.Context,
	batch *types.OutgoingTxBatch,
	minFeeInUSD float64,
) bool {
	if minFeeInUSD == 0 || s.priceFeeder == nil {
		return true
	}

	decimals, err := s.peggyContract.GetERC20Decimals(
		ctx,
		common.HexToAddress(batch.TokenContract),
		s.peggyContract.FromAddress(),
	)
	if err != nil {
		s.logger.Err(err).Str("token_contract", batch.TokenContract).Msg("failed to get token decimals")
		return false
	}

	s.logger.Debug().
		Uint8("decimals", decimals).
		Str("token_contract", batch.TokenContract).
		Msg("got token decimals")

	tokenPriceInUSD, err := s.priceFeeder.QueryUSDPrice(common.HexToAddress(batch.TokenContract))
	if err != nil {
		return false
	}

	// We calculate the total fee in ERC20 tokens
	totalBatchFees := cosmtypes.Int{}
	for _, tx := range batch.Transactions {
		totalBatchFees.Add(tx.Erc20Fee.Amount)
	}

	tokenPriceInUSDDec := decimal.NewFromFloat(tokenPriceInUSD)
	// decimals (uint8) can be safely casted into int32 because the max uint8 is 255 and the max int32 is 2147483647
	totalFeeInUSDDec := decimal.NewFromBigInt(totalBatchFees.BigInt(), -int32(decimals)).Mul(tokenPriceInUSDDec)
	minFeeInUSDDec := decimal.NewFromFloat(minFeeInUSD)

	s.logger.Debug().
		Str("token_contract", batch.TokenContract).
		Float64("token_price_in_usd", tokenPriceInUSD).
		Int64("total_fees", totalBatchFees.Int64()).
		Float64("total_fee_in_usd", totalFeeInUSDDec.InexactFloat64()).
		Float64("min_fee_in_usd", minFeeInUSDDec.InexactFloat64()).
		Msg("checking if token fees meet minimum batch fee threshold")

	return totalFeeInUSDDec.GreaterThan(minFeeInUSDDec)
}
