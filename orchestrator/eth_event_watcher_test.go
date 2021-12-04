package orchestrator

import (
	"context"
	"math/big"
	"os"
	"testing"
	"time"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/umee-network/peggo/mocks"
	"github.com/umee-network/peggo/orchestrator/ethereum/committer"
	"github.com/umee-network/peggo/orchestrator/ethereum/peggy"
	wrappers "github.com/umee-network/peggo/solidity/wrappers/Peggy.sol"
)

// TODO: This function will require quite some effort to get it tested.
func TestCheckForEvents(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	fromAddress := ethcmn.HexToAddress("0xd8da6bf26964af9d7eed9e03e53415d37aa96045")
	peggyAddress := ethcmn.HexToAddress("0x3bdf8428734244c9e5d82c95d125081939d6d42d")
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr})

	ethProvider := mocks.NewMockEVMProviderWithRet(mockCtrl)
	ethProvider.EXPECT().PendingNonceAt(gomock.Any(), fromAddress).Return(uint64(0), nil)
	ethProvider.EXPECT().HeaderByNumber(gomock.Any(), nil).Return(&ethtypes.Header{
		Number: big.NewInt(100),
	}, nil)
	ethProvider.EXPECT().FilterLogs(gomock.Any(), gomock.Any()).Return([]ethtypes.Log{}, errors.New("some error")).AnyTimes()

	ethGasPriceAdjustment := 1.0
	ethCommitter, _ := committer.NewEthCommitter(
		logger,
		fromAddress,
		ethGasPriceAdjustment,
		nil,
		ethProvider,
	)

	peggyContract, _ := peggy.NewPeggyContract(logger, ethCommitter, peggyAddress)

	orch := NewPeggyOrchestrator(
		zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}),
		nil,
		nil,
		peggyContract,
		fromAddress,
		nil,
		nil,
		nil,
		time.Second,
		time.Second,
		time.Second,
		100,
	)

	orch.CheckForEvents(context.Background(), 0, 1)

}

func TestFilterSendToCosmosEventsByNonce(t *testing.T) {
	// In testEv we'll add 2 valid and 1 past event.
	// This should result in only 2 events after the filter.
	testEv := []*wrappers.PeggySendToCosmosEvent{
		{EventNonce: big.NewInt(3)},
		{EventNonce: big.NewInt(4)},
		{EventNonce: big.NewInt(5)},
	}
	nonce := uint64(3)

	assert.Len(t, filterSendToCosmosEventsByNonce(testEv, nonce), 2)
}

func TestFilterTransactionBatchExecutedEventsByNonce(t *testing.T) {
	// In testEv we'll add 2 valid and 1 past event.
	// This should result in only 2 events after the filter.
	testEv := []*wrappers.PeggyTransactionBatchExecutedEvent{
		{EventNonce: big.NewInt(3)},
		{EventNonce: big.NewInt(4)},
		{EventNonce: big.NewInt(5)},
	}
	nonce := uint64(3)

	assert.Len(t, filterTransactionBatchExecutedEventsByNonce(testEv, nonce), 2)
}

func TestFilterValsetUpdateEventsByNonce(t *testing.T) {
	// In testEv we'll add 2 valid and 1 past event.
	// This should result in only 2 events after the filter.
	testEv := []*wrappers.PeggyValsetUpdatedEvent{
		{EventNonce: big.NewInt(3)},
		{EventNonce: big.NewInt(4)},
		{EventNonce: big.NewInt(5)},
	}
	nonce := uint64(3)

	assert.Len(t, filterValsetUpdateEventsByNonce(testEv, nonce), 2)
}

func TestFilterERC20DeployedEventsByNonce(t *testing.T) {
	// In testEv we'll add 2 valid and 1 past event.
	// This should result in only 2 events after the filter.
	testEv := []*wrappers.PeggyERC20DeployedEvent{
		{EventNonce: big.NewInt(3)},
		{EventNonce: big.NewInt(4)},
		{EventNonce: big.NewInt(5)},
	}
	nonce := uint64(3)

	assert.Len(t, filterERC20DeployedEventsByNonce(testEv, nonce), 2)
}

func TestIsUnknownBlockErr(t *testing.T) {
	gethErr := errors.New("unknown block")
	assert.True(t, isUnknownBlockErr(gethErr))

	parityErr := errors.New("One of the blocks specified in filter...")
	assert.True(t, isUnknownBlockErr(parityErr))

	otherErr := errors.New("other error")
	assert.False(t, isUnknownBlockErr(otherErr))
}
