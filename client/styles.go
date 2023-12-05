package main

import (
	lipgloss "github.com/charmbracelet/lipgloss"
)

const (
	mainMargin = 2
)

func selectedOptionStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Width(23).
		Height(3).
		Bold(true).
		Underline(true).
		Background(lipgloss.Color("#7D56F4")).
		Align(lipgloss.Center, lipgloss.Center).
		Padding(0, 3).
		MarginBottom(1)
}

func optionStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#444444")).
		Width(21).
		Height(1).
		Padding(0, 3).
		MarginBottom(1).
		Align(lipgloss.Center, lipgloss.Center)
}

func mainBlockStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("99")).
		Margin(mainMargin).
		Align(lipgloss.Center, lipgloss.Center)
}

func roomsListStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Inherit(mainBlockStyle()).
		Margin(mainMargin).
		UnsetAlign()
}
