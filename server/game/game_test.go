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
	err := game.broadcast(expected)
	if err != nil {
		t.Errorf("Could not broadcast to clients: %s", err)
	}

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

func TestCanGetWordListByCategory(t *testing.T) {
	wordsByCategoryForClient := map[int]map[shared.Category]string{
		0: {
			shared.Color:   "Rojo",
			shared.Animal:  "Rana",
			shared.Pais:    "Rumania",
			shared.Deporte: "Rugby",
		},
		1: {
			shared.Color:  "Rosa",
			shared.Animal: "Renacuajo",
			shared.Pais:   "Rusia",
		},
		2: {
			shared.Color:  "Rufo",
			shared.Animal: "Rana",
			shared.Pais:   "Rusia",
		},
	}

	expected := map[shared.Category][]string{
		shared.Color:   {"Rufo", "Rojo", "Rosa"},
		shared.Animal:  {"Rana", "Renacuajo"},
		shared.Pais:    {"Rusia", "Rumania"},
		shared.Deporte: {"Rugby"},
	}

	result := buildWordListByCategory(wordsByCategoryForClient)

	if !maps.EqualFunc(expected, result, func(m1, m2 []string) bool {
		slices.Sort(m1)
		slices.Sort(m2)

		return slices.Equal(m1, m2)
	}) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

}

func TestCanGetInvalidWordsFromScrutiny(t *testing.T) {
	scrutiny := make(map[shared.Category]map[string]int)

	cat1_voting := make(map[string]int)
	cat1_voting["Perro"] = 3
	cat1_voting["Gato"] = 1
	cat1_voting["Vaca"] = 2
	cat1_voting["Oveja"] = 4

	cat2_voting := make(map[string]int)
	cat2_voting["Perro"] = 3
	cat2_voting["Gato"] = 1
	cat2_voting["Vaca"] = 2

	scrutiny[shared.Pais] = cat1_voting
	scrutiny[shared.Animal] = cat2_voting

	expected := make(map[shared.Category]map[string]struct{})
	cat1_invalid := make(map[string]struct{})
	cat2_invalid := make(map[string]struct{})

	cat1_invalid["Perro"] = struct{}{}
	cat1_invalid["Oveja"] = struct{}{}
	cat2_invalid["Perro"] = struct{}{}

	expected[shared.Pais] = cat1_invalid
	expected[shared.Animal] = cat2_invalid

	got := getInvalidWords(scrutiny, 3)

	if !maps.EqualFunc(expected, got, maps.Equal) {
		t.Errorf("Expected %v, but got %v", expected, got)
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
