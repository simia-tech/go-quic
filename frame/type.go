package frame

import "fmt"

// Definition of frame types, flags and masks.
const (
	TypeStream          = 0x80
	TypeAcknowledge     = 0x40
	TypePadding         = 0x00
	TypeResetStream     = 0x01
	TypeConnectionClose = 0x02
	TypeGoAway          = 0x03
	TypeWindowUpdate    = 0x04
	TypeBlocked         = 0x05
	TypeStopWaiting     = 0x06
	TypePing            = 0x07

	// Flags for the stream type.
	FlagFinish       = 0x40
	FlagDataLen      = 0x20
	FlagOffsetLen8   = 0x07 << 2
	FlagOffsetLen7   = 0x06 << 2
	FlagOffsetLen6   = 0x05 << 2
	FlagOffsetLen5   = 0x04 << 2
	FlagOffsetLen4   = 0x03 << 2
	FlagOffsetLen3   = 0x02 << 2
	FlagOffsetLen2   = 0x01 << 2
	FlagOffsetLen0   = 0x00 << 2
	FlagStreamIDLen4 = 0x03
	FlagStreamIDLen3 = 0x02
	FlagStreamIDLen2 = 0x01
	FlagStreamIDLen1 = 0x00

	// Masks for the stream type.
	MaskStream      = 0x7f
	MaskOffsetLen   = 0x07 << 2
	MaskStreamIDLen = 0x03

	// Flags for the acknowledge type.
	FlagMultiple         = 0x20
	FlagLargestAckedLen6 = 0x03 << 2
	FlagLargestAckedLen4 = 0x02 << 2
	FlagLargestAckedLen2 = 0x01 << 2
	FlagLargestAckedLen1 = 0x00 << 2
	FlagAckBlockLen6     = 0x03
	FlagAckBlockLen4     = 0x02
	FlagAckBlockLen2     = 0x01
	FlagAckBlockLen1     = 0x00

	// Masks for the acknowledge type.
	MaskAcknowledge     = 0x3f
	MaskLargestAckedLen = 0x03 << 2
	MaskAckBlockLen     = 0x03
)

// Type defines the frame type.
type Type []byte

// SetType sets the frame type.
func (t Type) SetType(value uint8) {
	t.ensureLen(1)
	t[0] = value
}

// Type returns the frame type.
func (t Type) Type() uint8 {
	t.ensureLen(1)
	if t[0]&TypeStream != 0 {
		return TypeStream
	}
	if t[0]&TypeAcknowledge != 0 {
		return TypeAcknowledge
	}
	return t[0]
}

// SetFlags sets the flags for the stream and acknowledge type.
func (t Type) SetFlags(value uint8) {
	ft := t.Type()
	if ft == TypeStream {
		t[0] |= value & MaskStream
	}
	if ft == TypeAcknowledge {
		t[0] |= value & MaskAcknowledge
	}
}

// Flags returns the flags for the stream and acknowledge type.
func (t Type) Flags() uint8 {
	ft := t.Type()
	if ft == TypeStream {
		return t[0] & MaskStream
	}
	if ft == TypeAcknowledge {
		return t[0] & MaskAcknowledge
	}
	return 0
}

// Len returns the length of the frame type.
func (t Type) Len() int {
	return 1
}

func (t Type) ensureLen(l int) {
	if len(t) < l {
		panic(fmt.Sprintf("expected buffer to have at least %d bytes, got %d", l, len(t)))
	}
}
