package main

import (
	lipgloss "github.com/charmbracelet/lipgloss"
)

func selectedOptionStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Underline(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 3)
}

func mainBlockStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("99")).
		Margin(2).
		Align(lipgloss.Center, lipgloss.Center)
}
