package relayer

import (
	"context"
	"math/big"
	"testing"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/umee-network/peggo/mocks"
	peggyMocks "github.com/umee-network/peggo/mocks/peggy"
	"github.com/umee-network/umee/x/peggy/types"
)

func TestRelayValsets(t *testing.T) {
	t.Run("ok", func(t *testing.T) {

		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		mockQClient := mocks.NewMockQueryClient(mockCtrl)
		mockQClient.EXPECT().
			LastValsetRequests(gomock.Any(), &types.QueryLastValsetRequestsRequest{}).
			Return(&types.QueryLastValsetRequestsResponse{
				Valsets: []*types.Valset{
					{
						Nonce: 3,
						Members: []*types.BridgeValidator{
							{
								Power:           1000,
								EthereumAddress: "0x0000000000000000000000000000000000000000",
							},
							{
								Power:           1000,
								EthereumAddress: "0x1000000000000000000000000000000000000000",
							},
						},
						Height: 0,
					},
				},
			}, nil)

		mockQClient.EXPECT().ValsetConfirmsByNonce(
			gomock.Any(),
			&types.QueryValsetConfirmsByNonceRequest{
				Nonce: 3,
			}).Return(&types.QueryValsetConfirmsByNonceResponse{
			Confirms: []*types.MsgValsetConfirm{
				{
					Nonce:        0,
					Orchestrator: "aaa",
					EthAddress:   "0x0000000000000000000000000000000000000000",
					Signature:    "0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
				},
			},
		}, nil)

		peggyAddress := ethcmn.HexToAddress("0x3bdf8428734244c9e5d82c95d125081939d6d42d")
		fromAddress := ethcmn.HexToAddress("0xd8da6bf26964af9d7eed9e03e53415d37aa96045")

		mockPeggyContract := peggyMocks.NewMockContract(mockCtrl)
		mockPeggyContract.EXPECT().GetValsetNonce(gomock.Any(), fromAddress).Return(big.NewInt(2), nil)
		mockPeggyContract.EXPECT().FromAddress().Return(fromAddress).AnyTimes()
		mockPeggyContract.EXPECT().
			EncodeValsetUpdate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return([]byte{1, 2, 3}, nil)
		mockPeggyContract.EXPECT().Address().Return(peggyAddress).AnyTimes()
		mockPeggyContract.EXPECT().
			EstimateGas(gomock.Any(), peggyAddress, []byte{1, 2, 3}).
			Return(uint64(1000), big.NewInt(100), nil)

		mockPeggyContract.EXPECT().IsPendingTxInput([]byte{1, 2, 3}, gomock.Any()).Return(false)

		mockPeggyContract.EXPECT().SendTx(
			gomock.Any(),
			peggyAddress,
			[]byte{1, 2, 3},
			uint64(1000),
			big.NewInt(100),
		).Return(ethcmn.HexToHash("0x01010101"), nil)

		relayer := peggyRelayer{
			peggyContract:     mockPeggyContract,
			cosmosQueryClient: mockQClient,
		}

		relayer.RelayValsets(context.Background(), &types.Valset{})
	})
}
