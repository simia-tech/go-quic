package packet

import (
	"encoding/binary"
	"fmt"
)

// Regular defines the regular packet type.
type Regular []byte

// AddConnectionID adds the connection id.
func (r Regular) AddConnectionID(value uint64) {
	Header(r).AddConnectionID(value)
}

// ConnectionID returns the connection id.
func (r Regular) ConnectionID() uint64 {
	return Header(r).ConnectionID()
}

// AddVersion set the quic versions and the corresponding flag.
func (r Regular) AddVersion(version uint32) {
	header := Header(r)
	offset := header.Len()
	r.ensureLen(offset + 4)
	header.SetFlags(FlagVersion)
	binary.LittleEndian.PutUint32(r[offset:], version)
}

// Version returns the version.
func (r Regular) Version() uint32 {
	header := Header(r)
	offset := header.Len()
	r.ensureLen(offset + 4)
	return binary.LittleEndian.Uint32(r[offset:])
}

// AddPacketNumber sets the packet number and the corresponding header flags. The value has to be
// uint8, uint16, uint32 or uint64. Values of other types will cause a panic.
func (r Regular) AddPacketNumber(value interface{}) {
	header := Header(r)
	offset := header.Len() + r.versionLen()
	switch v := value.(type) {
	case uint64:
		r.ensureLen(offset + 6)
		header.SetFlags(PacketNumberLen6)
		binary.LittleEndian.PutUint32(r[offset:], uint32(v&0xffffffff))
		binary.LittleEndian.PutUint16(r[offset+4:], uint16((v&0xffff00000000)>>32))
	case uint32:
		r.ensureLen(offset + 4)
		header.SetFlags(PacketNumberLen4)
		binary.LittleEndian.PutUint32(r[offset:], v)
	case uint16:
		r.ensureLen(offset + 2)
		header.SetFlags(PacketNumberLen2)
		binary.LittleEndian.PutUint16(r[offset:], v)
	case uint8:
		r.ensureLen(offset + 1)
		header.SetFlags(PacketNumberLen1)
		r[offset] = v
	default:
		panic(fmt.Sprintf("cannot set packet number value of type %T", v))
	}
}

// PacketNumber returns the connection id.
func (r Regular) PacketNumber() interface{} {
	header := Header(r)
	offset := header.Len() + r.versionLen()
	switch r[0] & PacketNumberMask {
	case PacketNumberLen6:
		r.ensureLen(offset + 6)
		return uint64(binary.LittleEndian.Uint32(r[offset:])) |
			(uint64(binary.LittleEndian.Uint16(r[offset+4:])) << 32)
	case PacketNumberLen4:
		r.ensureLen(offset + 4)
		return binary.LittleEndian.Uint32(r[offset:])
	case PacketNumberLen2:
		r.ensureLen(offset + 2)
		return binary.LittleEndian.Uint16(r[offset:])
	case PacketNumberLen1:
		r.ensureLen(offset + 1)
		return r[offset]
	}
	return nil
}

// SetData sets the packet's payload data.
func (r Regular) SetData(data []byte) {
	header := Header(r)
	offset := header.Len() + r.versionLen() + r.packetNumberLen()
	r.ensureLen(offset + len(data))
	copy(r[offset:], data)
}

// Data returns the packet's payload data.
func (r Regular) Data() []byte {
	header := Header(r)
	offset := header.Len() + r.versionLen() + r.packetNumberLen()
	return r[offset:]
}

// Len returns the length of the packet including the header.
func (r Regular) Len() int {
	return len(r)
}

func (r Regular) ensureLen(l int) {
	if len(r) < l {
		panic(fmt.Sprintf("expected buffer to have at least %d bytes, got %d", l, len(r)))
	}
}

func (r Regular) versionLen() int {
	if Header(r).Flags()&FlagVersion == 0x00 {
		return 0
	}
	return 4
}

func (r Regular) packetNumberLen() int {
	switch Header(r).Flags() & PacketNumberMask {
	case PacketNumberLen6:
		return 6
	case PacketNumberLen4:
		return 4
	case PacketNumberLen2:
		return 2
	case PacketNumberLen1:
		return 1
	}
	return 0
}
