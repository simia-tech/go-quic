package packet_test

import (
	"testing"

	"github.com/simia-tech/go-quic/packet"
	"github.com/stretchr/testify/assert"
)

func TestRegular(t *testing.T) {
	testCases := []struct {
		name string

		connectionID interface{}
		version      interface{}
		nonce        []byte
		packetNumber interface{}
		data         []byte

		bytes []byte
	}{
		{"Version", uint64(1), uint32(2), nil, uint64(3), []byte{0x04, 0x05, 0x06},
			[]byte{0x39, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x05, 0x06}},
		{"PacketNumber6", uint64(1), nil, nil, uint64(2), []byte{0x04, 0x05, 0x06},
			[]byte{0x38, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x05, 0x06}},
		{"PacketNumber4", uint64(1), nil, nil, uint32(2), []byte{0x04, 0x05, 0x06},
			[]byte{0x28, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x04, 0x05, 0x06}},
		{"PacketNumber2", uint64(1), nil, nil, uint16(2), []byte{0x04, 0x05, 0x06},
			[]byte{0x18, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x04, 0x05, 0x06}},
		{"PacketNumber1", uint64(1), nil, nil, uint8(2), []byte{0x04, 0x05, 0x06},
			[]byte{0x08, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x04, 0x05, 0x06}},
	}

	t.Run("Write", func(t *testing.T) {
		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				buffer := make([]byte, len(testCase.bytes))

				regular := packet.Regular(buffer)
				if testCase.connectionID != nil {
					regular.AddConnectionID(testCase.connectionID.(uint64))
				}
				if testCase.version != nil {
					regular.AddVersion(testCase.version.(uint32))
				}
				if testCase.nonce != nil {
					// regular.AddNonce(testCase.nonce)
				}
				if testCase.packetNumber != nil {
					regular.AddPacketNumber(testCase.packetNumber)
				}
				regular.SetData(testCase.data)

				assert.Equal(t, len(testCase.bytes), regular.Len())
				assert.Equal(t, testCase.bytes, buffer)
			})
		}
	})

	t.Run("Read", func(t *testing.T) {
		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				regular := packet.Regular(testCase.bytes)
				assert.Equal(t, testCase.connectionID, regular.ConnectionID())
				if testCase.version != nil {
					assert.Equal(t, testCase.version, regular.Version())
				}
				if testCase.nonce != nil {
					// assert.Equal(t, testCase.nonce, regular.Nonce())
				}
				if testCase.packetNumber != nil {
					assert.Equal(t, testCase.packetNumber, regular.PacketNumber())
				}
				assert.Equal(t, testCase.data, regular.Data())
			})
		}
	})
}
