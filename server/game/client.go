package game

import (
	"encoding/gob"
	"net"
)

// Encapsulates a client connection in a game
type client struct {
	encoder *gob.Encoder
	decoder *gob.Decoder
}

// Creates a new client
func makeClient(connection net.Conn) client {
	encoder := gob.NewEncoder(connection)
	decoder := gob.NewDecoder(connection)

	return client{
		encoder,
		decoder,
	}
}

// Sends a gob-enconded structure to the client, as an interface.
func (c *client) send(structure any) error {
	return c.encoder.Encode(structure)
}
