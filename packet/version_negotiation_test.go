package packet_test

import (
	"testing"

	"github.com/simia-tech/go-quic/packet"
	"github.com/stretchr/testify/assert"
)

func TestVersionNegotitation(t *testing.T) {
	testCases := []struct {
		name         string
		connectionID uint64
		versions     []uint32
		bytes        []byte
	}{
		{"One", 1, []uint32{1},
			[]byte{0x09, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00}},
		{"Two", 1, []uint32{1, 2},
			[]byte{0x09, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00}},
	}

	t.Run("Write", func(t *testing.T) {
		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				buffer := make([]byte, len(testCase.bytes))

				vnp := packet.VersionNegotiation(buffer)
				vnp.SetConnectionID(testCase.connectionID)
				vnp.SetVersions(testCase.versions)

				assert.Equal(t, len(testCase.bytes), vnp.Len())
				assert.Equal(t, testCase.bytes, buffer)
			})
		}
	})

	t.Run("Read", func(t *testing.T) {
		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				vnp := packet.VersionNegotiation(testCase.bytes)
				assert.Equal(t, testCase.connectionID, vnp.ConnectionID())
				assert.Equal(t, testCase.versions, vnp.Versions())
			})
		}
	})
}
