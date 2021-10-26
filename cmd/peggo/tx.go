package peggo

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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

			// ethKeyFromAddress, _, personalSignFn, err := initEthereumAccountsManager(
			// 	0,
			// 	ethKeystoreDir,
			// 	ethKeyFrom,
			// 	ethPassphrase,
			// 	ethPrivKey,
			// 	ethUseLedger,
			// )
			// if err != nil {
			// 	log.WithError(err).Fatalln("failed to init Ethereum account")
			// }

			// log.Infoln("Using Cosmos ValAddress", valAddress.String())
			// log.Infoln("Using Ethereum address", ethKeyFromAddress.String())

			// actionConfirmed := *alwaysAutoConfirm || stdinConfirm("Confirm UpdatePeggyOrchestratorAddresses transaction? [y/N]: ")
			// if !actionConfirmed {
			// 	return
			// }

			// clientCtx, err := client.NewClientContext(*cosmosChainID, valAddress.String(), cosmosKeyring)
			// if err != nil {
			// 	log.WithError(err).Fatalln("failed to initialize cosmos client context")
			// }
			// clientCtx = clientCtx.WithNodeURI(*tendermintRPC)

			// tmRPC, err := rpchttp.New(*tendermintRPC, "/websocket")
			// if err != nil {
			// 	log.WithError(err)
			// }

			// clientCtx = clientCtx.WithClient(tmRPC)
			// daemonClient, err := client.NewCosmosClient(clientCtx, *cosmosGRPC, client.OptionGasPrices(*cosmosGasPrices))
			// if err != nil {
			// 	log.WithError(err).WithFields(log.Fields{
			// 		"endpoint": *cosmosGRPC,
			// 	}).Fatalln("failed to connect to Cosmos daemon")
			// }

			// log.Infoln("Waiting for injectived GRPC")
			// time.Sleep(1 * time.Second)

			// daemonWaitCtx, cancelWait := context.WithTimeout(context.Background(), time.Minute)
			// grpcConn := daemonClient.QueryClient()
			// waitForService(daemonWaitCtx, grpcConn)
			// peggyQuerier := types.NewQueryClient(grpcConn)
			// peggyBroadcaster := cosmos.NewPeggyBroadcastClient(
			// 	peggyQuerier,
			// 	daemonClient,
			// 	nil,
			// 	personalSignFn,
			// )
			// cancelWait()

			// broadcastCtx, cancelFn := context.WithTimeout(context.Background(), 15*time.Second)
			// defer cancelFn()

			// if err = peggyBroadcaster.UpdatePeggyOrchestratorAddresses(broadcastCtx, ethKeyFromAddress, valAddress); err != nil {
			// 	log.WithError(err).Errorln("failed to broadcast Tx")
			// 	time.Sleep(time.Second)
			// 	return
			// }

			// log.Infof("Registered Ethereum address %s for validator address %s",
			// 	ethKeyFromAddress, valAddress.String())

			return nil
		},
	}

	cmd.Flags().BoolP(flagAutoConfirm, "y", false, "Always auto-confirm actions, such as transaction sending")
	cmd.Flags().AddFlagSet(cosmosFlagSet())
	cmd.Flags().AddFlagSet(cosmosKeyringFlagSet())
	cmd.Flags().AddFlagSet(ethereumKeyOptsFlagSet())

	return cmd
}
