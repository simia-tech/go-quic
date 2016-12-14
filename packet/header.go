package packet

import (
	"encoding/binary"
	"fmt"
)

// Header flags
const (
	VersionFlag     = 0x01
	PublicResetFlag = 0x02

	ConnectionIDMask = 0x0c
	ConnectionIDLen8 = 0x0c
	ConnectionIDLen4 = 0x08
	ConnectionIDLen1 = 0x04
	ConnectionIDLen0 = 0x00

	PacketNumberMask = 0x30
	PacketNumberLen6 = 0x30
	PacketNumberLen4 = 0x20
	PacketNumberLen2 = 0x10
	PacketNumberLen1 = 0x00

	MaxHeaderSize = 19
)

// Header defines the packet header type.
type Header []byte

// SetFlag sets public header flags.
func (h Header) SetFlag(flag uint8) {
	if len(h) < 1 {
		panic("buffer too small")
	}
	h[0] |= flag
}

// SetConnectionID sets the connection id and the corresponding header flags. The value has to be
// uint8, uint16, uint32 or uint64. Values of other types will cause a panic.
func (h Header) SetConnectionID(value interface{}) {
	switch v := value.(type) {
	case uint64:
		h.SetFlag(ConnectionIDLen8)
		binary.LittleEndian.PutUint64(h[1:], v)
	default:
		panic(fmt.Sprintf("cannot set connection id value of type %T", v))
	}
}

// SetPacketNumber sets the packet number and the corresponding header flags. The value has to be
// uint8, uint16, uint32 or uint64. Values of other types will cause a panic.
func (h Header) SetPacketNumber(value interface{}) {
	offset := 1 + h.connectionIDLen()
	switch v := value.(type) {
	case uint64:
		h.SetFlag(PacketNumberLen6)
		binary.LittleEndian.PutUint32(h[offset:], uint32(v&0xffffffff))
		binary.LittleEndian.PutUint16(h[offset+4:], uint16((v&0xffff00000000)>>32))
	default:
		panic(fmt.Sprintf("cannot set connection id value of type %T", v))
	}
}

// Len returns the length of the header.
func (h Header) Len(special bool) int {
	return 1 + h.connectionIDLen() + h.packetNumberLen(special)
}

func (h Header) connectionIDLen() int {
	switch h[0] & ConnectionIDMask {
	case ConnectionIDLen8:
		return 8
	case ConnectionIDLen4:
		return 4
	case ConnectionIDLen1:
		return 1
	}
	return 0
}

func (h Header) packetNumberLen(special bool) int {
	switch h[0] & PacketNumberMask {
	case PacketNumberLen6:
		return 6
	case PacketNumberLen4:
		return 4
	case PacketNumberLen2:
		return 2
	}
	if special {
		return 0
	}
	return 1
}
