package packet_test

import (
	"testing"

	"github.com/simia-tech/go-quic/packet"
	"github.com/stretchr/testify/assert"
)

func TestHeader(t *testing.T) {
	testCases := []struct {
		name         string
		special      bool
		flags        uint8
		connectionID interface{}
		version      interface{}
		packetNumber interface{}

		bytes []byte
	}{
		{"VersionClient", false, 0, uint64(1), uint32(2), uint64(3),
			[]byte{0x3d, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{"VersionServer", true, 0, uint64(1), nil, nil,
			[]byte{0x0c, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{"PublicReset", true, packet.PublicResetFlag, uint64(1), nil, nil,
			[]byte{0x0e, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{"Regular8", false, 0, uint64(1), nil, uint64(2),
			[]byte{0x3c, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{"Regular4", false, 0, uint32(1), nil, uint32(2),
			[]byte{0x28, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00}},
		{"Regular2", false, 0, uint32(1), nil, uint16(2),
			[]byte{0x18, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00}},
		{"Regular1", false, 0, uint8(1), nil, uint8(2),
			[]byte{0x04, 0x01, 0x02}},
	}

	t.Run("Write", func(t *testing.T) {
		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				buffer := [packet.MaxHeaderSize]byte{}

				header := packet.Header(buffer[:len(testCase.bytes)])
				header.SetFlags(testCase.flags)
				header.SetConnectionID(testCase.connectionID)
				if testCase.version != nil {
					header.SetVersion(testCase.version.(uint32))
				}
				if testCase.packetNumber != nil {
					header.SetPacketNumber(testCase.packetNumber, testCase.special)
				}
				l := header.Len(testCase.special)

				assert.Equal(t, len(testCase.bytes), l)
				assert.Equal(t, testCase.bytes, buffer[:l])
			})
		}
	})

	t.Run("Read", func(t *testing.T) {
		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				header := packet.Header(testCase.bytes)
				assert.Equal(t, testCase.connectionID, header.ConnectionID())
				if testCase.version != nil {
					assert.Equal(t, testCase.version, header.Version())
				}
				if testCase.packetNumber != nil {
					assert.Equal(t, testCase.packetNumber, header.PacketNumber(testCase.special))
				}
			})
		}
	})
}
