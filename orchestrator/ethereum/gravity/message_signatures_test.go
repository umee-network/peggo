package gravity

import (
	"testing"

	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravity/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestEncodeValsetConfirm(t *testing.T) {
	gravityID := "defaultgravityid"

	valset := types.Valset{
		Nonce: 5,
		Members: []types.BridgeValidator{
			{Power: 1, EthereumAddress: "0x02fa1b44e2EF8436e6f35D5F56607769c658c225"},
			{Power: 123, EthereumAddress: "0x4f3a9f8f8f8f8f8f8f8f8f8f8f8f8f8f8f8f8f8f8f"},
		},
		Height:       111111,
		RewardAmount: sdk.NewInt(2),
		RewardToken:  "",
	}

	result := EncodeValsetConfirm(gravityID, valset)

	// Check the result with a previously calculated one.
	assert.Equal(t, "0x530516ded1a45852c4000d36e5da715a934b8f272ed09e70b049c78474f8343b", result.Hex())
}

func TestEncodeTxBatchConfirm(t *testing.T) {
	gravityID := "defaultgravityid"

	txBatch := types.OutgoingTxBatch{
		Transactions: []types.OutgoingTransferTx{
			{
				DestAddress: "0x02fa1b44e2EF8436e6f35D5F56607769c658c225",
				Erc20Token: types.ERC20Token{
					Contract: "0x4884e2a214dc5040f52a41c3f21c765283170b6e",
					Amount:   sdk.NewInt(100000),
				},
				Erc20Fee: types.ERC20Token{
					Contract: "0x4884e2a214dc5040f52a41c3f21c765283170b6e",
					Amount:   sdk.NewInt(2000),
				},
			},
		},
	}

	result := EncodeTxBatchConfirm(gravityID, txBatch)

	// Check the result with a previously calculated one.
	assert.Equal(t, "0x2c8418bc8093a21b04e82d0527b039084bca48cbbb6d413011a98181f7af5081", ethcmn.Bytes2Hex(result))
}
