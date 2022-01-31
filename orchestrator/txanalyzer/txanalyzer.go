package txanalyzer

import (
	"context"
	"errors"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	badger "github.com/dgraph-io/badger/v3"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"

	"github.com/umee-network/peggo/orchestrator/ethereum/provider"
	wrappers "github.com/umee-network/peggo/solwrappers/Gravity.sol"
)

// erc20DefaultValues is going to be used to fill any gaps in data until we figure out a better way.
var erc20DefaultValues = []uint64{575563, 582863, 591565, 600967, 611968, 621532, 630328, 642386, 653063, 661581, 668183, 678635, 685289, 696851, 704866, 712721, 708887, 721445, 727461, 734690, 742043, 752750, 767272, 760223, 769101, 784019, 773423, 798268, 806362, 802351, 807763, 831213, 828969, 814683, 843207, 875285, 847569, 877254, 870002, 873950, 882126, 887008, 919109, 911901, 927237, 911510, 936261, 947621, 918882, 935638, 920685, 933757, 948716, 976337, 970508, 965708, 998407, 995011, 999148, 1024643, 1016724, 1046768, 1053295, 1044177, 1035313, 1059293, 1084901, 1094332, 1103762, 1053903, 1114666, 1078123, 1073982, 1078022, 1082061, 1157889, 1108249, 1146072, 1126675, 1173955, 1136556, 1159855, 1191781, 1234944, 1188274, 1181863, 1171010, 1172318, 1226168, 1154187, 1194480, 1238819, 1227017, 1209858, 1256137, 1261934, 1228247, 1261745, 1258859, 1244511}

// KVStore key prefixes
var (
	KeyPrefixUnprocessedTx = []byte{0x01}
	KeyPrefixProcessedTx   = []byte{0x02}
	KeyPrefixEstimate      = []byte{0x03}
)

type TXAnalyzer struct {
	logger          zerolog.Logger
	db              *badger.DB
	evmProvider     provider.EVMProviderWithRet
	pruneKeepRecent uint64
}

func NewTXAnalyzer(
	logger zerolog.Logger,
	dbDir string,
	evmProvider provider.EVMProviderWithRet,
	pruneKeepRecent uint64,
) (*TXAnalyzer, error) {
	db, err := badger.Open(badger.DefaultOptions(dbDir).WithInMemory(true))
	if err != nil {
		logger.Fatal().AnErr("err", err).Msg("failed to open db for txanalyzer")
	}

	return &TXAnalyzer{
		logger:          logger.With().Str("module", "txanalyzer").Logger(),
		db:              db,
		evmProvider:     evmProvider,
		pruneKeepRecent: pruneKeepRecent,
	}, nil
}

