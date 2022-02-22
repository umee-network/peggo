package peggo

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/knadh/koanf"
	"github.com/pkg/errors"
)

var ethManager *EthRPCManager

type EthRPCManager struct {
	currentEndpoint int // the slice index of the endpoint currently used
	//client          *ethclient.Client // the current client
	client *rpc.Client
	konfig *koanf.Koanf

	// TODO: how to detect client failures so we can call DialNext()
	// maybe wrap methods
}

// initializes the single instance of EthRPCManager with a given config (uses flagEthRPCs)
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
