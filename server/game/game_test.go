package game

import (
	"encoding/gob"
	"net"
	"testing"

	"github.com/cbeuw/connutil"
)

const lobbySize = 3

func mockGame(lobbySize int) ([]mockRemote, Game) {
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

	game := MakeGame(local...)

	return remote, game
}

func TestBroadcastSendsInterfaceToAllClients(t *testing.T) {
	remotes, game := mockGame(lobbySize)

	expected := []string{"Hello World", "Hola Mundo", "Olá Mundo"}
	game.broadcast(expected)

	for _, remote := range remotes {
		received, err := remote.receive()
		if err != nil {
			t.Errorf("Could not read from remote: %s", err)
		}

		_, ok := received.([]string)
		if !ok {
			t.Errorf("Did not received a round")
		}
	}
}

func TestBroadcastAllButSendsInterfaceToAllClientsButOne(t *testing.T) {
	remotes, game := mockGame(lobbySize)

	expected := []string{"Hello World", "Hola Mundo", "Olá Mundo"}
	game.broadcastAllBut(expected, 0)

	ch := make(chan struct{})
	go func() {
		remotes[0].receive()
		ch <- struct{}{}
	}()

	for _, remote := range remotes[1:] {
		received, err := remote.receive()
		if err != nil {
			t.Errorf("Could not read from remote: %s", err)
		}

		_, ok := received.([]string)
		if !ok {
			t.Errorf("Did not received a round")
		}
	}

	select {
	case <-ch:
		t.Errorf("Should have not received a message")
	default:
	}
}

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
