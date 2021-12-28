package relayer

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/umee-network/umee/x/peggy/types"
)

func TestCheckIfValsetsDiffer(t *testing.T) {
	// this function doesn't return a value. Running different scenarios just to increase code coverage.

	t.Run("ok", func(t *testing.T) {
		logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr})

		relayer := peggyRelayer{
			logger: logger,
		}

		relayer.checkIfValsetsDiffer(&types.Valset{}, &types.Valset{})
		relayer.checkIfValsetsDiffer(nil, &types.Valset{})
		relayer.checkIfValsetsDiffer(nil, &types.Valset{Nonce: 2})
		relayer.checkIfValsetsDiffer(&types.Valset{Nonce: 12}, &types.Valset{Nonce: 11})
		relayer.checkIfValsetsDiffer(&types.Valset{}, &types.Valset{Members: []*types.BridgeValidator{{EthereumAddress: "0x0"}}})
	})

}
