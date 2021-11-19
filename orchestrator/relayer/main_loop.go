package relayer

import (
	"context"

	"errors"

	retry "github.com/avast/retry-go"

	"github.com/umee-network/peggo/orchestrator/loops"
)

func (s *peggyRelayer) Start(ctx context.Context) error {
	logger := s.logger.With().Str("loop", "RelayerMainLoop").Logger()

	// TODO: for now we multiply the blocktime by 3 to allow a bit more time to pass so we sent duplicated txs less
	// often
	return loops.RunLoop(ctx, s.logger, s.ethereumBlockTime*3, func() error {
		var pg loops.ParanoidGroup
		if s.valsetRelayEnabled {
			logger.Info().Msg("valset relay enabled; starting to relay valsets to Ethereum")
			pg.Go(func() error {
				return retry.Do(func() error {
					return s.RelayValsets(ctx)
				}, retry.Context(ctx), retry.OnRetry(func(n uint, err error) {
					logger.Err(err).Uint("retry", n).Msg("failed to relay valsets; retrying...")
				}))
			})
		}

		if s.batchRelayEnabled {
			logger.Info().Msg("batch relay enabled; starting to relay batches to Ethereum")
			pg.Go(func() error {
				return retry.Do(func() error {

					currentValset, err := s.FindLatestValset(ctx)
					if err != nil {
						return errors.New("failed to find latest valset")
					} else if currentValset == nil {
						return errors.New("latest valset not found")
					}

					possibleBatches, err := s.getBatchesAndSignatures(ctx, currentValset)
					if err != nil {
						return err
					}

					return s.RelayBatches(ctx, currentValset, possibleBatches)
				}, retry.Context(ctx), retry.OnRetry(func(n uint, err error) {
					logger.Err(err).Uint("retry", n).Msg("failed to relay tx batches; retrying...")
				}))
			})
		}

		if pg.Initialized() {
			if err := pg.Wait(); err != nil {
				logger.Err(err).Msg("main relay loop failed; exiting...")
				return err
			}
		}
		return nil
	})
}
