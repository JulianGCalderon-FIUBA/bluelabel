package game

import (
	"bluelabel/shared"
	"math/rand"
)

// Returns a random rune, from 'a' to 'z'
func randomRune() rune {
	return rune('a' + rand.Intn('c'+1-'a'))
}

// Represents a concrete message received from a specific client
type MessageFromClient[T any] struct {
	id  int
	msg T
}

// Given a map of words (as a map by category) by client, construct a map of words
// set (as a list) by category.
func buildWordListByCategory(words map[int]map[shared.Category]string) map[shared.Category][]string {
	wordsByCategory := make(map[shared.Category]map[string]struct{})
	for _, clientWords := range words {
		for category, word := range clientWords {
			if wordsByCategory[category] == nil {
				wordsByCategory[category] = make(map[string]struct{})
			}

			wordsByCategory[category][word] = struct{}{}
		}
	}

	wordListByCategory := make(map[shared.Category][]string)
	for category, words := range wordsByCategory {
		wordListByCategory[category] = make([]string, 0, len(words))
		for word := range words {
			wordListByCategory[category] = append(wordListByCategory[category], word)
		}
	}

	return wordListByCategory
}
