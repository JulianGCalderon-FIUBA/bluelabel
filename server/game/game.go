package game

import (
	"net"
)

// Encapsulates the game data and behaviour
type Game struct {
	clients []client
}

// Creates a new game, for the given connections
func MakeGame(clientConnections ...net.Conn) Game {
	clients := make([]client, len(clientConnections))
	for i, c := range clientConnections {
		clients[i] = makeClient(c)
	}

	return Game{
		clients,
	}
}

// Starts the game loop. This function is blocking.
func (g *Game) PlayGame() {}

// Sends a gob-encoded structure to each client, as an interface.
func (g *Game) broadcast(structure any) error {
	for _, client := range g.clients {
		err := client.send(&structure)
		if err != nil {
			return err
		}
	}

	return nil
}
