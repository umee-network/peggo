package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravity/types"
	"github.com/rs/zerolog"
)

type (
	// TransactionData interface{fmt.Stringer}
	Transactions struct {
		ToEth    map[StatKind][]types.BatchFees       `json:"to_eth"`
		ToCosmos map[StatKind][]types.OutgoingTxBatch `json:"to_cosmos"`
		// ToEth    map[StatKind][]fmt.Stringer `json:"to_cosmos"`
		// ToCosmos map[StatKind][]fmt.Stringer `json:"to_cosmos"`
	}
	Stats struct {
		Transactions `json:"transactions"`
	}
)

const (
	ToEth StatDst = iota
	ToCosmos
)

const (
	Waiting StatKind = iota
	Unprocessed
	Successful
)

// StatsDst is the destination of transaction batches
type StatDst int

func (s StatDst) String() string {
	kinds := [...]string{"to_eth", "to_cosmos"}
	return kinds[s]
}

func (s *StatDst) UnmarshalText(b []byte) error {
	// TODO lookup by string -> int
	return nil
}

func (s StatDst) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

// StatKind is the stage these peggo transaction batches are in
type StatKind int

func (s StatKind) String() string {
	kinds := [...]string{"waiting", "unprocessed", "successful"}
	return kinds[s]
}

func (s *StatKind) UnmarshalText(b []byte) error {
	// TODO lookup by string -> int
	return nil
}

func (s StatKind) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s Stats) Debug(input interface{}) {
	fmt.Printf("%+v\n", input)
}

func NewStats() *Stats {
	return &Stats{
		Transactions: Transactions{
			ToEth:    make(map[StatKind][]types.BatchFees),
			ToCosmos: make(map[StatKind][]types.OutgoingTxBatch),
			// ToCosmos: make(map[StatKind][]fmt.Stringer),
		},
	}
}

func (s *Stats) Add(dst StatDst, kind StatKind, batchFees types.BatchFees) (stats *Stats) {
	// TODO add kind of stat
	// TODO prune slice after N length (or MAX size-- better)
	return s
}

func (s *Stats) Set(dst StatDst, kind StatKind, batchFees interface{}) {
	// set kind of stat
	switch dst {
	case ToEth:
		s.ToEth[kind] = batchFees.([]types.BatchFees)
	case ToCosmos:
		s.ToCosmos[kind] = batchFees.([]types.OutgoingTxBatch)
	}
}

func Listen(port string, ctx context.Context, logger zerolog.Logger) {
	// Create http endpoint for peggo statistic requests
	server := http.Server{Addr: ":" + port}
	http.HandleFunc("/", writeStats)
	logger.Info().Msgf("Starting status API on %s\n", port)
	go func() {
		for {
			select {
			case <-ctx.Done():
				server.Shutdown(ctx)
				return
			default:
				if err := server.ListenAndServe(); err != nil {
					logger.Fatal().Msgf("Failed to start status API on %s: %s\n", port, err)
				}
			}
		}
	}()
}

// helper fns
// ---------
func writeStats(w http.ResponseWriter, r *http.Request) {
	// human-friendly, stringify keys
	transactions := make(map[string]map[string][]types.BatchFees)
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
