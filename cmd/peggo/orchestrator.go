package peggo

import (
	"github.com/spf13/cobra"
)

func getOrchestratorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "orchestrator",
		Short: "Starts the orchestrator",
		Long: `Starts the orchestrator's main relayer loop. Only start the orchestrator
if the Peggy (Gravity Bridge) contract has been deployed and initialized, which
requires all validators to set their delegate keys.

Inputs in the CLI commands can be provided via flags or environment variables. If
using the later, prefix the environment variable with PEGGO_ and the named of the
flag (e.g. PEGGO_COSMOS_PK).`,
		RunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
	}

	cmd.Flags().Bool(flagRelayValsets, false, "Relay validator set updates to Ethereum")
	cmd.Flags().Bool(flagRelayBatches, false, "Relay transaction batches to Ethereum")
	cmd.Flags().AddFlagSet(cosmosFlagSet())
	cmd.Flags().AddFlagSet(cosmosKeyringFlagSet())
	cmd.Flags().AddFlagSet(ethereumKeyOptsFlagSet())
	cmd.Flags().AddFlagSet(ethereumOptsFlagSet())

	// TODO: Support profitable batch arguments. Injective's Peggo provides a
	// '--min-batch-fee-usd' argument which uses coingecko to retrieve price
	// information.
	//
	// Ref: https://github.com/umee-network/peggo/issues/2

	return cmd
}
