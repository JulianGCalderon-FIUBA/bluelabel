package game

import (
	"bluelabel/shared"
	"encoding/gob"
	"maps"
	"net"
	"slices"
	"testing"

	"github.com/cbeuw/connutil"
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

func TestSkipsMessageIfNoWaiter(t *testing.T) {
	_, local := connutil.AsyncPipe()

	client := makeClient(local)

	expected := shared.StopRequest{}
	var message any = expected
	client.handleMessage(&message)

	stop := make(chan shared.StopRequest)

	go func() {
		stop <- client.receiveStop()
	}()

	go client.handleMessage(&message)

	received := <-stop

	if !maps.Equal(expected.Words, received.Words) {
		t.Fatalf("Expected %v, but received %v", message, received)
	}
}
