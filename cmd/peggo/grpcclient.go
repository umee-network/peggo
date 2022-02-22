package peggo

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/knadh/koanf"
	"github.com/pkg/errors"
)

var ethManager *EthRPCManager

type EthRPCManager struct {
	currentEndpoint int               // the slice index of the endpoint currently used
	client          *ethclient.Client // the current client
	konfig          *koanf.Koanf

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

// closes the current client and dials configured ethereum rpc endpoints in a roundrobin fashion until one
// is connected. returns an error if no endpoints exist or all dials failed
func (em *EthRPCManager) DialNext() error {
	rpcs := em.konfig.Strings(flagEthRPCs)

	if em.client != nil {
		em.client.Close()
		em.client = nil
	}

	dialIndex := func(i int) bool {
		if cli, err := ethclient.Dial(rpcs[i]); err == nil {
			em.currentEndpoint = i
			em.client = cli
			return true
		}
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

	return errors.New("could not dial any of the Ethereum RPC endpoints provided")
}

func (em *EthRPCManager) GetClient() (*ethclient.Client, error) {
	if em.client == nil {
		if err := em.DialNext(); err != nil {
			return nil, err
		}
	}
	return em.client, nil
}
