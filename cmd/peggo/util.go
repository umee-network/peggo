package peggo

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

func stdinConfirm(msg string) bool {
	var response string

	fmt.Fprint(os.Stderr, msg)

	if _, err := fmt.Scanln(&response); err != nil {
		fmt.Fprintf(os.Stderr, "failed to confirm action: %s", err)
		return false
	}

	switch strings.ToLower(strings.TrimSpace(response)) {
	case "y", "yes":
		return true

	default:
		return false
	}
}

// parseERC20ContractMapping converts list of address:denom pairs to a proper
// typed map.
func parseERC20ContractMapping(items []string) map[ethcmn.Address]string {
	res := make(map[ethcmn.Address]string)

	for _, item := range items {
		// item is a pair address:denom
		parts := strings.Split(item, ":")
		addr := ethcmn.HexToAddress(parts[0])

		if len(parts) != 2 || len(parts[0]) == 0 || addr == (ethcmn.Address{}) {
			fmt.Fprint(os.Stderr, "failed to parse ERC20 mapping: check that all inputs contain valid denom:address pairs")
			os.Exit(1)
		}

		denom := parts[1]
		res[addr] = denom
	}

	return res
}

// duration parses duration from string with a provided default fallback.
func duration(s string, defaults time.Duration) time.Duration {
	dur, err := time.ParseDuration(s)
	if err != nil {
		dur = defaults
	}

	return dur
}

func hexToBytes(str string) ([]byte, error) {
	if strings.HasPrefix(str, "0x") {
		str = str[2:]
	}

	data, err := hex.DecodeString(str)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// waitForService awaits an active connection to a gRPC service.
func waitForService(ctx context.Context, clientconn *grpc.ClientConn) {
	for {
		select {
		case <-ctx.Done():
			fmt.Fprint(os.Stderr, "gRPC service wait timed out")
			os.Exit(1)

		default:
			state := clientconn.GetState()

			if state != connectivity.Ready {
				// TODO: ...

				// fmt.Fprintf(os.Stderr, "")
				// log.WithField("state", state.String()).Warningln("state of gRPC connection not ready")
				time.Sleep(5 * time.Second)
				continue
			}

			return
		}
	}
}

func orShutdown(err error) {
	if err != nil && err != grpc.ErrServerStopped {
		fmt.Fprint(os.Stderr, "unable to start peggo orchestrator")
		os.Exit(1)
	}
}
