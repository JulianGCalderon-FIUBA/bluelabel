package conntest

import (
	"encoding/gob"
	"net"
)

type MockRemote struct {
	gob.Decoder
	gob.Encoder
}

func NewMockRemote(c net.Conn) *MockRemote {
	return &MockRemote{
		*gob.NewDecoder(c),
		*gob.NewEncoder(c),
	}
}

func (c *MockRemote) Send(structure any) error {
	return c.Encoder.Encode(structure)
}

func (c *MockRemote) Receive() (message any, err error) {
	err = c.Decoder.Decode(&message)
	return
}
