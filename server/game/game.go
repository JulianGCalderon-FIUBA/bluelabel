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

	// TODO: Probablemente sea mejor guardar estos campos en una estructura
	// auxiliar de "ronda", para que sea mas ordenado
	character rune
	words     []map[shared.Category]string
}

// Creates a new game, for the given connections
func MakeGame(clientConnections ...net.Conn) Game {
	clients := make([]client, len(clientConnections))
	for i, c := range clientConnections {
		clients[i] = makeClient(c)
	}

	gob.Register(shared.Round{})

	return Game{
		clients: clients,
	}
}

// Starts the game loop. This function is blocking.
func (g *Game) PlayGame() {
	err := g.startRound()
	if err != nil {
		fmt.Printf("Could not broadcast new round: %s", err)
		return
	}

	g.waitStop()

	g.broadcastWords()

	g.waitVoting()

	g.broadcastResults()
}

func (g *Game) startRound() error {
	round := shared.Round{Character: randomRune()}
	g.character = round.Character
	return g.broadcast(round)
}

func (g *Game) waitStop() error         { return nil }
func (g *Game) broadcastWords() error   { return nil }
func (g *Game) waitVoting() error       { return nil }
func (g *Game) broadcastResults() error { return nil }

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
