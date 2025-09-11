package ui

import "github.com/charmbracelet/lipgloss"

var (
	// appStyle is the base style for the application
	appStyle = lipgloss.NewStyle().Padding(1, 2)
	
	// titleStyle is the style for titles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)
			
	// inputStyle is the style for input fields
	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(0, 1)
			
	// resultStyle is the style for search results
	resultStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2)
			
	// cardStyle is the style for card details
	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2)
			
	// paginationStyle is the style for pagination information
	paginationStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Padding(1, 0)
			
	// helpStyle is the style for help text
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).Render
)