package packet_test

import (
	"testing"

	"github.com/simia-tech/go-quic/packet"
	"github.com/stretchr/testify/assert"
)

func TestHeaderWrite(t *testing.T) {
	testCases := []struct {
		name          string
		special       bool
		flags         uint8
		connectionID  interface{}
		packetNumber  interface{}
		expectedBytes []byte
	}{
		{"Version", true, packet.VersionFlag, uint64(1), nil,
			[]byte{0x0d, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{"PublicReset", true, packet.PublicResetFlag, uint64(1), nil,
			[]byte{0x0e, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{"Regular", false, 0, uint64(1), uint64(2),
			[]byte{0x3c, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			buffer := [packet.MaxHeaderSize]byte{}

			header := packet.Header(buffer[:])
			header.SetFlag(testCase.flags)
			header.SetConnectionID(testCase.connectionID)
			if testCase.packetNumber != nil {
				header.SetPacketNumber(testCase.packetNumber)
			}
			l := header.Len(testCase.special)

			assert.Equal(t, len(testCase.expectedBytes), l)
			assert.Equal(t, testCase.expectedBytes, buffer[:l])
		})
	}
}
