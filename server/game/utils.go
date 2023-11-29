package game

import "math/rand"

// Returns a random rune, from 'a' to 'z'
func randomRune() rune {
	return rune('a' + rand.Intn('c'+1-'a'))
}

// Represents a concrete message received from a specific client
type MessageWithId[T any] struct {
	id  int
	msg T
}
