package frame_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/simia-tech/go-quic/frame"
)

func TestType(t *testing.T) {
	testCases := []struct {
		name string

		frameType uint8
		flags     uint8

		bytes []byte
	}{
		{"StreamCreate", frame.TypeStream, frame.FlagDataLen | frame.FlagOffsetLen8 | frame.FlagStreamIDLen4,
			[]byte{0xbf}},
		{"StreamDelete", frame.TypeStream, frame.FlagFinish,
			[]byte{0xc0}},
		{"Acknowledge", frame.TypeAcknowledge, frame.FlagMultiple | frame.FlagLargestAckedLen6 | frame.FlagAckBlockLen6,
			[]byte{0x6f}},
		{"Padding", frame.TypePadding, 0,
			[]byte{0x00}},
		{"ResetStream", frame.TypeResetStream, 0,
			[]byte{0x01}},
		{"ConnectionClose", frame.TypeConnectionClose, 0,
			[]byte{0x02}},
		{"GoAway", frame.TypeGoAway, 0,
			[]byte{0x03}},
		{"WindowUpdate", frame.TypeWindowUpdate, 0,
			[]byte{0x04}},
		{"Blocked", frame.TypeBlocked, 0,
			[]byte{0x05}},
		{"StopWaiting", frame.TypeStopWaiting, 0,
			[]byte{0x06}},
		{"Ping", frame.TypePing, 0,
			[]byte{0x07}},
	}

	t.Run("Write", func(t *testing.T) {
		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				buffer := make([]byte, len(testCase.bytes))

				ft := frame.Type(buffer)
				ft.SetType(testCase.frameType)
				ft.SetFlags(testCase.flags)

				assert.Equal(t, len(testCase.bytes), ft.Len())
				assert.Equal(t, testCase.bytes, buffer)
			})
		}
	})

	t.Run("Read", func(t *testing.T) {
		for _, testCase := range testCases {
			t.Run(testCase.name, func(t *testing.T) {
				ft := frame.Type(testCase.bytes)
				assert.Equal(t, testCase.frameType, ft.Type())
				assert.Equal(t, testCase.flags, ft.Flags())
			})
		}
	})
}
