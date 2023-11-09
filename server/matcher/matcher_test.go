package matcher_test

import (
	"bluelabel/server/matcher"
	"io"
	"net"
	"slices"
	"testing"
)

const (
	dummyMessage = "dummy"
	maxLobbySize = 2
)

func broadcastDummyMessage(connections ...net.Conn) {
	for _, c := range connections {
		go func(c net.Conn) {
			io.WriteString(c, dummyMessage)
			c.Close()
		}(c)
	}
}

func TestLobbyConnectsUsers(t *testing.T) {
	for lobbySize := 2; lobbySize <= maxLobbySize; lobbySize++ {
		pipes := make([][2]net.Conn, lobbySize)
		for i := range pipes {
			pipes[i][0], pipes[i][1] = net.Pipe()
		}

		starter := func(connections ...net.Conn) {
			if lobbySize != len(connections) {
				t.Errorf("Expected lobby size of %v, but got %v", lobbySize, len(connections))
			}

			broadcastDummyMessage(connections...)
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
