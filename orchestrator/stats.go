package orchestrator

import (
	"net/http"

	"github.com/rs/zerolog"
)

func Listen(port string, logger zerolog.Logger) {
	// Create http endpoint for peggo statistic requests
	http.HandleFunc("/", writeStats)
	logger.Info().Msgf("Starting status API on %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.Fatal().Msgf("Failed to start status API on %s: %s\n", port, err)
	}
}

func writeStats(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if !query.Has("test") {
		httpError(w, "test is required")
		return
	}
}

func httpError(w http.ResponseWriter, err string) {
	w.Header().Set("x-peggo-error", err)
	w.Write([]byte(err))
	defer w.WriteHeader(http.StatusBadRequest)
}
