package frame

import (
	"encoding/binary"
	"fmt"
)

// Stream defines the stream frame.
type Stream []byte

// SetStreamID sets the stream id.
func (s Stream) SetStreamID(value interface{}) {
	frameType := Type(s)
	frameType.SetType(TypeStream)
	offset := frameType.Len()

	switch v := value.(type) {
	case uint32:
		s.ensureLen(offset + 4)
		frameType.SetFlags(FlagStreamIDLen4)
		binary.LittleEndian.PutUint32(s[offset:], v)
	case uint16:
		s.ensureLen(offset + 2)
		frameType.SetFlags(FlagStreamIDLen2)
		binary.LittleEndian.PutUint16(s[offset:], v)
	case uint8:
		s.ensureLen(offset + 1)
		frameType.SetFlags(FlagStreamIDLen1)
		s[offset] = v
	default:
		panic(fmt.Sprintf("cannot set stream id value of type %T", v))
	}
}

// StreamID returns the stream id.
func (s Stream) StreamID() interface{} {
	frameType := Type(s)
	offset := frameType.Len()

	switch v := frameType.Flags() & MaskStreamIDLen; v {
	case FlagStreamIDLen4:
		s.ensureLen(offset + 4)
		return binary.LittleEndian.Uint32(s[offset:])
	case FlagStreamIDLen2:
		s.ensureLen(offset + 2)
		return binary.LittleEndian.Uint16(s[offset:])
	case FlagStreamIDLen1:
		s.ensureLen(offset + 1)
		return s[offset]
	default:
		panic(fmt.Sprintf("cannot get stream id value of type %x", v))
	}
}

// AddOffset adds the offset.
func (s Stream) AddOffset(value interface{}) {
	frameType := Type(s)
	offset := frameType.Len() + s.streamIDLen()

	switch v := value.(type) {
	case uint64:
		s.ensureLen(offset + 8)
		frameType.SetFlags(FlagOffsetLen8)
		binary.LittleEndian.PutUint64(s[offset:], v)
	case uint32:
		s.ensureLen(offset + 4)
		frameType.SetFlags(FlagOffsetLen4)
		binary.LittleEndian.PutUint32(s[offset:], v)
	case uint16:
		s.ensureLen(offset + 2)
		frameType.SetFlags(FlagOffsetLen2)
		binary.LittleEndian.PutUint16(s[offset:], v)
	}
}

// Offset returns the offset.
func (s Stream) Offset() interface{} {
	frameType := Type(s)
	offset := frameType.Len() + s.streamIDLen()

	switch v := frameType.Flags() & MaskOffsetLen; v {
	case FlagOffsetLen8:
		s.ensureLen(offset + 8)
		return binary.LittleEndian.Uint64(s[offset:])
	case FlagOffsetLen4:
		s.ensureLen(offset + 4)
		return binary.LittleEndian.Uint32(s[offset:])
	case FlagOffsetLen2:
		s.ensureLen(offset + 2)
		return binary.LittleEndian.Uint16(s[offset:])
	case FlagOffsetLen0:
		return nil
	default:
		panic(fmt.Sprintf("cannot get stream id value of type %x", v))
	}
}

// SetData sets the payload data.
func (s Stream) SetData(data []byte) {
	frameType := Type(s)
	if len(data) == 0 {
		return
	}

	offset := frameType.Len() + s.streamIDLen() + s.offsetLen()
	s.ensureLen(offset + 2 + len(data))

	frameType.SetFlags(FlagDataLen)
	binary.LittleEndian.PutUint16(s[offset:], uint16(len(data)))

	copy(s[offset+2:], data)
}

// Data returns the payload data.
func (s Stream) Data() []byte {
	frameType := Type(s)
	if frameType.Flags()&FlagDataLen == 0x00 {
		return []byte{}
	}

	offset := frameType.Len() + s.streamIDLen() + s.offsetLen()
	s.ensureLen(offset + 2)

	l := int(binary.LittleEndian.Uint16(s[offset:]))
	offset += 2
	s.ensureLen(offset + l)

	return s[offset : offset+l]
}

// Len returns the length of the stream frame.
func (s Stream) Len() int {
	return Type(s).Len() + s.streamIDLen() + s.offsetLen() + s.dataLen()
}

func (s Stream) ensureLen(l int) {
	if len(s) < l {
		panic(fmt.Sprintf("expected buffer to have at least %d bytes, got %d", l, len(s)))
	}
}

func (s Stream) streamIDLen() int {
	switch Type(s).Flags() & MaskStreamIDLen {
	case FlagStreamIDLen4:
		return 4
	case FlagStreamIDLen3:
		return 3
	case FlagStreamIDLen2:
		return 2
	case FlagStreamIDLen1:
		return 1
	}
	return 1
}

func (s Stream) offsetLen() int {
	switch Type(s).Flags() & MaskOffsetLen {
	case FlagOffsetLen8:
		return 8
	case FlagOffsetLen7:
		return 7
	case FlagOffsetLen6:
		return 6
	case FlagOffsetLen5:
		return 5
	case FlagOffsetLen4:
		return 4
	case FlagOffsetLen3:
		return 3
	case FlagOffsetLen2:
		return 2
	case FlagOffsetLen0:
		return 0
	}
	return 0
}

func (s Stream) dataLen() int {
	frameType := Type(s)
	if frameType.Flags()&FlagDataLen == 0x00 {
		return 0
	}

	offset := frameType.Len() + s.streamIDLen() + s.offsetLen()
	s.ensureLen(offset + 2)

	return 2 + int(binary.LittleEndian.Uint16(s[offset:]))
}
