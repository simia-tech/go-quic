package packet

type PublicReset []byte

func (pr PublicReset) SetConnectionID(value uint64) {
	header := Header(pr)
	header.SetFlags(PublicResetFlag)
	header.SetConnectionID(value)
}

func (pr PublicReset) ConnectionID() uint64 {
	return Header(pr).ConnectionID().(uint64)
}

func (pr PublicReset) Len() int {
	return len(pr)
}
