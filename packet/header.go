package packet

import (
	"encoding/binary"
	"fmt"
)

// Header flags
const (
	FlagMask        = 0x03
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

// SetFlags sets public header flags.
func (h Header) SetFlags(flag uint8) {
	h.ensureLen(1)
	h[0] |= flag
}

// Flags returns public header flags.
func (h Header) Flags() uint8 {
	h.ensureLen(1)
	return h[0]
}

// SetConnectionID sets the connection id and the corresponding header flags. The value has to be
// uint8, uint16, uint32 or uint64. Values of other types will cause a panic.
func (h Header) SetConnectionID(value interface{}) {
	switch v := value.(type) {
	case uint64:
		h.ensureLen(9)
		h.SetFlags(ConnectionIDLen8)
		binary.LittleEndian.PutUint64(h[1:], v)
	case uint32:
		h.ensureLen(5)
		h.SetFlags(ConnectionIDLen4)
		binary.LittleEndian.PutUint32(h[1:], v)
	case uint8:
		h.ensureLen(2)
		h.SetFlags(ConnectionIDLen1)
		h[1] = v
	default:
		panic(fmt.Sprintf("cannot set connection id value of type %T", v))
	}
}

// ConnectionID returns the connection id.
func (h Header) ConnectionID() interface{} {
	switch h[0] & ConnectionIDMask {
	case ConnectionIDLen8:
		h.ensureLen(9)
		return binary.LittleEndian.Uint64(h[1:])
	case ConnectionIDLen4:
		h.ensureLen(5)
		return binary.LittleEndian.Uint32(h[1:])
	case ConnectionIDLen1:
		h.ensureLen(2)
		return h[1]
	}
	return nil
}

// SetVersion set the quic versions and the corresponding flag.
func (h Header) SetVersion(version uint32) {
	offset := 1 + h.connectionIDLen()
	h.ensureLen(offset + 4)
	h.SetFlags(VersionFlag)
	binary.LittleEndian.PutUint32(h[offset:], version)
}

// Versions returns the versions.
func (h Header) Version() uint32 {
	offset := 1 + h.connectionIDLen()
	h.ensureLen(offset + 4)
	return binary.LittleEndian.Uint32(h[offset:])
}

// SetPacketNumber sets the packet number and the corresponding header flags. The value has to be
// uint8, uint16, uint32 or uint64. Values of other types will cause a panic.
func (h Header) SetPacketNumber(value interface{}, special bool) {
	offset := 1 + h.connectionIDLen() + h.versionsLen(special)
	switch v := value.(type) {
	case uint64:
		h.ensureLen(offset + 6)
		h.SetFlags(PacketNumberLen6)
		binary.LittleEndian.PutUint32(h[offset:], uint32(v&0xffffffff))
		binary.LittleEndian.PutUint16(h[offset+4:], uint16((v&0xffff00000000)>>32))
	case uint32:
		h.ensureLen(offset + 4)
		h.SetFlags(PacketNumberLen4)
		binary.LittleEndian.PutUint32(h[offset:], v)
	case uint16:
		h.ensureLen(offset + 2)
		h.SetFlags(PacketNumberLen2)
		binary.LittleEndian.PutUint16(h[offset:], v)
	case uint8:
		h.ensureLen(offset + 1)
		h.SetFlags(PacketNumberLen1)
		h[offset] = v
	default:
		panic(fmt.Sprintf("cannot set packet number value of type %T", v))
	}
}

// PacketNumber returns the connection id.
func (h Header) PacketNumber(special bool) interface{} {
	offset := 1 + h.connectionIDLen() + h.versionsLen(special)
	switch h[0] & PacketNumberMask {
	case PacketNumberLen6:
		h.ensureLen(offset + 6)
		return uint64(binary.LittleEndian.Uint32(h[offset:])) |
			(uint64(binary.LittleEndian.Uint16(h[offset+4:])) << 32)
	case PacketNumberLen4:
		h.ensureLen(offset + 4)
		return binary.LittleEndian.Uint32(h[offset:])
	case PacketNumberLen2:
		h.ensureLen(offset + 2)
		return binary.LittleEndian.Uint16(h[offset:])
	case PacketNumberLen1:
		if len(h) < offset+1 {
			return nil
		}
		return h[offset]
	}
	return nil
}

// Len returns the length of the header.
func (h Header) Len(special bool) int {
	return 1 + h.connectionIDLen() + h.versionsLen(special) + h.packetNumberLen(special)
}

func (h Header) ensureLen(l int) {
	if len(h) < l {
		panic(fmt.Sprintf("expected buffer to have at least %d bytes, got %d", l, len(h)))
	}
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

func (h Header) versionsLen(special bool) int {
	if !special && h[0]&VersionFlag != 0 {
		return 4
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
