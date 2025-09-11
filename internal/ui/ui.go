package ui

import (
	"ygocdb-tui/internal/log"
	tea "github.com/charmbracelet/bubbletea"
)

// Start initializes and starts the TUI application
func Start() error {
	log.Info("Starting TUI application")
	p := tea.NewProgram(NewModel())
	_, err := p.Run()
	
	if err != nil {
		log.Error("TUI application error: %v", err)
	} else {
		log.Info("TUI application exited normally")
	}
	
	return err
}