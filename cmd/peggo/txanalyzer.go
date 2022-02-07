package peggo

import (
	"context"
	"fmt"
	"os"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/umee-network/peggo/orchestrator/ethereum/provider"
	"github.com/umee-network/peggo/orchestrator/txanalyzer"
	"golang.org/x/sync/errgroup"
)

func getTxAnalyzerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "txanalyzer [gravity-addr]",
		Short: "Analyzes transactions on the Gravity Bridge contract to make cost estimates",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			konfig, err := parseServerConfig(cmd)
			if err != nil {
				return err
			}

			logger, err := getLogger(cmd)
			if err != nil {
				return err
			}

			ethRPCEndpoint := konfig.String(flagEthRPC)
			ethRPC, err := ethrpc.Dial(ethRPCEndpoint)
			if err != nil {
				return fmt.Errorf("failed to dial Ethereum RPC node: %w", err)
			}

			fmt.Fprintf(os.Stderr, "Connected to Ethereum RPC: %s\n", ethRPCEndpoint)
			ethProvider := provider.NewEVMProvider(ethRPC)

			gravityAddr := ethcmn.HexToAddress(args[0])

			txa, err := txanalyzer.NewTXAnalyzer(
				logger,
				"./txanalyzer",
				ethProvider,
				gravityAddr,
				172800,
			)

			if err != nil {
				return fmt.Errorf("failed to create TX Analyzer: %w", err)
			}

			ctx, cancel := context.WithCancel(context.Background())
			g, errCtx := errgroup.WithContext(ctx)

			g.Go(func() error {
				return startTXAnalyzer(errCtx, logger, txa)
			})

			// listen for and trap any OS signal to gracefully shutdown and exit
			trapSignal(cancel)

			return g.Wait()
		},
	}

	cmd.Flags().String(flagEthRPC, "http://localhost:8545", "Specify the RPC address of an Ethereum node")

	return cmd
}

func startTXAnalyzer(ctx context.Context, logger zerolog.Logger, txa *txanalyzer.TXAnalyzer) error {
	srvErrCh := make(chan error, 1)
	go func() {
		logger.Info().Msg("starting tx analyzer...")
		srvErrCh <- txa.Start(ctx)
	}()

	for {
		select {
		case <-ctx.Done():
			return nil

		case err := <-srvErrCh:
			logger.Error().Err(err).Msg("failed to start orchestrator")
			return err
		}
	}
}
