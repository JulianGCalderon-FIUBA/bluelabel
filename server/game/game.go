package game

import (
	"net"
)

func PlayGame(clientConnections ...net.Conn) {
	game := newGame(clientConnections...)
	game.playGame()
}

type game struct {
	clients []*client
}

func newGame(clientConnections ...net.Conn) *game {
	clients := make([]*client, len(clientConnections))
	for i, c := range clientConnections {
		clients[i] = newClient(c)
	}

	return &game{
		clients,
	}
}

func (g *game) playGame() {}
