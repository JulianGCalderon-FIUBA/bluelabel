package game

import (
	"bluelabel/shared"
	"encoding/gob"
	"net"
)

// Encapsulates a client connection in a game
type client struct {
	encoder *gob.Encoder
	decoder *gob.Decoder
	stops   chan shared.StopRequest
}

// Creates a new client
func makeClient(connection net.Conn) client {
	encoder := gob.NewEncoder(connection)
	decoder := gob.NewDecoder(connection)

	return client{
		encoder: encoder,
		decoder: decoder,
		stops:   make(chan shared.StopRequest),
	}
}

// Sends a gob-enconded structure to the client
func (c *client) send(structure any) error {
	return c.encoder.Encode(structure)
}

// Sends a gob-enconded interface from the client
func (c *client) receive() (message any) {
	c.decoder.Decode(&message)
	return
}

func (c *client) loop() {
	for {
		msg := c.receive()
		switch received := msg.(type) {
		case shared.StopRequest:
			c.stops <- received
		default:
		}
	}
}
