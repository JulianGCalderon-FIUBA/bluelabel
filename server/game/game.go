package game

import (
	"bluelabel/shared"
	"encoding/gob"
	"fmt"
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

	gob.Register(shared.Round{})

	return Game{
		clients,
	}
}

// Starts the game loop. This function is blocking.
func (g *Game) PlayGame() {
	round := shared.Round{Character: randomRune()}
	err := g.broadcast(round)
	if err != nil {
		fmt.Printf("Could not broadcast message: %s", err)
	}
}

// Sends a gob-encoded structure to each client.
func (g *Game) broadcast(structure any) error {
	for _, client := range g.clients {
		err := client.send(&structure)
		if err != nil {
			return err
		}
	}

	return nil
}

// Returns a random rune, from 'a' to 'z'
func randomRune() rune {
	return rune('a' + rand.Intn('c'+1-'a'))
}
