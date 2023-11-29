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

	var received any
	decoder.Decode(&received)

	if !slices.Equal(expected, received.([]string)) {
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
	received, err := client.receive()
	if err != nil {
		t.Errorf("Could not receive from connection: %s", err)
	}

	if !slices.Equal(expected.([]string), received.([]string)) {
		t.Errorf("Expected %v, but received %v", expected, received)
	}
}

func TestCanReceiveAGobEncodedConcreteStructure(t *testing.T) {
	remote, local := net.Pipe()

	remoteClient := makeClient(remote)

	expected := []string{"Hello World", "Hola Mundo", "Olá Mundo"}
	go remoteClient.send(expected)

	localClient := makeClient(local)
	received, err := receiveConcrete[[]string](localClient)
	if err != nil {
		t.Errorf("Could not receive from connection: %s", err)
	}

	if !slices.Equal(expected, received) {
		t.Errorf("Expected %v, but received %v", expected, received)
	}
}

func TestReceivingConcreteIgnoresIncorrectTypes(t *testing.T) {
	remote, local := net.Pipe()

	remoteClient := makeClient(remote)

	incorrect := []int{1, 2, 3, 4, 5}
	go remoteClient.send(incorrect)
	expected := []string{"Hello World", "Hola Mundo", "Olá Mundo"}
	go remoteClient.send(expected)

	localClient := makeClient(local)
	received, err := receiveConcrete[[]string](localClient)
	if err != nil {
		t.Errorf("Could not receive from connection: %s", err)
	}

	if !slices.Equal(expected, received) {
		t.Errorf("Expected %v, but received %v", expected, received)
	}
}