func (txa *TXAnalyzer) StoreBatches(batches []wrappers.GravityTransactionBatchExecutedEvent) error {
	err := txa.db.Update(func(txn *badger.Txn) error {
		for _, batch := range batches {
			err := txn.Set(unprocessedTxKey(batch.Raw.TxHash), batch.Token.Bytes())
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (txa *TXAnalyzer) GetUnprocessedTXsByToken() (map[ethcmn.Address][]ethcmn.Hash, error) {
	unprocessedTxs := map[ethcmn.Address][]ethcmn.Hash{}

	err := txa.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = KeyPrefixUnprocessedTx

		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()

			if item.KeySize() != 33 {
				return errors.New("wrong key length")
			}

			err := item.Value(func(v []byte) error {
				unprocessedTxs[ethcmn.BytesToAddress(v)] = append(unprocessedTxs[ethcmn.BytesToAddress(v)], ethcmn.BytesToHash(k[1:]))
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	return unprocessedTxs, err
}

func (txa *TXAnalyzer) ProcessTXs(txs map[ethcmn.Address][]ethcmn.Hash) error {
	for tokenAddr, txHashes := range txs {
		for _, txHash := range txHashes {
			receipt, err := txa.evmProvider.TransactionReceipt(context.Background(), txHash)
			if err != nil {
				return err
			}

			k := processedTxKey(receipt.BlockNumber, tokenAddr, txHash)

			txCount := uint8(len(receipt.Logs) - 2) // -2 removes 2 logs that are not outgoing txs
			gasUsed := sdk.Uint64ToBigEndian(receipt.GasUsed)

			value := []byte{txCount}
			value = append(value, gasUsed...)

			// Store the tx's data and delete the unprocessed marker
			err = txa.db.Update(func(txn *badger.Txn) error {
				err := txn.Set(k, value)
				if err != nil {
					return err
				}

				err = txn.Delete(unprocessedTxKey(txHash))
				return err
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (txa *TXAnalyzer) PruneTXs() (int, error) {
	count := 0

	lastBlock, err := txa.evmProvider.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return 0, err
	}

	minimumBlock := lastBlock.Number.Uint64() - txa.pruneKeepRecent

	err = txa.db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = KeyPrefixProcessedTx

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			key := it.Item().Key()
			blockNumber := sdk.BigEndianToUint64(key[1:9])

			if blockNumber > minimumBlock {
				break
			}

			err := txn.Delete(key)
			if err != nil {
				return err
			}
			count++

		}
		return nil
	})

	return count, err
}

func (txa *TXAnalyzer) RecalculateEstimates() error {
	// 1. Get all processed txs and sum them up in a map
	// 2. Go through the totals and calculate the average gas used per batch count
	// 3. Fill gaps somehow?

	totals := map[ethcmn.Address][][]uint64{}

	// map[tokenAddr][][]uint64{ {totalBatchCount, totalGasUsed}  }

	err := txa.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = KeyPrefixProcessedTx

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			key := it.Item().Key()
			tokenAddr := ethcmn.BytesToAddress(key[9:29])

			if _, ok := totals[tokenAddr]; !ok {
				// initialize averages slices
				totals[tokenAddr] = make([][]uint64, 100)
				for i := range totals[tokenAddr] {
					totals[tokenAddr][i] = make([]uint64, 2)
				}
			}

			outTxCount := uint64(0)
			gasUsed := uint64(0)
			err := it.Item().Value(func(v []byte) error {
				outTxCount = uint64(v[0])
				gasUsed = sdk.BigEndianToUint64(v[1:])
				return nil
			})
			if err != nil {
				return err
			}

			totals[tokenAddr][outTxCount-1][0] += 1
			totals[tokenAddr][outTxCount-1][1] += gasUsed
		}
		return nil
	})

	// 2. calculate the averages
	// 3. fill the gaps
	for tokenAddr, v := range totals {
		for i := range v {
			if v[i][0] == 0 {
				// no data points here, so we fill the gap with our hardcoded data
				totals[tokenAddr][i][1] = erc20DefaultValues[i]
			} else {
				// average gas used
				totals[tokenAddr][i][1] = v[i][1] / v[i][0]
			}
		}
	}
	// We can use the tx count as an accuracy indicator for the gas estimate, if it's 0, then it might be off by a lot.
	// If it's +20 then it means we might have a pretty accurate estimate as long as we are keeping only recent data.

	// 4. Store the averages
	txa.db.Update(func(txn *badger.Txn) error {
		for tokenAddr, v := range totals {
			value := []byte{}

			for i := range v {
				value = append(value, sdk.Uint64ToBigEndian(v[i][0])...)
				value = append(value, sdk.Uint64ToBigEndian(v[i][1])...)
			}

			err := txn.Set(estimateKey(tokenAddr), value)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

func (txa *TXAnalyzer) GetEstimatesOfToken(tokenAddr ethcmn.Address) ([][]uint64, error) {

	estimates := make([][]uint64, 100)
	for i := range estimates {
		estimates[i] = make([]uint64, 2)
	}

	err := txa.db.View(func(txn *badger.Txn) error {
		it, err := txn.Get(estimateKey(tokenAddr))
		if err != nil {
			return err
		}

		err = it.Value(func(v []byte) error {
			for i := 0; i < len(v); i += 16 {
				estimates[i/16][0] = sdk.BigEndianToUint64(v[i : i+8])
				estimates[i/16][1] = sdk.BigEndianToUint64(v[i+8 : i+16])
			}
			return nil
		})

		return err

	})

	return estimates, err
}

// TODO: figure out where to call this
func (txa *TXAnalyzer) Close() error {
	return txa.db.Close()
}

func unprocessedTxKey(txHash ethcmn.Hash) []byte {
	return append(KeyPrefixUnprocessedTx, txHash.Bytes()...)
}
func processedTxKey(blockNumber *big.Int, tokenAddr ethcmn.Address, txHash ethcmn.Hash) []byte {
	key := append(KeyPrefixProcessedTx, sdk.Uint64ToBigEndian(blockNumber.Uint64())...)
	key = append(key, tokenAddr.Bytes()...)
	key = append(key, txHash.Bytes()...)
	return key
}

func estimateKey(tokenAddr ethcmn.Address) []byte {
	return append(KeyPrefixEstimate, tokenAddr.Bytes()...)
}
