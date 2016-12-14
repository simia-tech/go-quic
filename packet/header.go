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
	if len(h) < 1 {
		panic("buffer too small")
	}
	h[0] |= flag
}

// Flags returns public header flags.
func (h Header) Flags() uint8 {
	if len(h) < 1 {
		panic("buffer too small")
	}
	return h[0]
}

// SetConnectionID sets the connection id and the corresponding header flags. The value has to be
// uint8, uint16, uint32 or uint64. Values of other types will cause a panic.
func (h Header) SetConnectionID(value interface{}) {
	switch v := value.(type) {
	case uint64:
		h.SetFlags(ConnectionIDLen8)
		binary.LittleEndian.PutUint64(h[1:], v)
	case uint32:
		h.SetFlags(ConnectionIDLen4)
		binary.LittleEndian.PutUint32(h[1:], v)
	case uint8:
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
		if len(h) < 9 {
			panic("buffer too small")
		}
		return binary.LittleEndian.Uint64(h[1:])
	case ConnectionIDLen4:
		if len(h) < 5 {
			panic("buffer too small")
		}
		return binary.LittleEndian.Uint32(h[1:])
	case ConnectionIDLen1:
		if len(h) < 2 {
			panic("buffer too small")
		}
		return h[1]
	}
	return nil
}

// SetVersion set the quic versions and the corresponding flag.
func (h Header) SetVersion(version uint32) {
	offset := 1 + h.connectionIDLen()
	if len(h) < offset+4 {
		panic("buffer too small")
	}
	h.SetFlags(VersionFlag)
	binary.LittleEndian.PutUint32(h[offset:], version)
}

// Versions returns the versions.
func (h Header) Version() uint32 {
	offset := 1 + h.connectionIDLen()
	if len(h) < offset+4 {
		panic("buffer too small")
	}
	return binary.LittleEndian.Uint32(h[offset:])
}

// SetPacketNumber sets the packet number and the corresponding header flags. The value has to be
// uint8, uint16, uint32 or uint64. Values of other types will cause a panic.
func (h Header) SetPacketNumber(value interface{}) {
	offset := 1 + h.connectionIDLen() + h.versionsLen()
	switch v := value.(type) {
	case uint64:
		if len(h) < offset+6 {
			panic("buffer too small")
		}
		h.SetFlags(PacketNumberLen6)
		binary.LittleEndian.PutUint32(h[offset:], uint32(v&0xffffffff))
		binary.LittleEndian.PutUint16(h[offset+4:], uint16((v&0xffff00000000)>>32))
	case uint32:
		if len(h) < offset+4 {
			panic("buffer too small")
		}
		h.SetFlags(PacketNumberLen4)
		binary.LittleEndian.PutUint32(h[offset:], v)
	case uint16:
		if len(h) < offset+2 {
			panic("buffer too small")
		}
		h.SetFlags(PacketNumberLen2)
		binary.LittleEndian.PutUint16(h[offset:], v)
	case uint8:
		if len(h) < offset+1 {
			panic("buffer too small")
		}
		h.SetFlags(PacketNumberLen1)
		h[offset] = v
	default:
		panic(fmt.Sprintf("cannot set packet number value of type %T", v))
	}
}

// PacketNumber returns the connection id.
func (h Header) PacketNumber() interface{} {
	offset := 1 + h.connectionIDLen() + h.versionsLen()
	switch h[0] & PacketNumberMask {
	case PacketNumberLen6:
		if len(h) < offset+6 {
			panic("buffer too small")
		}
		return uint64(binary.LittleEndian.Uint32(h[offset:])) |
			(uint64(binary.LittleEndian.Uint16(h[offset+4:])) << 32)
	case PacketNumberLen4:
		if len(h) < offset+4 {
			panic("buffer too small")
		}
		return binary.LittleEndian.Uint32(h[offset:])
	case PacketNumberLen2:
		if len(h) < offset+2 {
			panic("buffer too small")
		}
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
	return 1 + h.connectionIDLen() + h.versionsLen() + h.packetNumberLen(special)
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

func (h Header) versionsLen() int {
	if h[0]&VersionFlag != 0 {
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
