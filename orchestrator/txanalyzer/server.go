package txanalyzer

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
)

func (txa *TXAnalyzer) serveEstimates(ctx context.Context) error {
	r := mux.NewRouter()
	r.HandleFunc("/estimates/{token_address}", func(w http.ResponseWriter, r *http.Request) {
		address := ethcmn.HexToAddress(mux.Vars(r)["token_address"])

		w.Header().Set("Content-Type", "application/json")

		vals, err := txa.GetEstimatesOfToken(address)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "` + err.Error() + `"}`))
			txa.logger.Err(err).Msg("failed to get estimates")
			return
		}

		response, _ := json.Marshal(vals)

		w.WriteHeader(http.StatusOK)
		w.Write(response)

	})

	writeTimeout, err := time.ParseDuration("15s")
	if err != nil {
		return err
	}
	readTimeout, err := time.ParseDuration("15s")
	if err != nil {
		return err
	}

	srvErrCh := make(chan error, 1)
	srv := &http.Server{
		Handler:      r,
		Addr:         txa.listenAddr,
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
	}

	go func() {
		txa.logger.Info().Str("listen_addr", txa.listenAddr).Msg("starting txanalyzer server...")
		srvErrCh <- srv.ListenAndServe()
	}()

	for {
		select {
		case <-ctx.Done():
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			txa.logger.Info().Str("listen_addr", txa.listenAddr).Msg("shutting down tx analyzer server...")
			if err := srv.Shutdown(shutdownCtx); err != nil {
				txa.logger.Error().Err(err).Msg("failed to gracefully shutdown tx analyzer server")
				return err
			}

			return nil

		case err := <-srvErrCh:
			txa.logger.Error().Err(err).Msg("failed to start tx analyzer server")
			return err
		}
	}
}
