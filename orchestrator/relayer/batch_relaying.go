package relayer

import (
	"context"
	"sort"

	cosmtypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"github.com/umee-network/peggo/orchestrator/ethereum/peggy"
	"github.com/umee-network/umee/x/peggy/types"
)

type SubmittableBatch struct {
	Batch      *types.OutgoingTxBatch
	Signatures []*types.MsgConfirmBatch
}

// getBatchesAndSignatures retrieves the latest batches from the Cosmos module and then iterates through the signatures
// for each batch, determining if they are ready to submit. It is possible for a batch to not have valid signatures for
// two reasons one is that not enough signatures have been collected yet from the validators two is that the batch is
// old enough that the signatures do not reflect the current validator set on Ethereum. In both the later and the former
// case the correct solution is to wait through timeouts, new signatures, or a later valid batch being submitted old
// batches will always be resolved.
func (s *peggyRelayer) getBatchesAndSignatures(
	ctx context.Context,
	currentValset *types.Valset,
) (map[common.Address][]SubmittableBatch, error) {
	possibleBatches := map[common.Address][]SubmittableBatch{}

	latestBatches, err := s.cosmosQueryClient.LatestTransactionBatches(ctx)
	if err != nil {
		s.logger.Err(err).Msg("failed to get latest batches")
		return possibleBatches, err
	}

	for _, batch := range latestBatches {
		// get the signatures for the batch
		sigs, err := s.cosmosQueryClient.TransactionBatchSignatures(
			ctx,
			batch.BatchNonce,
			common.HexToAddress(batch.TokenContract),
		)
		if err != nil {
			// If we can't get the signatures for a batch we will continue to the next batch
			s.logger.Err(err).
				Uint64("batch_nonce", batch.BatchNonce).
				Str("token_contract", batch.TokenContract).
				Msg("failed to get batch's signatures")
			continue
		}

		// this checks that the signatures for the batch are actually possible to submit to the chain
		// we only need to know if the signatures are good, we won't use the other returned values
		_, _, _, _, _, err = peggy.CheckBatchSigsAndRepack(currentValset, sigs)

		if err != nil {
			// this batch is not ready to be relayed
			s.logger.
				Debug().
				AnErr("err", err).
				Uint64("batch_nonce", batch.BatchNonce).
				Str("token_contract", batch.TokenContract).
				Msg("batch can't be submitted yet, waiting for more signatures")
		}

		// if the previous check didn't fail, we can add the batch to the list of possible batches
		possibleBatches[common.HexToAddress(batch.TokenContract)] = append(
			possibleBatches[common.HexToAddress(batch.TokenContract)],
			SubmittableBatch{Batch: batch, Signatures: sigs},
		)
	}

	// order batches by nonce ASC. That means that the next/oldest batch is [0].
	for tokenAddress := range possibleBatches {
		address := tokenAddress // use this because of scopelint
		sort.SliceStable(possibleBatches[address], func(i, j int) bool {
			return possibleBatches[address][i].Batch.BatchNonce > possibleBatches[address][j].Batch.BatchNonce
		})
	}

	return possibleBatches, nil
}

// RelayBatches ttempts to submit batches with valid signatures, checking the state of the Ethereum chain to ensure that
// it is valid to submit a given batch more specifically that the correctly signed batch has not timed out or already
// been submitted. The goal of this function is to submit batches in chronological order of their creation, submitting
// batches newest first will invalidate old batches and is less efficient if those old batches are profitable.
// This function estimates the cost of submitting a batch before actually submitting it to Ethereum, if it is determined
// that the ETH cost to submit is too high the batch will be skipped and a later, more profitable, batch may be
// submitted.
// Keep in mind that many other relayers are making this same computation and some may have different standards for
// their profit margin, therefore there may be a race not only to submit individual batches but also batches in
// different orders
func (s *peggyRelayer) RelayBatches(
	ctx context.Context,
	currentValset *types.Valset,
	possibleBatches map[common.Address][]SubmittableBatch,
) error {
	// first get current block height to check for any timeouts
	lastEthereumHeader, err := s.ethProvider.HeaderByNumber(ctx, nil)
	if err != nil {
		s.logger.Err(err).Msg("failed to get last ethereum header")
		return err
	}

	ethBlockHeight := lastEthereumHeader.Number.Uint64()

	for tokenContract, batches := range possibleBatches {

		// requests data from Ethereum only once per token type, this is valid because we are
		// iterating from oldest to newest, so submitting a batch earlier in the loop won't
		// ever invalidate submitting a batch later in the loop. Another relayer could always
		// do that though.
		latestEthereumBatch, err := s.peggyContract.GetTxBatchNonce(
			ctx,
			tokenContract,
			s.peggyContract.FromAddress(),
		)
		if err != nil {
			s.logger.Err(err).Msg("failed to get latest Ethereum batch")
			return err
		}

		// now we iterate through batches per token type
		for _, batch := range batches {

			if batch.Batch.BatchTimeout < ethBlockHeight {
				s.logger.Debug().
					Uint64("batch_nonce", batch.Batch.BatchNonce).
					Str("token_contract", batch.Batch.TokenContract).
					Msg("batch has timed out and can't be submitted")
				continue
			}

			// if the batch is newer than the latest Ethereum batch, we can submit it
			if batch.Batch.BatchNonce <= latestEthereumBatch.Uint64() {
				continue
			}

			// TODO: estimate gas cost and check if this tx is profitable
			// If the batch is not profitable, move on to the next one.
			if !s.IsBatchProfitable(ctx, batch.Batch, s.minBatchFeeUSD) {
				continue
			}

			s.logger.Info().
				Uint64("latest_batch", batch.Batch.BatchNonce).
				Uint64("latest_ethereum_batch", latestEthereumBatch.Uint64()).
				Msg("we have detected latest batch but Ethereum has a different one. Sending an update!")

			// Send SendTransactionBatch to Ethereum
			txHash, err := s.peggyContract.SendTransactionBatch(ctx, currentValset, batch.Batch, batch.Signatures)
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
