package orchestrator

import (
	"encoding/json"
	"fmt"
	"testing"

	gravitytypes "github.com/Gravity-Bridge/Gravity-Bridge/module/x/gravity/types"
)

// TODO match a mocked payload to a testcase
const payload = `{
	"transactions": {
		"to_eth": {
			"waiting": []
		},
	},
}`

func TestStats(t *testing.T) {
	type test struct {
		name     string
		input    string
		expected string
	}

	var tests = []test{
		{name: "good payload"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := NewStats()
			s.Set(ToEth, Waiting, []gravitytypes.BatchFees{})
			s.Set(ToEth, Successful, []gravitytypes.BatchFees{})
			s.Set(ToEth, Unprocessed, []gravitytypes.BatchFees{})
			s.Set(ToCosmos, Waiting, []gravitytypes.OutgoingTxBatch{})
			s.Set(ToCosmos, Successful, []gravitytypes.OutgoingTxBatch{})
			s.Set(ToCosmos, Unprocessed, []gravitytypes.OutgoingTxBatch{})
			x, _ := json.Marshal(s)
			fmt.Printf("%s\n", string(x))
		})
	}

}
