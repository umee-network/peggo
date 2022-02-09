package txanalyzer

import (
	"log"
	"math/big"
	"os"
	"testing"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/umee-network/peggo/orchestrator/ethereum/provider"

	wrappers "github.com/umee-network/peggo/solwrappers/Gravity.sol"
)

func TestTXAnalyzer(t *testing.T) {
	ethRPC, err := ethrpc.Dial("https://goerli.infura.io/v3/")
	if err != nil {
		log.Fatal(err)
	}

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr})

	txAnalyzer, err := NewTXAnalyzer(
		logger,
		"",
		provider.NewEVMProvider(ethRPC),
		ethcmn.HexToAddress("0x0000000000000000000000000000000000000001"),
		200000,
		6000000,
	)

	assert.Nil(t, err)

	// Add some TXs to be processed
	evt := &wrappers.GravityTransactionBatchExecutedEvent{
		BatchNonce: &big.Int{},
		Token:      ethcmn.HexToAddress("0xe54fbaecc50731afe54924c40dfd1274f718fe02"),
		EventNonce: &big.Int{},
		Raw:        types.Log{TxHash: ethcmn.HexToHash("0x354ee8700da020fb1d1794ad9bed5c82b63274b3b2c61db0f88edc70333bb13a")},
	}

	err = txAnalyzer.StoreBatches([]*wrappers.GravityTransactionBatchExecutedEvent{evt})
	assert.Nil(t, err)

	// Query the unprocessed txs
	unprocessedTxs, err := txAnalyzer.GetUnprocessedTXsByToken()
	assert.Nil(t, err)
	assert.Len(t, unprocessedTxs[ethcmn.HexToAddress("0xe54fbaecc50731afe54924c40dfd1274f718fe02")], 1)

	// Process the txs
	assert.Nil(t, txAnalyzer.ProcessTXs(unprocessedTxs))

	// Now they should be processed, so no unprocessed txs should be returned
	unprocessedTxs, err = txAnalyzer.GetUnprocessedTXsByToken()
	assert.Nil(t, err)
	assert.Len(t, unprocessedTxs, 0)

	deletedTxs, err := txAnalyzer.PruneTXs()
	assert.Nil(t, err)
	assert.Equal(t, 0, deletedTxs)

	assert.Nil(t, txAnalyzer.RecalculateEstimates())

	estimates, err := txAnalyzer.GetEstimatesOfToken(ethcmn.HexToAddress("0xe54fbaecc50731afe54924c40dfd1274f718fe02"))
	assert.Nil(t, err)
	assert.Equal(t, uint64(1), estimates[99][0])
	assert.Equal(t, uint64(1245478), estimates[99][1])
	assert.Equal(t, uint64(0), estimates[98][0])
	assert.Equal(t, uint64(1258859), estimates[98][1])

	assert.Nil(t, txAnalyzer.Close())

}
