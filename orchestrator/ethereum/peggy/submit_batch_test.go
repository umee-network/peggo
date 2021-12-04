package peggy

import (
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/umee-network/umee/x/peggy/types"
)

func TestGetBatchCheckpointValues(t *testing.T) {
	batch := &types.OutgoingTxBatch{
		Transactions: []*types.OutgoingTransferTx{
			{
				DestAddress: common.HexToAddress("0x2").Hex(),
				Erc20Token: &types.ERC20Token{
					Contract: common.HexToAddress("0x1").Hex(),
					Amount:   sdk.NewInt(10000),
				},
				Erc20Fee: &types.ERC20Token{
					Contract: common.HexToAddress("0x1").Hex(),
					Amount:   sdk.NewInt(100),
				},
			},
		},
	}

	amounts, destinations, fees := getBatchCheckpointValues(batch)
	assert.Equal(t, []*big.Int{big.NewInt(10000)}, amounts)
	assert.Equal(t, []common.Address{common.HexToAddress("0x2")}, destinations)
	assert.Equal(t, []*big.Int{big.NewInt(100)}, fees)
}

func TestCheckBatchSigsAndRepack(t *testing.T) {
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

	confirms := []*types.MsgConfirmBatch{
		{
			EthSigner: common.HexToAddress("0x0").Hex(),
			Signature: "0xaae54ee7e285fbb0275279143abc4c554e5314e7b417ecac83a5984a964facbaad68866a2841c3e83ddf125a2985566261c4014f9f960ec60253aebcda9513a9b4",
		},
		{
			EthSigner: common.HexToAddress("0x1").Hex(),
			Signature: "0xaae54ee7e285fbb0275279143abc4c554e5314e7b417ecac83a5984a964facbaad68866a2841c3e83ddf125a2985566261c4014f9f960ec60253aebcda9513a9b4",
		},
	}

	repackedSigs, err := checkBatchSigsAndRepack(valset, confirms)
	assert.Nil(t, err)

	assert.Equal(t, []common.Address{common.HexToAddress("0x0"), common.HexToAddress("0x1"), common.HexToAddress("0x2")}, repackedSigs.validators)
	assert.Equal(t, []*big.Int{big.NewInt(1111111111), big.NewInt(2212121212), big.NewInt(123456)}, repackedSigs.powers)

}
