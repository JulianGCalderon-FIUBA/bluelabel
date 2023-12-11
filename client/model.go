package main

import (
	"bluelabel/shared"
	"encoding/gob"
	"fmt"
	"net"

	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
)

// func read_from_connection(connection net.Conn) {
// 	for {
// 		tmp_buffer := make([]byte, 256)
// 		tmp_buffer_msg_len, err := connection.Read(tmp_buffer)

// 		if err != nil {
// 			break
// 		}

// 		fmt.Println(string(tmp_buffer[:tmp_buffer_msg_len]))
// 	}
// }

type gameState int

const (
	InMenu gameState = iota
	InMatch
	InRoomSelectorMenu
	WaitingForMorePlayers
)

type model struct {
	viewportWidth    int
	viewportHeight   int
	state            gameState
	menu             menu
	roomSelectorMenu roomSelectorMenu
	server           server
}

func initialModel() model {
	connection, err := net.Dial("tcp", "localhost:4000")

	gob.Register(shared.Round{})
	gob.Register(shared.StopRequest{})
	gob.Register(shared.StopNotify{})
	gob.Register(shared.Words{})
	gob.Register(shared.WordsValidation{})
	gob.Register(shared.Score{})
	gob.Register(shared.End{})

	if err == nil {
		return model{
			state:            InMenu,
			menu:             initialMenu(),
			roomSelectorMenu: initialRoomSelector(),
			server:           makeServer(connection),
		}
	} else {
		panic("No se pudo conectar con el servidor.")
	}
}

func (m *model) quickmatch() {
	m.server.send(shared.QuickMatch{})
	m.state = WaitingForMorePlayers

	//received, err := m.server.receive()

	//connection, err := net.Dial("tcp", "localhost:4000")

	//if err == nil {
	//go read_from_connection(connection)
	//go write_to_connection(connection)
	//}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewportWidth = msg.Width - 6
		m.viewportHeight = msg.Height - 6
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+q":
			return m, tea.Quit
		}
	}

	switch m.state {
	case InMenu:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "up", "down":
				m.menu.moveThroughOptions(msg)

			case "enter":
				switch m.menu.selectedOption {
				case int(QuickMatchOption):
					m.quickmatch()
				case int(RoomSelectorOption):
					m.state = InRoomSelectorMenu
				case int(QuitOption):
					return m, tea.Quit
				}
			}
		case tea.MouseMsg:
			if zone.Get(fmt.Sprint("option ", QuickMatchOption)).InBounds(msg) {
				if msg.Type == tea.MouseLeft {
					m.quickmatch()
				} else {
					m.menu.selectOption(int(QuickMatchOption))
				}

			} else if zone.Get(fmt.Sprint("option ", RoomSelectorOption)).InBounds(msg) {
				if msg.Type == tea.MouseLeft {
					m.state = InRoomSelectorMenu
				} else {
					m.menu.selectOption(int(RoomSelectorOption))
				}
			} else if zone.Get(fmt.Sprint("option ", QuitOption)).InBounds(msg) {
				if msg.Type == tea.MouseLeft {
					return m, tea.Quit
				} else {
					m.menu.selectOption(int(QuitOption))
				}
			}
		}
	case InRoomSelectorMenu:
		switch msg := msg.(type) {
		case tea.MouseMsg:
			if zone.Get("quitMenuOption").InBounds(msg) {
				if msg.Type == tea.MouseLeft {
					m.state = InMenu
				}
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	switch m.state {
	case InMenu:
		return zone.Scan(m.menu.show(m.viewportHeight, m.viewportWidth))
	case InRoomSelectorMenu:
		return zone.Scan(m.roomSelectorMenu.show(m.viewportHeight, m.viewportWidth))
	case WaitingForMorePlayers:
		return newWindow("Esperando jugadores...").show(m.viewportHeight, m.viewportWidth)
	}

	return ""
}
