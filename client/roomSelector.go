package main

import (
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

type roomSelectorMenu struct {
	rooms []string
}

func initialRoomSelector() roomSelectorMenu {
	return roomSelectorMenu{
		rooms: []string{},
	}
}

func (roomSelectorMenu roomSelectorMenu) show(viewportHeight int, viewportWidth int) string {
	s := lipgloss.JoinHorizontal(
		lipgloss.Top,
		roomsListStyle().Width((viewportWidth/8)*7).Height(viewportHeight).Render(),
		zone.Mark("quitMenuOption", optionStyle().Width(viewportWidth/8-2).MarginTop(2).MarginRight(2).Render("Volver")),
	)

	return s
}
