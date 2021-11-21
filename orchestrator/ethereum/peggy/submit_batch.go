package peggy

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/umee-network/umee/x/peggy/types"
)

func (s *peggyContract) EncodeTransactionBatch(
	ctx context.Context,
	currentValset *types.Valset,
	batch *types.OutgoingTxBatch,
	confirms []*types.MsgConfirmBatch,
) ([]byte, error) {

	validators, powers, sigV, sigR, sigS, err := CheckBatchSigsAndRepack(currentValset, confirms)
	if err != nil {
		err = errors.Wrap(err, "confirmations check failed")
		return nil, err
	}

	amounts, destinations, fees := getBatchCheckpointValues(batch)
	currentValsetNonce := new(big.Int).SetUint64(currentValset.Nonce)
	batchNonce := new(big.Int).SetUint64(batch.BatchNonce)
	batchTimeout := new(big.Int).SetUint64(batch.BatchTimeout)

	currentValsetArs := ValsetArgs{
		Validators:   validators,
		Powers:       powers,
		ValsetNonce:  currentValsetNonce,
		RewardAmount: currentValset.RewardAmount.BigInt(),
		RewardToken:  common.HexToAddress(currentValset.RewardToken),
	}

	txData, err := peggyABI.Pack("submitBatch",
		currentValsetArs,
		sigV, sigR, sigS,
		amounts,
		destinations,
		fees,
		batchNonce,
		common.HexToAddress(batch.TokenContract),
		batchTimeout,
	)
	if err != nil {
		s.logger.Err(err).Msg("ABI Pack (Peggy submitBatch) method")
		return nil, err
	}

	return txData, nil
}

func getBatchCheckpointValues(batch *types.OutgoingTxBatch) (
	amounts []*big.Int,
	destinations []common.Address,
	fees []*big.Int,
) {
	amounts = make([]*big.Int, len(batch.Transactions))
	destinations = make([]common.Address, len(batch.Transactions))
	fees = make([]*big.Int, len(batch.Transactions))

	for i, tx := range batch.Transactions {
		amounts[i] = tx.Erc20Token.Amount.BigInt()
		destinations[i] = common.HexToAddress(tx.DestAddress)
		fees[i] = tx.Erc20Fee.Amount.BigInt()
	}

	return
}

func CheckBatchSigsAndRepack(
	valset *types.Valset,
	confirms []*types.MsgConfirmBatch,
) (
	validators []common.Address,
	powers []*big.Int,
	v []uint8,
	r []common.Hash,
	s []common.Hash,
	err error,
) {
	if len(confirms) == 0 {
		err = errors.New("no signatures in batch confirmation")
		return
	}

	signerToSig := make(map[string]*types.MsgConfirmBatch, len(confirms))
	for _, sig := range confirms {
		signerToSig[sig.EthSigner] = sig
	}

	powerOfGoodSigs := new(big.Int)

	for _, m := range valset.Members {
		mPower := big.NewInt(0).SetUint64(m.Power)
		if sig, ok := signerToSig[m.EthereumAddress]; ok && sig.EthSigner == m.EthereumAddress {
			powerOfGoodSigs.Add(powerOfGoodSigs, mPower)

			validators = append(validators, common.HexToAddress(m.EthereumAddress))
			powers = append(powers, mPower)

			sigV, sigR, sigS := sigToVRS(sig.Signature)
			v = append(v, sigV)
			r = append(r, sigR)
			s = append(s, sigS)
		} else {
			validators = append(validators, common.HexToAddress(m.EthereumAddress))
			powers = append(powers, mPower)
			v = append(v, 0)
			r = append(r, [32]byte{})
			s = append(s, [32]byte{})
		}
	}
	if peggyPowerToPercent(powerOfGoodSigs) < 66 {
		err = ErrInsufficientVotingPowerToPass
		return validators, powers, v, r, s, err
	}

	return validators, powers, v, r, s, err
}
