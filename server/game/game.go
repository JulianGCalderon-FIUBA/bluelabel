package game

import (
	"bluelabel/shared"
	"encoding/gob"
	"fmt"
	"math/rand"
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
		go clients[i].loop()
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
func (g *Game) waitStop() error {
	stops := buildStopArray(g.clients...)

	indexed_stops := mergeChannels(stops...)

	first_stop_id := g.waitOneStop(indexed_stops)

	g.broadcastAllBut(shared.StopNotify{}, first_stop_id)

	g.waitAllStop(indexed_stops, stopTimeoutDuration)

	return nil
}

// Builds an array with the stop channel for every client
func buildStopArray(clients ...client) []chan shared.StopRequest {
	stops := make([]chan shared.StopRequest, len(clients))
	for i := range stops {
		stops[i] = clients[i].stops
	}
	return stops
}

// Wait for one stop request, and return the id of the client
// who send it.
func (g *Game) waitOneStop(stops chan indexedMessage[shared.StopRequest]) int {
	stop := <-stops
	g.words[stop.id] = stop.msg.Words

	return stop.id
}

// Wait for all the remaining stop requests
func (g *Game) waitAllStop(stops chan indexedMessage[shared.StopRequest], timeoutDuration time.Duration) {
	timeout := time.NewTimer(timeoutDuration)
	for len(g.words) < len(g.clients) {
		select {
		case stop := <-stops:
			g.words[stop.id] = stop.msg.Words
		case <-timeout.C:
			return
		}
	}
}

func (g *Game) broadcastWords() error { return nil }
func (g *Game) waitValidation() error { return nil }
func (g *Game) broadcastScore() error { return nil }
func (g *Game) broadcastEnd() error   { return nil }

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

// Sends a gob-encoded structure to every client, except one.
func (g *Game) broadcastAllBut(message any, first_id int) {
	for i, c := range g.clients {
		if i != first_id {
			c.send(&message)
		}
	}
}

// Returns a random rune, from 'a' to 'z'
func randomRune() rune {
	return rune('a' + rand.Intn('c'+1-'a'))
}

// Represents a message received from a merge channel
type indexedMessage[T any] struct {
	id  int
	msg T
}

// Merges an array of channels into a single merge channel, receiving
// indexed messages.
// The id of the message received corresponds to the index of the channel which
// sent it.
func mergeChannels[T any](chs ...chan T) chan indexedMessage[T] {
	merge := make(chan indexedMessage[T])
	for i, c := range chs {
		go func(i int, c chan T) {
			merge <- indexedMessage[T]{
				i, <-c,
			}
		}(i, c)
	}

	return merge
}
