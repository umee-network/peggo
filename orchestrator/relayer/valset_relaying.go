package relayer

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/umee-network/umee/x/peggy/types"
)

// RelayValsets checks the last validator set on Ethereum, if it's lower than our latest validator
// set then we should package and submit the update as an Ethereum transaction
func (s *peggyRelayer) RelayValsets(ctx context.Context, currentValset *types.Valset) error {
	// we should determine if we need to relay one
	// to Ethereum for that we will find the latest confirmed valset and compare it to the ethereum chain
	latestValsets, err := s.cosmosQueryClient.LatestValsets(ctx)
	if err != nil {
		err = errors.Wrap(err, "failed to fetch latest valsets from cosmos")
		return err
	}

	var latestCosmosSigs []*types.MsgValsetConfirm
	var latestCosmosConfirmed *types.Valset
	for _, set := range latestValsets {
		sigs, err := s.cosmosQueryClient.AllValsetConfirms(ctx, set.Nonce)
		if err != nil {
			err = errors.Wrapf(err, "failed to get valset confirms at nonce %d", set.Nonce)
			return err
		} else if len(sigs) == 0 {
			continue
		}

		latestCosmosSigs = sigs
		latestCosmosConfirmed = set
		break
	}

	if latestCosmosConfirmed == nil {
		s.logger.Debug().Msg("no confirmed valsets found, nothing to relay")
		return nil
	}

	s.logger.Debug().
		Uint64("current_eth_valset_nonce", currentValset.Nonce).
		Uint64("latest_cosmos_confirmed_nonce", latestCosmosConfirmed.Nonce).
		Msg("found latest valsets")

	if latestCosmosConfirmed.Nonce > currentValset.Nonce {
		latestEthereumValsetNonce, err := s.peggyContract.GetValsetNonce(ctx, s.peggyContract.FromAddress())
		if err != nil {
			err = errors.Wrap(err, "failed to get latest Valset nonce")
			return err
		}

		// Check if latestCosmosConfirmed already submitted by other validators in mean time
		if latestCosmosConfirmed.Nonce > latestEthereumValsetNonce.Uint64() {
			s.logger.Info().
				Uint64("latest_cosmos_confirmed_nonce", latestCosmosConfirmed.Nonce).
				Uint64("latest_ethereum_valset_nonce", latestEthereumValsetNonce.Uint64()).
				Msg("detected latest cosmos valset nonce, but latest valset on Ethereum is different. Sending update to Ethereum")

			txData, err := s.peggyContract.EncodeValsetUpdate(
				ctx,
				currentValset,
				latestCosmosConfirmed,
				latestCosmosSigs,
			)
			if err != nil {
				return err
			}

			//TODO: estimate gas and profitability using "valset reward" param

			// Checking in pending txs (mempool) if tx with same input is already submitted
			// We have to check this at the very last moment because any other relayer could have submitted
			// TODO: remove hardcoded time
			if s.peggyContract.IsPendingTxInput(txData, time.Minute) {
				s.logger.Error().
					Msg("Transaction with same valset input data is already present in mempool")
				return nil
			}

			// Send Valset Update to Ethereum
			txHash, err := s.peggyContract.SendTx(ctx, s.peggyContract.Address(), txData)
			if err != nil {
				s.logger.Err(err).
					Str("tx_hash", txHash.Hex()).
					Msg("failed to sign and submit (Peggy updateValset) to EVM")
				return err
			}

			s.logger.Info().Str("tx_hash", txHash.Hex()).Msg("sent Tx (Peggy updateValset)")

		}

	}

	return nil
}
