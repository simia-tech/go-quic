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
	Nonce           = 0x04
	ConnectionID    = 0x08

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

// AddConnectionID adds the connection id and sets the corresponding header flags. The value has to be
// uint8, uint16, uint32 or uint64. Values of other types will cause a panic.
func (h Header) AddConnectionID(value uint64) {
	h.ensureLen(9)
	h.SetFlags(ConnectionID)
	binary.LittleEndian.PutUint64(h[1:], value)
}

// ConnectionID returns the connection id.
func (h Header) ConnectionID() uint64 {
	if h[0]&ConnectionID == 0x00 {
		return 0
	}
	h.ensureLen(9)
	return binary.LittleEndian.Uint64(h[1:])
}

// Len returns the length of the header.
func (h Header) Len() int {
	return 1 + h.connectionIDLen()
}

func (h Header) ensureLen(l int) {
	if len(h) < l {
		panic(fmt.Sprintf("expected buffer to have at least %d bytes, got %d", l, len(h)))
	}
}

func (h Header) connectionIDLen() int {
	if h[0]&ConnectionID == 0x00 {
		return 0
	}
	return 8
}
