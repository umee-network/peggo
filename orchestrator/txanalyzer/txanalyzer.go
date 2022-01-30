package txanalyzer

import (
	"context"
	"errors"
	"log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	badger "github.com/dgraph-io/badger/v3"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	wrappers "github.com/umee-network/peggo/solwrappers/Gravity.sol"
)

// KVStore key prefixes
var (
	KeyPrefixRawTx          = []byte{0x01}
	KeyPrefixHistoricalData = []byte{0x02}
)

type TXAnalyzer struct {
	db     *badger.DB
	ethRPC *ethclient.Client
}

// type RawTX struct {
// 	GasUsed uint64 `json:"gu"`
// 	TxCount uint8  `json:"tc"` // uint8 is enough as our max is 100 txs per batch
// }

// type ProcessedData struct {
// 	GasUsedPerTXCount []uint64 `json:"gm"`
// }

/*

2 things are going to get stored:
- "Raw" Ethereum TXs, containing tx hash, gas usage, token address, amount of outgoing txs, etc.
 {prefix}{erc20 contract address}{txhash}:struct
- "Processed data": per token address we'll have the estimated gas usage for each amount of outgoing tx (from 1 to 100).
  This "estimation" is going to be the running average of the gas usage of the last 100 matching transactions.
  Any hole in the data is going to be filled in using near values to guesstimate the missing data.

*/

func NewTXAnalyzer(dbDir string) (*TXAnalyzer, error) {
	db, err := badger.Open(badger.DefaultOptions(dbDir).WithInMemory(true))
	if err != nil {
		log.Fatal(err)
	}

	ethRPC, err := ethclient.Dial("https://goerli.infura.io/v3/")
	if err != nil {
		log.Fatal(err)
	}

	return &TXAnalyzer{db: db, ethRPC: ethRPC}, nil
}

func (txa *TXAnalyzer) StoreBatches(batches []wrappers.GravityTransactionBatchExecutedEvent) error {
	err := txa.db.Update(func(txn *badger.Txn) error {
		for _, batch := range batches {

			err := txn.Set(rawTXKey(batch), []byte{})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (txa *TXAnalyzer) GetUnprocessedRawTXs(tokenAddress ethcmn.Address) ([][]byte, error) {
	unprocessedTxs := make([][]byte, 0)

	err := txa.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = append(KeyPrefixRawTx, tokenAddress.Bytes()...)

		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()

			if item.KeySize() != 53 {
				return errors.New("wrong key length")
			}

			if item.ValueSize() == 0 {
				// not processed
				unprocessedTxs = append(unprocessedTxs, k)
			}

			// err := item.Value(func(v []byte) error {
			// 	fmt.Printf("erc20=%s, txhash=%s, value=%s\n", ethcmn.Bytes2Hex(k[1:21]), ethcmn.Bytes2Hex(k[21:53]), v)
			// 	return nil
			// })
			// if err != nil {
			// 	return err
			// }
		}
		return nil
	})

	return unprocessedTxs, err
}

func (txa *TXAnalyzer) ProcessTXs(keys [][]byte) error {
	for _, k := range keys {
		// tokenAddr := ethcmn.BytesToAddress(k[1:21])
		txHash := ethcmn.BytesToHash(k[21:53])

		receipt, err := txa.ethRPC.TransactionReceipt(context.Background(), txHash)
		if err != nil {
			return err
		}

		txCount := uint8(len(receipt.Logs) - 2) // -2 removes 2 logs that are not outgoing txs
		gasUsed := sdk.Uint64ToBigEndian(receipt.GasUsed)

		err = txa.db.Update(func(txn *badger.Txn) error {
			err := txn.Set(k, append([]byte{txCount}, gasUsed...))
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func rawTXKey(evt wrappers.GravityTransactionBatchExecutedEvent) []byte {
	return append(KeyPrefixRawTx, append(evt.Token.Bytes(), evt.Raw.TxHash.Bytes()...)...)
}

// TODO: figure out where to call this
func (txa *TXAnalyzer) Close() {
	txa.db.Close()
}
