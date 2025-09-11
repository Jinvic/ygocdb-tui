package ui

import "github.com/charmbracelet/lipgloss"

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)
	
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)
			
	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(0, 1)
			
	resultStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2)
			
	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2)
			
	paginationStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Padding(1, 0)
			
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).Render
)