package game

import (
	"encoding/gob"
	"net"
	"slices"
	"sync"
	"testing"
	"time"
)

const lobbySize = 3

func mockGame(lobbySize int) ([]net.Conn, Game) {
	remote := make([]net.Conn, lobbySize)
	local := make([]net.Conn, lobbySize)

	for i := 0; i < lobbySize; i++ {
		remote[i], local[i] = net.Pipe()
	}

	game := MakeGame(local...)

	return remote, game
}

func TestBroadcastSendsMessageToAllConnections(t *testing.T) {
	remoteConnections, game := mockGame(lobbySize)

	mockMessage := []string{"Hello World", "Hola Mundo", "Olá Mundo"}

	go func() {
		err := game.broadcast(mockMessage)
		if err != nil {
			t.Errorf("Could not broadcast message: %v", err)
		}
	}()

	connectionWaitGroup := sync.WaitGroup{}
	connectionWaitGroup.Add(lobbySize)
	for _, c := range remoteConnections {
		go func(c net.Conn) {
			c.SetDeadline(time.Now().Add(1 * time.Second))
			decoder := gob.NewDecoder(c)

			var receivedInterface any
			err := decoder.Decode(&receivedInterface)
			if err != nil {
				t.Errorf("Could not decode message from server: %v", err)
			}

			receivedMessage, ok := receivedInterface.([]string)
			if !ok {
				t.Errorf("Message received was not the correct type: %v", err)
			}

			if !slices.Equal(mockMessage, receivedMessage) {
				t.Errorf("Expected %v, but got %v", mockMessage, receivedMessage)
			}

			connectionWaitGroup.Done()
		}(c)
	}

	connectionWaitGroup.Wait()
}
