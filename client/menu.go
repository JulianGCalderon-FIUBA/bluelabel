package main

import tea "github.com/charmbracelet/bubbletea"

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
	QuickMath MenuOption = iota
	RoomSearcher
	Quit
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

func (m *menu) selectOption(msg tea.KeyMsg) {
	switch msg.String() {
	case "up":
		if m.selectedOption == 0 {
			m.selectedOption = len(m.options) - 1
		} else {
			m.selectedOption--
		}
	case "down":
		if m.selectedOption == len(m.options)-1 {
			m.selectedOption = 0
		} else {
			m.selectedOption++
		}
	}
}

func (m menu) show(viewportHeight int, viewportWidth int) string {
	s := title

	for i, option := range m.options {
		if i == m.selectedOption {
			s += selectedOptionStyle().Render(option) + "\n\n"
		} else {
			s += option + "\n\n"
		}
	}

	return mainBlockStyle().
		Width(viewportWidth).
		Height(viewportHeight).
		Render(s)
}
