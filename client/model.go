package main

import (
	"bluelabel/shared"
)

type GameState int

const (
	InMenu GameState = iota
	InMatch
)

type model struct {
	viewportWidth    int
	viewportHeight   int
	state            GameState
	menu             menu
	currentCharacter rune
	categories       []string
	scores           []int
	words            map[shared.Category]string
}

func initialModel() model {
	return model{
		state:            InMenu,
		menu:             initialMenu(),
		currentCharacter: 'A',
		categories:       make([]string, 0),
		scores:           make([]int, 0),
		words:            make(map[shared.Category]string, 0),
	}
}
