package peggy

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPeggyPowerToPercent(t *testing.T) {
	percent := peggyPowerToPercent(big.NewInt(213192100))
	assert.Equal(t, percent, float32(4.9637656))

}
