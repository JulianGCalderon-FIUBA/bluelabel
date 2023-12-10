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
	character    rune
	words        map[int]map[shared.Category]string
	invalidWords map[shared.Category]map[string]struct{}
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

	err = g.broadcastWords()
	if err != nil {
		fmt.Printf("Could not broadcast words: %s", err)
	}

	err = g.waitValidation()
	if err != nil {
		fmt.Printf("Could not receive validations: %s", err)
	}

	// g.broadcastScore()
	//
	// g.broadcastEnd()
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
//
// TODO: Esta función es batante larga y confusa, pero todavía no la separe ya
// que voy a ver primero que otras funcionalidades comparte con el resto de
// etapas de la ronda.
func (g *Game) waitStop() error {
	stopsFromClients := make(chan MessageFromClient[shared.StopRequest])
	for i, c := range g.clients {
		go func(i int, c client) {
			stopRequest, err := receiveConcrete[shared.StopRequest](c)
			if err != nil {
				// FIX: ¿Cómo se puede manejar este error?
				log.Printf("Could not receive from client: %s", err)
				return
			}

			stopsFromClients <- MessageFromClient[shared.StopRequest]{
				i, stopRequest,
			}
			return

		}(i, c)
	}

	firstStop := <-stopsFromClients
	g.words[firstStop.id] = firstStop.msg.Words

	err := g.broadcastAllBut(shared.StopNotify{}, firstStop.id)
	if err != nil {
		return err
	}

	timeout := time.NewTimer(stopTimeoutDuration)
	for len(g.words) < len(g.clients) {
		select {
		case stop := <-stopsFromClients:
			g.words[stop.id] = stop.msg.Words
		case <-timeout.C:
			return nil
		}
	}

	return nil
}

func (g *Game) broadcastWords() error {
	wordListByCategory := buildWordListByCategory(g.words)

	words := shared.Words{Words: wordListByCategory}

	return g.broadcast(words)
}

func (g *Game) waitValidation() error {
	validationFromClients := make(chan shared.WordsValidation)
	for i, c := range g.clients {
		go func(i int, c client) {
			validation, err := receiveConcrete[shared.WordsValidation](c)
			if err != nil {
				// FIX: ¿Cómo se puede manejar este error?
				log.Printf("Could not receive from client: %s", err)
				return
			}

			validationFromClients <- validation
			return

		}(i, c)
	}

	scrutiny := make(map[shared.Category]map[string]int)

	timeout := time.NewTimer(stopTimeoutDuration)
	for i := 0; i < len(g.clients); i++ {
		select {
		case validation := <-validationFromClients:
			for category, words := range validation.Invalid {
				for _, word := range words {
					scrutiny[category][word] += 1
				}
			}
		case <-timeout.C:
			break
		}
	}

	g.invalidWords = getInvalidWords(scrutiny, len(g.clients)/2+1)

	return nil
}

func getInvalidWords(scrutiny map[shared.Category]map[string]int, necessaryVotes int) map[shared.Category]map[string]struct{} {
	invalidWords := make(map[shared.Category]map[string]struct{})

	for category, words := range scrutiny {
		invalidWords[category] = make(map[string]struct{})
		for word, votes := range words {
			if votes >= necessaryVotes {
				invalidWords[category][word] = struct{}{}
			}
		}
	}

	return invalidWords
}

func (g *Game) broadcastScore() error { return nil }
func (g *Game) broadcastEnd() error   { return nil }

// Sends a gob-encoded structure to each client.
// FIX: ¿Es correcto devolver ante el primer error?
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
func (g *Game) broadcastAllBut(message any, first_id int) error {
	for i, client := range g.clients {
		if i != first_id {
			err := client.send(message)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
