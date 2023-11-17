package game_test

import (
	"bluelabel/server/game"
	"bluelabel/shared"
	"encoding/gob"
	"net"
	"testing"

	"github.com/cbeuw/connutil"
)

const lobbySize = 3

type mockRemote struct {
	*gob.Decoder
	*gob.Encoder
}

func (c *mockRemote) send(structure any) error {
	return c.Encoder.Encode(structure)
}

func (c *mockRemote) receive() (message any, err error) {
	err = c.Decoder.Decode(&message)
	return
}

func mockGame(lobbySize int) []mockRemote {
	local := make([]net.Conn, lobbySize)
	remote := make([]mockRemote, lobbySize)
	for i := range local {
		localConn, remoteConn := connutil.AsyncPipe()

		local[i] = localConn
		remote[i] = mockRemote{
			gob.NewDecoder(remoteConn),
			gob.NewEncoder(remoteConn),
		}
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
