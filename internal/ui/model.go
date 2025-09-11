package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	"ygocdb-tui/internal/api"
)

const (
	// PageSize is the number of items to display per page
	PageSize = 10
)

type model struct {
	textInput   textinput.Model
	results     []api.Card // All cached results
	currentPage int        // Current page index (0-based)
	totalPages  int        // Total number of pages
	card        *api.GetCardResponse
	selected    int
	err         error
	mode        mode
	loading     bool
	apiClient   *api.Client
	query       string
	nextStart   int   // Next start position for API request
	pageHistory []int // Record page history for navigation
}

type mode int

const (
	searchMode mode = iota
	resultMode
	cardMode
)

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "输入卡片名称或ID"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 40

	return model{
		textInput:   ti,
		results:     []api.Card{},
		currentPage: 0,
		totalPages:  0,
		card:        nil,
		selected:    -1,
		err:         nil,
		mode:        searchMode,
		loading:     false,
		apiClient:   api.NewClient(),
		query:       "",
		nextStart:   0,
		pageHistory: make([]int, 0),
	}
}