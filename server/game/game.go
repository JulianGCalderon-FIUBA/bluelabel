package game

import (
	"net"
)

func PlayGame(clientConnections ...net.Conn) {
	game := makeGame(clientConnections...)
	game.playGame()
}

type game struct {
	clients []client
}

func makeGame(clientConnections ...net.Conn) game {
	clients := make([]client, len(clientConnections))
	for i, c := range clientConnections {
		clients[i] = makeClient(c)
	}

	return game{
		clients,
	}
}

func (g *game) playGame() {}
