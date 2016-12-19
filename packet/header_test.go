package packet_test

import (
	"testing"

	"github.com/simia-tech/go-quic/packet"
	"github.com/stretchr/testify/assert"
)

func TestHeader(t *testing.T) {
	testCases := []struct {
		name         string
		flags        uint8
		connectionID interface{}

		bytes []byte
	}{
		{"Regular", packet.ConnectionID, uint64(1),
			[]byte{0x08, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{"Version", packet.VersionFlag | packet.ConnectionID, uint64(1),
			[]byte{0x09, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{"PublicReset", packet.PublicResetFlag | packet.ConnectionID, uint64(1),
			[]byte{0x0a, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	}

	t.Run("Write", func(t *testing.T) {
		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				buffer := make([]byte, len(testCase.bytes))

				header := packet.Header(buffer)
				header.SetFlags(testCase.flags)
				if testCase.connectionID != nil {
					header.AddConnectionID(testCase.connectionID.(uint64))
				}

				assert.Equal(t, len(testCase.bytes), header.Len())
				assert.Equal(t, testCase.bytes, buffer)
			})
		}
	})

	t.Run("Read", func(t *testing.T) {
		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				header := packet.Header(testCase.bytes)
				assert.Equal(t, testCase.flags, header.Flags())
				if testCase.connectionID != nil {
					assert.Equal(t, testCase.connectionID, header.ConnectionID())
				}
			})
		}
	})
}
