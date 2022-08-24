package orchestrator

import (
	"encoding/json"
	"testing"

	gravitytypes "github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravity/types"
)

// mocked payload, used for initial sanity test
const payload = `{"transactions":{"to_eth":{"successful":[],"unprocessed":[],"waiting":[]},"to_cosmos":{"successful":[],"unprocessed":[],"waiting":[]}}}`

func TestStats(t *testing.T) {
	type test struct {
		name     string
		expected string
	}

	var tests = []test{
		{name: "sane initial payload", expected: payload},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// compare the expected and actual for empty `Set`s
			s := NewStats()
			s.Set(ToEth, Waiting, []gravitytypes.BatchFees{})
			s.Set(ToEth, Successful, []gravitytypes.BatchFees{})
			s.Set(ToEth, Unprocessed, []gravitytypes.BatchFees{})
			s.Set(ToCosmos, Waiting, []gravitytypes.OutgoingTxBatch{})
			s.Set(ToCosmos, Successful, []gravitytypes.OutgoingTxBatch{})
			s.Set(ToCosmos, Unprocessed, []gravitytypes.OutgoingTxBatch{})
			x, _ := json.Marshal(s)
			got := string(x)
			if got != tc.expected {
				t.Errorf("expected: %s, got: %s", tc.expected, got)
			}

			// TODO run peggo

			// TODO compare the expected and actual for non-empty `Set`s

		})
	}
}
