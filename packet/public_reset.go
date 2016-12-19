package packet

// PublicReset defines the public reset packet type.
type PublicReset []byte

// SetConnectionID sets the connection id.
func (pr PublicReset) SetConnectionID(value uint64) {
	header := Header(pr)
	header.SetFlags(PublicResetFlag)
	header.AddConnectionID(value)
}

// ConnectionID returns the connection id.
func (pr PublicReset) ConnectionID() uint64 {
	return Header(pr).ConnectionID()
}

// Len returns the length of the packet including the header.
func (pr PublicReset) Len() int {
	return len(pr)
}
