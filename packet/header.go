package packet

import (
	"encoding/binary"
	"fmt"
)

// Header flags
const (
	PublicVersionFlag = 0x01

	PublicConnectionIDSize64 = 0x0c
)

// Header defines the packet header type.
type Header []byte

// SetPublic sets public header flags.
func (h Header) SetPublic(flag uint8) {
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
		h.SetPublic(PublicConnectionIDSize64)
		binary.LittleEndian.PutUint64(h[1:], v)
	default:
		panic(fmt.Sprintf("cannot set connection id value of type %T", v))
	}
}

// Len returns the length of the header.
func (h Header) Len() int {
	return 0
}
