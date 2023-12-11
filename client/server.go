package main

import (
	"encoding/gob"
	"net"
)

// Encapsulates a server connection in a game.
type server struct {
	encoder *gob.Encoder
	decoder *gob.Decoder
}

// Creates a new server
func makeServer(connection net.Conn) server {
	encoder := gob.NewEncoder(connection)
	decoder := gob.NewDecoder(connection)

	return server{
		encoder: encoder,
		decoder: decoder,
	}
}

// Sends a gob-enconded structure to the server, as an interface.
func (c *server) send(structure any) error {
	return c.encoder.Encode(&structure)
}

// Receives a gob-enconded interface from the server.
func (c *server) receive() (message any, err error) {
	err = c.decoder.Decode(&message)
	return
}
