package game

import (
	"net"
)

type Game struct {
	clients []client
}

func MakeGame(clientConnections ...net.Conn) Game {
	clients := make([]client, len(clientConnections))
	for i, c := range clientConnections {
		clients[i] = makeClient(c)
	}

	return Game{
		clients,
	}
}

func (g *Game) PlayGame() {}
