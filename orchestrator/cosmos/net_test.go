package cosmos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProtocolAndAddress(t *testing.T) {
	protocol, address := ProtocolAndAddress("tcp://127.0.0.1:8080")
	assert.Equal(t, "tcp", protocol)
	assert.Equal(t, "127.0.0.1:8080", address)
}

func TestConnect(t *testing.T) {
	conn, err := Connect("tcp://localhost:9999")
	assert.Nil(t, conn)
	assert.NotNil(t, err)
}
