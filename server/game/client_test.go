package game

import (
	"encoding/gob"
	"net"
	"slices"
	"testing"
)

func TestCanSendAGobEncodedStructure(t *testing.T) {
	remote, local := net.Pipe()

	client := makeClient(local)

	expected := []string{"Hello World", "Hola Mundo", "Olá Mundo"}
	go client.send(expected)

	decoder := gob.NewDecoder(remote)

	var received []string
	decoder.Decode(&received)

	if !slices.Equal(expected, received) {
		t.Errorf("Expected %v, but received %v", expected, received)
	}
}

func TestCanReceiveAGobEncodedStructure(t *testing.T) {
	remote, local := net.Pipe()

	var expected any
	expected = []string{"Hello World", "Hola Mundo", "Olá Mundo"}
	encoder := gob.NewEncoder(remote)
	go encoder.Encode(&expected)

	client := makeClient(local)
	received := client.receive()

	if !slices.Equal(expected.([]string), received.([]string)) {
		t.Errorf("Expected %v, but received %v", expected, received)
	}
}
