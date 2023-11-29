package game

import (
	"bluelabel/shared"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"time"
)

const stopTimeoutDuration = 1 * time.Second

// Encapsulates the game data and behaviour
type Game struct {
	clients []client

	// TODO: Probablemente sea mejor guardar estos campos en una estructura
	// auxiliar de "ronda", para que sea mas ordenado
	character rune
	words     map[int]map[shared.Category]string
}

// Creates a new game, for the given connections
func MakeGame(clientConnections ...net.Conn) Game {
	clients := make([]client, len(clientConnections))
	for i, c := range clientConnections {
		clients[i] = makeClient(c)
	}

	gob.Register(shared.Round{})
	gob.Register(shared.StopRequest{})
	gob.Register(shared.StopNotify{})
	gob.Register(shared.Words{})
	gob.Register(shared.WordsValidation{})
	gob.Register(shared.Score{})
	gob.Register(shared.End{})

	return Game{
		clients: clients,
		words:   make(map[int]map[shared.Category]string),
	}
}

// Starts the game loop. This function is blocking.
func (g *Game) PlayGame() {
	err := g.startRound()
	if err != nil {
		fmt.Printf("Could not broadcast new round: %s", err)
		return
	}

	err = g.waitStop()
	if err != nil {
		fmt.Printf("Could not wait for stop: %s", err)
	}

	g.broadcastWords()

	g.waitValidation()

	g.broadcastScore()

	g.broadcastEnd()
}

// Sends a round message to every client
func (g *Game) startRound() error {
	round := shared.Round{Character: randomRune()}
	g.character = round.Character
	return g.broadcast(round)
}

// Waits for all clients to send a StopRequest message
// Sends a StopNotify to every other client after the first received StopRequest.
// If the clients don't send a StopRequest after certain timeout time, the
// function returns prematurely.
func (g *Game) waitStop() error {
	messagesWithIds := make(chan MessageWithId[shared.StopRequest])
	for i, c := range g.clients {
		go func(i int, c client) {
			stopRequest, err := receiveConcrete[shared.StopRequest](c)
			if err != nil {
				log.Printf("Could not receive from client: %s", err)
				return
			}

			messagesWithIds <- MessageWithId[shared.StopRequest]{
				i, stopRequest,
			}
			return

		}(i, c)
	}

	firstStop := <-messagesWithIds
	g.words[firstStop.id] = firstStop.msg.Words

	g.broadcastAllBut(shared.StopNotify{}, firstStop.id)

	timeout := time.NewTimer(stopTimeoutDuration)
	for len(g.words) < len(g.clients) {
		select {
		case stop := <-messagesWithIds:
			g.words[stop.id] = stop.msg.Words
		case <-timeout.C:
			return nil
		}
	}

	return nil
}

func (g *Game) broadcastWords() error { return nil }
func (g *Game) waitValidation() error { return nil }
func (g *Game) broadcastScore() error { return nil }
func (g *Game) broadcastEnd() error   { return nil }

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

// Sends a gob-encoded structure to every client, except one.
func (g *Game) broadcastAllBut(message any, first_id int) {
	for i, c := range g.clients {
		if i != first_id {
			c.send(message)
		}
	}
}
