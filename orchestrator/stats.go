package orchestrator

import (
	"encoding/json"
	"net/http"

	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravity/types"
	"github.com/rs/zerolog"
)

type StatKind int

const (
	Waiting StatKind = iota
	Unprocessed
	Successful
)

type Stats struct {
	Transactions map[StatKind][]types.BatchFees `json:"transactions"`
}

func (s StatKind) String() string {
	kinds := [...]string{"waiting", "unprocessed", "successful"}
	return kinds[s]
}

var s = Stats{Transactions: make(map[StatKind][]types.BatchFees)}

func init() {
	// testing stats
	s.Add(Waiting, types.BatchFees{})
	s.Add(Unprocessed, types.BatchFees{})
	s.Add(Successful, types.BatchFees{})
}

func (s Stats) Add(kind StatKind, batchFees types.BatchFees) (stats Stats) {
	// add kind of stat
	s.Transactions[kind] = append(s.Transactions[kind], batchFees)
	// TODO prune slice after N length (or MAX size-- better)
	return s
}

func Listen(port string, logger zerolog.Logger) {
	// Create http endpoint for peggo statistic requests
	http.HandleFunc("/", writeStats)
	logger.Info().Msgf("Starting status API on %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.Fatal().Msgf("Failed to start status API on %s: %s\n", port, err)
	}
}

// helper fns
// ---------
func writeStats(w http.ResponseWriter, r *http.Request) {
	// human-friendly, stringify keys
	transactions := make(map[string][]types.BatchFees)
	for kind, batchFees := range s.Transactions {
		transactions[kind.String()] = batchFees
	}
	// output json
	if data, err := json.Marshal(transactions); err == nil {
		w.Write(data)
	} else {
		httpError(w, "failed to marshal stats")
	}
}

func httpError(w http.ResponseWriter, err string) {
	w.Header().Set("x-peggo-error", err)
	w.Write([]byte(err))
	defer w.WriteHeader(http.StatusBadRequest)
}
