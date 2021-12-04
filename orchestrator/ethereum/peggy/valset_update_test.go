package peggy

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/umee-network/umee/x/peggy/types"
)

func TestEncodeValsetUpdate(t *testing.T) {
	// TODO
}

func TestValidatorsAndPowers(t *testing.T) {
	valset := &types.Valset{
		Members: []*types.BridgeValidator{
			{
				EthereumAddress: common.HexToAddress("0x0").Hex(),
				Power:           123456,
			},
			{
				EthereumAddress: common.HexToAddress("0x1").Hex(),
				Power:           7891011,
			},
		},
	}
	validators, powers := validatorsAndPowers(valset)

	expectedValidators := []common.Address{
		common.HexToAddress("0x0"),
		common.HexToAddress("0x1"),
	}

	expectedPowers := []*big.Int{
		big.NewInt(123456),
		big.NewInt(7891011),
	}

	assert.Equal(t, expectedValidators, validators)
	assert.Equal(t, expectedPowers, powers)

}

func TestCheckValsetSigsAndRepack(t *testing.T) {
	// TODO: These are not real signatures. Would be cool to use real data here.

	valset := &types.Valset{
		Members: []*types.BridgeValidator{
			{
				EthereumAddress: common.HexToAddress("0x0").Hex(),
				Power:           1111111111,
			},
			{
				EthereumAddress: common.HexToAddress("0x1").Hex(),
				Power:           2212121212,
			},
			{
				EthereumAddress: common.HexToAddress("0x2").Hex(),
				Power:           123456,
			},
		},
	}

	confirms := []*types.MsgValsetConfirm{
		{
			EthAddress: common.HexToAddress("0x0").Hex(),
			Signature:  "0xaae54ee7e285fbb0275279143abc4c554e5314e7b417ecac83a5984a964facbaad68866a2841c3e83ddf125a2985566261c4014f9f960ec60253aebcda9513a9b4",
		},
		{
			EthAddress: common.HexToAddress("0x1").Hex(),
			Signature:  "0xaae54ee7e285fbb0275279143abc4c554e5314e7b417ecac83a5984a964facbaad68866a2841c3e83ddf125a2985566261c4014f9f960ec60253aebcda9513a9b4",
		},
	}

	repackedSigs, err := checkValsetSigsAndRepack(valset, confirms)
	assert.Nil(t, err)

	assert.Equal(t, []common.Address{common.HexToAddress("0x0"), common.HexToAddress("0x1"), common.HexToAddress("0x2")}, repackedSigs.validators)
	assert.Equal(t, []*big.Int{big.NewInt(1111111111), big.NewInt(2212121212), big.NewInt(123456)}, repackedSigs.powers)

}
