package main

import (
	"fmt"
	"net"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func read_from_connection(connection net.Conn) {
	for {
		tmp_buffer := make([]byte, 256)
		tmp_buffer_msg_len, err := connection.Read(tmp_buffer)

		if err != nil {
			break
		}

		fmt.Println(string(tmp_buffer[:tmp_buffer_msg_len]))
	}
}

// func write_to_connection(connection net.Conn) {
// 	input_scanner := bufio.NewScanner(os.Stdin)
// 	for {
// 		if input_scanner.Scan() {
// 			line := input_scanner.Text()
// 			io.WriteString(connection, line)
// 		}
// 	}
// }

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.state == InMenu {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {

			case "ctrl+c", "q":
				return m, tea.Quit

			case "up", "down":
				m.menu.selectOption(msg)

			case "enter":
				switch m.menu.selectedOption {
				case int(QuickMath):
					connection, err := net.Dial("tcp", "localhost:4000")

					if err == nil {
						go read_from_connection(connection)
						//go write_to_connection(connection)
					}
				case int(RoomSearcher):

				case int(Quit):
					return m, tea.Quit
				}
			}
		case tea.WindowSizeMsg:
			m.viewportWidth = msg.Width - 6
			m.viewportHeight = msg.Height - 5
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.state == InMenu {
		return m.menu.show(m.viewportHeight, m.viewportWidth)
	} else {
		return ""
	}
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Hubo un error inesperado: %v", err)
		os.Exit(1)
	}
}
