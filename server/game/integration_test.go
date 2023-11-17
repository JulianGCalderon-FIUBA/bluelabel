package game_test

import (
	"bluelabel/server/conntest"
	"bluelabel/server/game"
	"bluelabel/shared"
	"net"
	"testing"

	"github.com/cbeuw/connutil"
)

const lobbySize = 3

func mockGame(lobbySize int) []*conntest.MockRemote {
	local := make([]net.Conn, lobbySize)
	remote := make([]*conntest.MockRemote, lobbySize)
	for i := range local {
		localConn, remoteConn := connutil.AsyncPipe()

		local[i] = localConn
		remote[i] = conntest.NewMockRemote(remoteConn)
	}

	game := game.MakeGame(local...)
	go game.PlayGame()

	return remote
}

func TestClientsReceiveInitialRound(t *testing.T) {
	remotes := mockGame(lobbySize)

	for _, remote := range remotes {
		received, err := remote.Receive()
		if err != nil {
			t.Errorf("Could not read from remote: %s", err)
		}

		_, ok := received.(shared.Round)
		if !ok {
			t.Errorf("Did not received a round")
		}
	}
}
