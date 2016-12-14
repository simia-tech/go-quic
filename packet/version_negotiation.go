package packet

import (
	"encoding/binary"
	"fmt"
)

type VersionNegotiation []byte

func (vn VersionNegotiation) SetConnectionID(value uint64) {
	header := Header(vn)
	header.SetFlags(VersionFlag)
	header.SetConnectionID(value)
}

func (vn VersionNegotiation) ConnectionID() uint64 {
	return Header(vn).ConnectionID().(uint64)
}

func (vn VersionNegotiation) SetVersions(values []uint32) {
	header := Header(vn)
	offset := header.Len(true)
	vn.ensureLen(offset + (len(values) * 4))

	for _, value := range values {
		binary.LittleEndian.PutUint32(vn[offset:], value)
		offset += 4
	}
}

func (vn VersionNegotiation) Versions() []uint32 {
	header := Header(vn)
	offset := header.Len(true)
	count := (len(vn) - offset) / 4
	versions := make([]uint32, count)
	for index := 0; index < count; index++ {
		versions[index] = binary.LittleEndian.Uint32(vn[offset:])
		offset += 4
	}
	return versions
}

func (vn VersionNegotiation) Len() int {
	return len(vn)
}

func (vn VersionNegotiation) ensureLen(l int) {
	if len(vn) < l {
		panic(fmt.Sprintf("expected buffer to have at least %d bytes, got %d", l, len(vn)))
	}
}
