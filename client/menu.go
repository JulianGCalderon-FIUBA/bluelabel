package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
)

type MenuOption int

var (
	title = `████████╗██╗   ██╗████████╗████████╗██╗    ███████╗██████╗ ██╗   ██╗████████╗████████╗██╗
╚══██╔══╝██║   ██║╚══██╔══╝╚══██╔══╝██║    ██╔════╝██╔══██╗██║   ██║╚══██╔══╝╚══██╔══╝██║
   ██║   ██║   ██║   ██║      ██║   ██║    █████╗  ██████╔╝██║   ██║   ██║      ██║   ██║
   ██║   ██║   ██║   ██║      ██║   ██║    ██╔══╝  ██╔══██╗██║   ██║   ██║      ██║   ██║
   ██║   ╚██████╔╝   ██║      ██║   ██║    ██║     ██║  ██║╚██████╔╝   ██║      ██║   ██║
   ╚═╝    ╚═════╝    ╚═╝      ╚═╝   ╚═╝    ╚═╝     ╚═╝  ╚═╝ ╚═════╝    ╚═╝      ╚═╝   ╚═╝` + "\n\n\n"
)

const (
	QuickMatchOption MenuOption = iota
	RoomSelectorOption
	QuitOption
)

type menu struct {
	options        []string
	selectedOption int
}

func initialMenu() menu {
	return menu{
		options:        []string{"Partida Rápida", "Buscar Sala", "Salir"},
		selectedOption: 0,
	}
}

func (m *menu) moveThroughOptions(msg tea.KeyMsg) {
	switch msg.String() {
	case "up":
		m.selectOption(m.selectedOption - 1)
	case "down":
		m.selectOption(m.selectedOption + 1)
	}
}

func (m *menu) selectOption(option int) {
	i := int(option)

	if i >= 0 && i < len(m.options) {
		m.selectedOption = i
	} else if i < 0 {
		m.selectedOption = len(m.options) - 1
	} else {
		m.selectedOption = 0
	}
}

func (m menu) show(viewportHeight int, viewportWidth int) string {
	s := title

	for i, option := range m.options {
		zoneMarkName := fmt.Sprint("option ", i)

		if i == m.selectedOption {
			s += zone.Mark(zoneMarkName, selectedOptionStyle().Render(option))
			s += "\n"
		} else {
			s += zone.Mark(zoneMarkName, optionStyle().Render(option))
			s += "\n"
		}
	}

	return mainBlockStyle().
		Width(viewportWidth).
		Height(viewportHeight).
		Render(s)
}
