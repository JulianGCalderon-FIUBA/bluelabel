package game

import (
	"bluelabel/shared"
	"encoding/gob"
	"net"
	"testing"
	"time"

	"github.com/cbeuw/connutil"
)

const lobbySize = 3

func mockGame(lobbySize int) ([]mockRemote, Game) {
	local := make([]net.Conn, lobbySize)
	remote := make([]mockRemote, lobbySize)
	for i := range local {
		localConn, remoteConn := connutil.AsyncPipe()

		mockRemote := newMockRemote(remoteConn)

		local[i] = localConn
		remote[i] = *mockRemote
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

func TestCanWaitForFirstStopRequest(t *testing.T) {
	remotes, game := mockGame(lobbySize)

	stops := buildStopArray(game.clients...)
	indexedStops := mergeChannels(stops...)

	remotes[0].sendInterface(shared.StopRequest{})

	senderId := game.waitOneStop(indexedStops)

	if senderId != 0 {
		t.Errorf("Expected sender id %v, but got %v", 0, senderId)
	}
}

func TestCanWaitForAllStopRequest(t *testing.T) {
	remotes, game := mockGame(lobbySize)

	stops := buildStopArray(game.clients...)
	indexedStops := mergeChannels(stops...)

	for i := range remotes {
		remotes[i].sendInterface(shared.StopRequest{})
	}

	game.waitAllStop(indexedStops, time.Hour)
}

func TestCanWaitForAllStopRequestWithTimeout(t *testing.T) {
	remotes, game := mockGame(lobbySize)

	stops := buildStopArray(game.clients...)
	indexedStops := mergeChannels(stops...)

	for i := 0; i < len(remotes)/2; i++ {
		remotes[i].sendInterface(shared.StopRequest{})
	}

	game.waitAllStop(indexedStops, 100*time.Millisecond)

	if len(game.words) != len(remotes)/2 {
		t.Errorf("Expected to receive %v stop requests, but received %v", len(remotes)/2, len(game.words))
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
	return c.Encoder.Encode(structure)
}

func (c *mockRemote) sendInterface(structure any) error {
	return c.Encoder.Encode(&structure)
}

func (c *mockRemote) receive() (message any, err error) {
	err = c.Decoder.Decode(&message)
	return
}
