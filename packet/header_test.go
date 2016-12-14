package packet_test

import (
	"testing"

	"github.com/simia-tech/go-quic/packet"
	"github.com/stretchr/testify/assert"
)

func TestHeader(t *testing.T) {
	testCases := []struct {
		name          string
		flags         uint8
		connectionID  interface{}
		expectedLen   int
		expectedBytes []byte
	}{
		{"Version", packet.PublicVersionFlag, uint64(1),
			0, []byte{0x0d, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			buffer := [19]byte{}

			header := packet.Header(buffer[:])
			header.SetPublic(testCase.flags)
			header.SetConnectionID(testCase.connectionID)

			assert.Equal(t, testCase.expectedLen, header.Len())
			assert.Equal(t, testCase.expectedBytes, buffer[:])
		})
	}
}
