package peggo

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	"github.com/umee-network/peggo/cmd/peggo/client"
	"github.com/umee-network/peggo/orchestrator/cosmos"
	peggytypes "github.com/umee-network/umee/x/peggy/types"
)

func getTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions for Peggy governance and maintenance on the Cosmos chain",
		Long: `Transactions for Peggy governance and maintenance on the Cosmos chain.

Inputs in the CLI commands can be provided via flags or environment variables. If
using the later, prefix the environment variable with PEGGO_ and the named of the
flag (e.g. PEGGO_COSMOS_PK).`,
	}

	cmd.AddCommand(
		getRegisterEthKeyCmd(),
	)

	return cmd
}

func getRegisterEthKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-eth-key",
		Short: "Submits an Ethereum key that will be used to sign messages on behalf of your Validator",
		RunE: func(cmd *cobra.Command, args []string) error {
			konfig, err := parseServerConfig(cmd)
			if err != nil {
				return err
			}

			if konfig.Bool(flagEthUseLedger) {
				fmt.Fprintf(os.Stderr, "WARNING: Uou cannot use Ledger for orchestrator, so make sure the Ethereum key is accessible outside of it")
			}

			valAddress, cosmosKeyring, err := initCosmosKeyring(konfig)
			if err != nil {
				return fmt.Errorf("failed to initialize Cosmos keyring")
			}

			ethKeyFromAddress, _, personalSignFn, err := initEthereumAccountsManager(0, konfig)
			if err != nil {
				return err
			}

			fmt.Fprintf(os.Stderr, "Using Cosmos validator address: %s", valAddress)
			fmt.Fprintf(os.Stderr, "Using Ethereum address: %s", ethKeyFromAddress)

			autoConfirm := konfig.Bool(flagAutoConfirm)
			actionConfirmed := autoConfirm || stdinConfirm("Confirm UpdatePeggyOrchestratorAddresses transaction? [y/N]: ")
			if !actionConfirmed {
				return nil
			}

			cosmosChainID := konfig.String(flagCosmosChainID)

			clientCtx, err := client.NewClientContext(cosmosChainID, valAddress.String(), cosmosKeyring)
			if err != nil {
				return err
			}

			tmRPCEndpoint := konfig.String(flagTendermintRPC)
			cosmosGRPC := konfig.String(flagCosmosGRPC)
			cosmosGasPrices := konfig.String(flagCosmosGasPrices)

			tmRPC, err := rpchttp.New(tmRPCEndpoint, "/websocket")
			if err != nil {
				return fmt.Errorf("failed to create Tendermint RPC client: %w", err)
			}

			clientCtx = clientCtx.WithClient(tmRPC).WithNodeURI(tmRPCEndpoint)

			daemonClient, err := client.NewCosmosClient(clientCtx, cosmosGRPC, client.OptionGasPrices(cosmosGasPrices))
			if err != nil {
				return err
			}

			// TODO: Clean this up to be more ergonomic and clean. We can probably
			// encapsulate all of this into a single utility function that gracefully
			// checks for the gRPC status/health.
			//
			// Ref: https://github.com/umee-network/peggo/issues/2

			fmt.Fprint(os.Stderr, "Waiting for cosmos gRPC service...")
			time.Sleep(time.Second)

			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			grpcConn := daemonClient.QueryClient()
			waitForService(ctx, grpcConn)

			peggyQuerier := peggytypes.NewQueryClient(grpcConn)
			peggyBroadcaster := cosmos.NewPeggyBroadcastClient(peggyQuerier, daemonClient, nil, personalSignFn)

			ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			if err = peggyBroadcaster.UpdatePeggyOrchestratorAddresses(ctx, ethKeyFromAddress, valAddress); err != nil {
				return fmt.Errorf("failed to broadcast transaction: %w", err)
			}

			fmt.Fprintf(os.Stderr, "Registered Ethereum Address %s for validator %s", ethKeyFromAddress, valAddress)
			return nil
		},
	}

	cmd.Flags().BoolP(flagAutoConfirm, "y", false, "Auto-confirm actions (e.g. transaction sending)")
	cmd.Flags().AddFlagSet(cosmosFlagSet())
	cmd.Flags().AddFlagSet(cosmosKeyringFlagSet())
	cmd.Flags().AddFlagSet(ethereumKeyOptsFlagSet())

	return cmd
}
