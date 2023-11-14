package main

import (
	"bluelabel/shared"
	"bufio"
	"fmt"
	"io"
	"net"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

var selectedOption = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4"))

type GameState int

const (
	InMenu GameState = iota
	InMatch
)

type model struct {
	state              GameState
	menuOptions        []string
	selectedMenuOption int
	currentCharacter   rune
	categories         []string
	scores             []int
	words              map[shared.Category]string
}

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

func write_to_connection(connection net.Conn) {
	input_scanner := bufio.NewScanner(os.Stdin)

	for {
		if input_scanner.Scan() {
			line := input_scanner.Text()
			io.WriteString(connection, line)
		}
	}
}

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

			case "up":
				if m.selectedMenuOption == 0 {
					m.selectedMenuOption = len(m.menuOptions) - 1
				} else {
					m.selectedMenuOption--
				}
			case "down":
				if m.selectedMenuOption == len(m.menuOptions)-1 {
					m.selectedMenuOption = 0
				} else {
					m.selectedMenuOption++
				}

			case "enter":
				if m.selectedMenuOption == 0 {
					connection, err := net.Dial("tcp", "localhost:4000")

					if err == nil {
						go read_from_connection(connection)
						go write_to_connection(connection)
					}
				} else if m.selectedMenuOption == 1 {
					return m, tea.Quit
				}
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "¡Tutti Frutti!\n\n"

	for i, option := range m.menuOptions {
		if i == m.selectedMenuOption {
			s += selectedOption.Render(option) + "\n"
		} else {
			s += option + "\n"
		}
	}

	return s
}

func main() {
	p := tea.NewProgram(model{
		state: InMenu,
		menuOptions: []string{
			"Jugar",
			"Salir",
		},
		selectedMenuOption: 0,
		currentCharacter:   'A',
		categories:         make([]string, 0),
		scores:             make([]int, 0),
		words:              make(map[shared.Category]string, 0),
	})
	if _, err := p.Run(); err != nil {
		fmt.Printf("Hubo un error inesperado: %v", err)
		os.Exit(1)
	}
}
