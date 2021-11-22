package peggy

import (
	"bytes"
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

type PendingTxInput struct {
	InputData    hexutil.Bytes
	ReceivedTime time.Time
}

type PendingTxInputList []PendingTxInput

func (p *PendingTxInputList) AddPendingTxInput(pendingTx *RPCTransaction) {

	submitBatchMethod := peggyABI.Methods["submitBatch"]
	valsetUpdateMethod := peggyABI.Methods["updateValset"]

	// If it's not a submitBatch or updateValset transaction, ignore it
	if !bytes.Equal(submitBatchMethod.ID, pendingTx.Input[:4]) &&
		!bytes.Equal(valsetUpdateMethod.ID, pendingTx.Input[:4]) {
		return
	}

	pendingTxInput := PendingTxInput{
		InputData:    pendingTx.Input,
		ReceivedTime: time.Now(),
	}

	// Enqueue pending tx input
	*p = append(*p, pendingTxInput)
	// Persisting top 100 pending txs of peggy contract only.
	if len(*p) > 100 {
		(*p)[0] = PendingTxInput{} // to avoid memory leak
		// Dequeue pending tx input
		*p = (*p)[1:]
	}
}

// IsPendingTxInput returns true if the input data is found in the pending tx list. If the tx is found but the tx is
// older than pendingTxWaitDuration, we consider it stale and return false, so the validator re-sends it.
func (s *peggyContract) IsPendingTxInput(txData []byte, pendingTxWaitDuration time.Duration) bool {
	for _, pendingTxInput := range s.pendingTxInputList {
		if bytes.Equal(pendingTxInput.InputData, txData) {
			// If this tx was for too long in the pending list, consider it stale
			return time.Now().Before(pendingTxInput.ReceivedTime.Add(pendingTxWaitDuration))
		}
	}
	return false
}

func (s *peggyContract) SubscribeToPendingTxs(alchemyWebsocketURL string) {
	args := map[string]interface{}{
		"address": s.peggyAddress.Hex(),
	}

	wsClient, err := rpc.Dial(alchemyWebsocketURL)
	if err != nil {
		s.logger.Fatal().
			AnErr("err", err).
			Str("endpoint", alchemyWebsocketURL).
			Msg("Failed to connect to Alchemy Websocket")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Subscribe to Transactions
	ch := make(chan *RPCTransaction)
	_, err = wsClient.EthSubscribe(ctx, ch, "alchemy_filteredNewFullPendingTransactions", args)
	if err != nil {
		s.logger.Fatal().
			AnErr("err", err).
			Str("endpoint", alchemyWebsocketURL).
			Msg("Failed to subscribe to pending transactions")
		return
	}

	for {
		// Check that the transaction was send over the channel
		pendingTransaction := <-ch
		s.pendingTxInputList.AddPendingTxInput(pendingTransaction)
	}
}

// RPCTransaction represents a transaction that will serialize to the RPC representation of a transaction
type RPCTransaction struct {
	BlockHash        *common.Hash         `json:"blockHash"`
	BlockNumber      *hexutil.Big         `json:"blockNumber"`
	From             common.Address       `json:"from"`
	Gas              hexutil.Uint64       `json:"gas"`
	GasPrice         *hexutil.Big         `json:"gasPrice"`
	GasFeeCap        *hexutil.Big         `json:"maxFeePerGas,omitempty"`
	GasTipCap        *hexutil.Big         `json:"maxPriorityFeePerGas,omitempty"`
	Hash             common.Hash          `json:"hash"`
	Input            hexutil.Bytes        `json:"input"`
	Nonce            hexutil.Uint64       `json:"nonce"`
	To               *common.Address      `json:"to"`
	TransactionIndex *hexutil.Uint64      `json:"transactionIndex"`
	Value            *hexutil.Big         `json:"value"`
	Type             hexutil.Uint64       `json:"type"`
	Accesses         *ethTypes.AccessList `json:"accessList,omitempty"`
	ChainID          *hexutil.Big         `json:"chainId,omitempty"`
	V                *hexutil.Big         `json:"v"`
	R                *hexutil.Big         `json:"r"`
	S                *hexutil.Big         `json:"s"`
}
