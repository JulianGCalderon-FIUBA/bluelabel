package matcher_test

import (
	"bluelabel/server/matcher"
	"net"
	"testing"
)

func createPipeArray(lobbySize int) [][2]net.Conn {
	pipes := make([][2]net.Conn, lobbySize)
	for i := range pipes {
		pipes[i][0], pipes[i][1] = net.Pipe()
	}
	return pipes
}

func TestLobbySizeIsCorrect(t *testing.T) {
	lobbySize := 8
	pipes := createPipeArray(lobbySize)

	ch := make(chan int)
	starter := func(connections ...net.Conn) {
		ch <- len(connections)
	}
	matcher := matcher.NewMatcher(lobbySize, matcher.MatchStarterFunc(starter))
	for _, c := range pipes {
		matcher.Match(c[1])
	}

	actualLobbySize := <-ch
	if lobbySize != actualLobbySize {
		t.Errorf("Expected lobby size of %v, but got %v", lobbySize, actualLobbySize)
	}
}
