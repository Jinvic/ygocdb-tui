package ui

import tea "github.com/charmbracelet/bubbletea"

// Start initializes and starts the TUI application
func Start() error {
	p := tea.NewProgram(initialModel())
	_, err := p.Run()
	return err
}