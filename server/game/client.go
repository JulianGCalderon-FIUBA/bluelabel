package game

import (
	"encoding/gob"
	"net"
)

type client struct {
	encoder *gob.Encoder
	decoder *gob.Decoder
}

func newClient(connection net.Conn) *client {
	encoder := gob.NewEncoder(connection)
	decoder := gob.NewDecoder(connection)

	return &client{
		encoder,
		decoder,
	}
}
