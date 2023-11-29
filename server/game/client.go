package game

import (
	"encoding/gob"
	"net"
)

// Encapsulates a client connection in a game. The client can be passed by copy
type client struct {
	encoder *gob.Encoder
	decoder *gob.Decoder
}

// Creates a new client
func makeClient(connection net.Conn) client {
	encoder := gob.NewEncoder(connection)
	decoder := gob.NewDecoder(connection)

	return client{
		encoder: encoder,
		decoder: decoder,
	}
}

// Sends a gob-enconded structure to the client, as an interface
// This function is blocking and thread-safe.
func (c *client) send(structure any) error {
	return c.encoder.Encode(&structure)
}

// Receives a gob-enconded interface from the client
// This function is blocking and thread-safe.
func (c *client) receive() (message any, err error) {
	err = c.decoder.Decode(&message)
	return
}

// Receives a gob-encoded concrete message the client. If the interface received
// is not from the correct type, it is ignored. If it does not receive an
// interface, it fails.
func receiveConcrete[T any](c client) (concrete T, err error) {
	for {
		var msg any
		msg, err = c.receive()
		if err != nil {
			return
		}

		var ok bool
		concrete, ok = msg.(T)
		if ok {
			return
		}
	}
}
