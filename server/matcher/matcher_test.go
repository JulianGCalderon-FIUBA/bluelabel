package matcher_test

import (
	"bluelabel/server/matcher"
	"io"
	"net"
	"slices"
	"sync"
	"testing"
)

const (
	dummyMessage   = "dummy"
	maxLobbySize   = 2
	maxLobbyAmount = 5
)

func createPipeArray(lobbySize int) [][2]net.Conn {
	pipes := make([][2]net.Conn, lobbySize)
	for i := range pipes {
		pipes[i][0], pipes[i][1] = net.Pipe()
	}
	return pipes
}

// This tests creates pipes for a number of lobbies with fixed size and matches them, with a
// game starter function that asserts the expected lobby size.
//
// This proves that all lobbies have the requested size.
func TestLobbySizeIsCorrect(t *testing.T) {
	for lobbySize := 2; lobbySize <= maxLobbySize; lobbySize++ {
		pipes := createPipeArray(lobbySize * maxLobbyAmount)

		wg := sync.WaitGroup{}
		wg.Add(maxLobbyAmount)

		starter := func(connections ...net.Conn) {
			if lobbySize != len(connections) {
				t.Errorf("Expected lobby size of %v, but got %v", lobbySize, len(connections))
			}

			wg.Done()
		}

		matcher := matcher.NewMatcher(lobbySize, matcher.MatchStarterFunc(starter))

		for _, c := range pipes {
			go matcher.Match(c[1])
		}

		wg.Wait()
	}
}

// This tests creates pipes for a number of lobbies with fixed size and matches them,
// with a game starter function that broadcast a dummy message to all connections.
// Each remote connection should read the correct message exactly once.
//
// This proves that a connections is not involved in more than one match.
func TestLobbyConnectsAllUsers(t *testing.T) {
	for lobbySize := 2; lobbySize <= maxLobbySize; lobbySize++ {
		pipes := createPipeArray(lobbySize * maxLobbyAmount)

		starter := func(connections ...net.Conn) {
			for _, c := range connections {
				go func(c net.Conn) {
					_, err := io.WriteString(c, dummyMessage)
					if err != nil {
						t.Errorf("Two matches where started with the same connection")
					}
					c.Close()
				}(c)
			}
		}

		matcher := matcher.NewMatcher(lobbySize, matcher.MatchStarterFunc(starter))

		for _, c := range pipes {
			go matcher.Match(c[1])
		}

		for _, c := range pipes {
			buf, err := io.ReadAll(c[0])
			if err != nil {
				t.Errorf("Could not read from pipe")
			}

			if !slices.Equal(buf, []byte(dummyMessage)) {
				t.Errorf("Expected %s from remote pipe, but read %s", dummyMessage, buf)
			}
		}
	}
}
