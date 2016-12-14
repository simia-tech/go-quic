package packet_test

import (
	"testing"

	"github.com/simia-tech/go-quic/packet"
	"github.com/stretchr/testify/assert"
)

func TestPublicReset(t *testing.T) {
	testCases := []struct {
		name         string
		connectionID uint64
		bytes        []byte
	}{
		{"Regular", 1,
			[]byte{0x0e, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	}

	t.Run("Write", func(t *testing.T) {
		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				buffer := make([]byte, len(testCase.bytes))

				prp := packet.PublicReset(buffer)
				prp.SetConnectionID(testCase.connectionID)

				assert.Equal(t, len(testCase.bytes), prp.Len())
				assert.Equal(t, testCase.bytes, buffer)
			})
		}
	})

	t.Run("Read", func(t *testing.T) {
		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				prp := packet.PublicReset(testCase.bytes)
				assert.Equal(t, testCase.connectionID, prp.ConnectionID())
			})
		}
	})
}
