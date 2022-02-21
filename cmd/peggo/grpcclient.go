package peggo

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/knadh/koanf"
	"github.com/pkg/errors"
)

// dials configured ethereum rpc endpoints, returning a client connected to the first successfully dialed
func getEthClient(konfig *koanf.Koanf) (*ethclient.Client, error) {
	rpcs := konfig.Strings(flagEthRPCs)

	for _, endpoint := range rpcs {
		if cli, err := ethclient.Dial(endpoint); err == nil {
			return cli, nil
		}
	}

	// todo #196: this doesn't try to remember if an endpoint is failing frequently
	// which would be needed to rotate / avoid trying to dial it every single time.
	//
	// a cheap way to do basically that would be to store the last endpoint *successfully*
	// dialed using a var at the top of this file, always trying that first if non-empty,
	// clearing it on fail, and setting it to a new string whenever a dial succeeds on
	// an endpoint in the for loop
	//
	// also, if there are other reasons we want to avoid an endpoint (e.g. dialing succeeds
	// but the rpcs are behaving badly) then this function won't be able to help

	return nil, errors.New("Could not connect to any of the Ethereum RPC endpoints provided")
}
