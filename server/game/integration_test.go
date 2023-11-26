package game_test

import (
	"bluelabel/server/game"
	"bluelabel/shared"
	"encoding/gob"
	"net"
	"testing"

	"github.com/cbeuw/connutil"
)

// Este archivo tiene codigo repetido con los `game_tests.go`, pero creo que esta repeticion es accidental, y mas adelante puede variar.
// Debido a esto, no me apuro a corregirlo. La escencia de los tests es distinta.

const lobbySize = 3

func mockGame(lobbySize int) []mockRemote {
	local := make([]net.Conn, lobbySize)
	remote := make([]mockRemote, lobbySize)
	for i := range local {
		localConn, remoteConn := connutil.AsyncPipe()

		mockRemote := newMockRemote(remoteConn)

		local[i] = localConn
		remote[i] = *mockRemote

	}

	game := game.MakeGame(local...)
	go game.PlayGame()

	return remote
}

func TestClientsReceiveInitialRound(t *testing.T) {
	remotes := mockGame(lobbySize)

	for _, remote := range remotes {
		received, err := remote.receive()
		if err != nil {
			t.Errorf("Could not read from remote: %s", err)
		}

		_, ok := received.(shared.Round)
		if !ok {
			t.Errorf("Did not received a round")
		}
	}
}

func TestClientsReceiveStopNotify(t *testing.T) {
	remotes := mockGame(lobbySize)

	for _, remote := range remotes {
		remote.receive()
	}

	err := remotes[0].send(shared.StopRequest{})
	if err != nil {
		t.Errorf("Could not send stop request to server: %s", err)
	}

	for _, remote := range remotes[1:] {
		received, err := remote.receive()
		if err != nil {
			t.Errorf("Could not read from remote: %s", err)
		}

		_, ok := received.(shared.StopNotify)
		if !ok {
			t.Errorf("Did not received a stop notify")
		}
	}
}

type mockRemote struct {
	*gob.Decoder
	*gob.Encoder
}

func newMockRemote(c net.Conn) *mockRemote {
	return &mockRemote{
		gob.NewDecoder(c),
		gob.NewEncoder(c),
	}
}

func (c *mockRemote) send(structure any) error {
	return c.Encoder.Encode(&structure)
}

func (c *mockRemote) receive() (message any, err error) {
	err = c.Decoder.Decode(&message)
	return
}
