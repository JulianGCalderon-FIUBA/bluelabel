package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
)

func write_to_connection(connection net.Conn) {
	input_scanner := bufio.NewScanner(os.Stdin)
	for {
		if input_scanner.Scan() {
			line := input_scanner.Text()
			io.WriteString(connection, line)
		}
	}
}

func main() {
	zone.NewGlobal()
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseAllMotion())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Hubo un error inesperado: %v", err)
		os.Exit(1)
	}
}
