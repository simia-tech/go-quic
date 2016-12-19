package frame_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/simia-tech/go-quic/frame"
)

func TestStream(t *testing.T) {
	testCases := []struct {
		name string

		streamID interface{}
		offset   interface{}
		data     []byte

		bytes []byte
	}{
		{"RegularStreamID4Offset8", uint32(1), uint64(2), []byte{0x03, 0x04, 0x05},
			[]byte{0xbf, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x00, 0x03, 0x04, 0x05}},
		{"RegularStreamID2Offset4", uint16(1), uint32(2), []byte{0x03, 0x04, 0x05},
			[]byte{0xad, 0x01, 0x00, 0x02, 0x00, 0x00, 0x00, 0x03, 0x00, 0x03, 0x04, 0x05}},
		{"RegularStreamID1Offset2", uint8(1), uint16(2), []byte{0x03, 0x04, 0x05},
			[]byte{0xa4, 0x01, 0x02, 0x00, 0x03, 0x00, 0x03, 0x04, 0x05}},
		{"RegularStreamID1Offset0", uint8(1), nil, []byte{0x03, 0x04, 0x05},
			[]byte{0xa0, 0x01, 0x03, 0x00, 0x03, 0x04, 0x05}},
		{"RegularNoData", uint32(1), uint64(2), []byte{},
			[]byte{0x9f, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	}

	t.Run("Write", func(t *testing.T) {
		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				buffer := make([]byte, len(testCase.bytes))

				stream := frame.Stream(buffer)
				stream.SetStreamID(testCase.streamID)
				if testCase.offset != nil {
					stream.AddOffset(testCase.offset)
				}
				stream.SetData(testCase.data)

				assert.Equal(t, len(testCase.bytes), stream.Len())
				assert.Equal(t, testCase.bytes, buffer)
			})
		}
	})

	t.Run("Read", func(t *testing.T) {
		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				stream := frame.Stream(testCase.bytes)
				assert.Equal(t, testCase.streamID, stream.StreamID())
				assert.Equal(t, testCase.offset, stream.Offset())
				assert.Equal(t, testCase.data, stream.Data())
			})
		}
	})
}
