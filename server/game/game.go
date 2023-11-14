package game

import (
	"bluelabel/shared"
	"math/rand"
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
func (g *Game) PlayGame() {
	round := shared.Round{Character: randomRune()}
	g.broadcast(&round)
}

// Sends a gob-encoded structure to each client.
func (g *Game) broadcast(structure any) error {
	for _, client := range g.clients {
		err := client.send(structure)
		if err != nil {
			return err
		}
	}

	return nil
}

func randomRune() rune {
	return rune('a' + rand.Intn('c'+1-'a'))
}
