package txanalyzer

import (
	"math/big"
	"testing"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	wrappers "github.com/umee-network/peggo/solwrappers/Gravity.sol"
)

func TestTXAnalyzer(t *testing.T) {
	txAnalyzer, err := NewTXAnalyzer("")

	assert.Nil(t, err)

	// Add some TXs to be processed
	evt := wrappers.GravityTransactionBatchExecutedEvent{
		BatchNonce: &big.Int{},
		Token:      ethcmn.HexToAddress("0xe54fbaecc50731afe54924c40dfd1274f718fe02"),
		EventNonce: &big.Int{},
		Raw:        types.Log{TxHash: ethcmn.HexToHash("0x354ee8700da020fb1d1794ad9bed5c82b63274b3b2c61db0f88edc70333bb13a")},
	}

	err = txAnalyzer.StoreBatches([]wrappers.GravityTransactionBatchExecutedEvent{evt})
	assert.Nil(t, err)

	// Query the unprocessed txs
	unprocessedTxs, err := txAnalyzer.GetUnprocessedRawTXs(ethcmn.HexToAddress("0xe54fbaecc50731afe54924c40dfd1274f718fe02"))
	assert.Nil(t, err)
	assert.Len(t, unprocessedTxs, 1)

	// Process the txs
	assert.Nil(t, txAnalyzer.ProcessTXs(unprocessedTxs))

	// Now they should be processed, so no unprossed txs should be returned
	unprocessedTxs, err = txAnalyzer.GetUnprocessedRawTXs(ethcmn.HexToAddress("0xe54fbaecc50731afe54924c40dfd1274f718fe02"))
	assert.Nil(t, err)
	assert.Len(t, unprocessedTxs, 0)

}
