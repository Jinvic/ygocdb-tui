package ui

import (
	"ygocdb-tui/internal/api"
	"ygocdb-tui/internal/log"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
)

const (
	// PageSize is the number of items to display per page
	PageSize = 10
)

// Mode represents the current UI mode
type Mode int

const (
	// SearchMode is the mode for searching cards
	SearchMode Mode = iota
	// ResultMode is the mode for displaying search results
	ResultMode
	// CardMode is the mode for displaying card details
	CardMode
)

// Model represents the application state
type Model struct {
	textInput   textinput.Model
	results     []api.Card // All cached results
	currentPage int        // Current page index (0-based)
	totalPages  int        // Total number of pages
	card        *api.GetCardResponse
	selected    int
	err         error
	mode        Mode
	loading     bool
	apiClient   *api.Client
	query       string
	nextStart   int   // Next start position for API request
}

// NewModel creates a new UI model
func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "输入卡片名称或ID"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 40

	return Model{
		textInput:   ti,
		results:     []api.Card{},
		currentPage: 0,
		totalPages:  0,
		card:        nil,
		selected:    -1,
		err:         nil,
		mode:        SearchMode,
		loading:     false,
		apiClient:   api.NewClient(),
		query:       "",
		nextStart:   0,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	log.Info("Initializing UI model")
	return textinput.Blink
}