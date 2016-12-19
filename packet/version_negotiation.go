package packet

import (
	"encoding/binary"
	"fmt"
)

// VersionNegotiation defines the version negotiation packet type.
type VersionNegotiation []byte

// SetConnectionID sets the connection id.
func (vn VersionNegotiation) SetConnectionID(value uint64) {
	header := Header(vn)
	header.SetFlags(FlagVersion)
	header.AddConnectionID(value)
}

// ConnectionID returns the connection id.
func (vn VersionNegotiation) ConnectionID() uint64 {
	return Header(vn).ConnectionID()
}

// SetVersions sets the versions.
func (vn VersionNegotiation) SetVersions(values []uint32) {
	header := Header(vn)
	offset := header.Len()
	vn.ensureLen(offset + (len(values) * 4))

	for _, value := range values {
		binary.LittleEndian.PutUint32(vn[offset:], value)
		offset += 4
	}
}

// Versions returns the versions.
func (vn VersionNegotiation) Versions() []uint32 {
	header := Header(vn)
	offset := header.Len()
	count := (len(vn) - offset) / 4
	versions := make([]uint32, count)
	for index := 0; index < count; index++ {
		versions[index] = binary.LittleEndian.Uint32(vn[offset:])
		offset += 4
	}
	return versions
}

// Len returns the length of the packet including the header.
func (vn VersionNegotiation) Len() int {
	return len(vn)
}

func (vn VersionNegotiation) ensureLen(l int) {
	if len(vn) < l {
		panic(fmt.Sprintf("expected buffer to have at least %d bytes, got %d", l, len(vn)))
	}
}
