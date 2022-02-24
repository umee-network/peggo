package peggo

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/knadh/koanf"
	"github.com/pkg/errors"
)

var ethManager *EthRPCManager

type EthRPCManager struct {
	currentEndpoint int // the slice index of the endpoint currently used
	client          *rpc.Client
	konfig          *koanf.Koanf
}

// initializes the single instance of EthRPCManager with a given config (uses flagEthRPCs).
// no-op if already initialized, even if konfig would be different.
func InitEthRPCManager(konfig *koanf.Koanf) {
	if ethManager == nil {
		ethManager = &EthRPCManager{
			konfig: konfig,
		}
	}
}

// closes and sets to nil the stored eth RPC client
func (em *EthRPCManager) CloseClient() {
	if em.client != nil {
		em.client.Close()
		em.client = nil
	}
}

// closes the current client and dials configured ethereum rpc endpoints in a roundrobin fashion until one
// is connected. returns an error if no endpoints ar configured or all dials failed
func (em *EthRPCManager) DialNext() error {
	if em.konfig == nil {
		return errors.New("ethRPCManager konfig is nil")
	}
	rpcs := em.konfig.Strings(flagEthRPCs)

	em.CloseClient()

	dialIndex := func(i int) bool {
		if cli, err := rpc.Dial(rpcs[i]); err == nil {
			em.currentEndpoint = i
			em.client = cli
			return true
		}
		// todo: should likely log the error
		return false
	}

	// first tries all endpoints in the slice after the current index
	for i := range rpcs {
		if i > em.currentEndpoint && dialIndex(i) {
			return nil
		}
	}

	// then tries remaining endpoints from the beginning of the slice
	for i := range rpcs {
		if i <= em.currentEndpoint && dialIndex(i) {
			return nil
		}
	}

	return errors.New(fmt.Sprintf("could not dial any of the %d Ethereum RPC endpoints configured", len(rpcs)))
}

// returns the current eth RPC client, dialing one first if nonexistent
func (em *EthRPCManager) GetClient() (*rpc.Client, error) {
	if em.client == nil {
		if err := em.DialNext(); err != nil {
			return nil, err
		}
	}
	return em.client, nil
}

// returns the current eth RPC client, dialing one first if nonexistent
func (em *EthRPCManager) GetEthClient() (*ethclient.Client, error) {
	cli, err := em.GetClient()
	if err != nil {
		return nil, err
	}
	return ethclient.NewClient(cli), nil
}

// wraps ethclient.PendingNonceAt, also closing client if PendingNonceAt returns an error
func (em *EthRPCManager) PendingNonceAt(ctx context.Context, addr common.Address) (uint64, error) {
	cli, err := em.GetEthClient()
	if err != nil {
		return 0, err
	}
	nonce, err := cli.PendingNonceAt(ctx, addr)
	if err != nil {
		em.CloseClient()
		return 0, err
	}
	return nonce, nil
}

// wraps ethclient.ChainID, also closing client if ChainID returns an error
func (em *EthRPCManager) ChainID(ctx context.Context) (*big.Int, error) {
	cli, err := em.GetEthClient()
	if err != nil {
		return nil, err
	}
	id, err := cli.ChainID(ctx)
	if err != nil {
		em.CloseClient()
		return nil, err
	}
	return id, nil
}

// wraps ethclient.SuggestGasPrice, also closing client if SuggestGasPrice returns an error
func (em *EthRPCManager) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	cli, err := em.GetEthClient()
	if err != nil {
		return nil, err
	}
	price, err := cli.SuggestGasPrice(ctx)
	if err != nil {
		em.CloseClient()
		return nil, err
	}
	return price, nil
}
