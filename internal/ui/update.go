package ui

import (
	"fmt"
	"ygocdb-tui/internal/api"
	"ygocdb-tui/internal/log"
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles UI updates
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	
	log.Debug("Processing message of type: %T", msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		log.Debug("Processing key message: %v", msg)
		
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			if m.mode == SearchMode {
				log.Info("Received exit key, quitting application")
				return m, tea.Quit
			} else if m.mode == ResultMode || m.mode == CardMode {
				log.Info("Returning to search mode")
				// Go back to search mode
				m.mode = SearchMode
				m.results = []api.Card{}
				m.card = nil
				m.selected = -1
				m.textInput.Focus()
				m.currentPage = 0
				m.totalPages = 0
				m.nextStart = 0
				return m, nil
			}

		case tea.KeyEnter:
			if m.mode == SearchMode && !m.loading {
				// Search for cards
				query := m.textInput.Value()
				if query != "" {
					log.Info("Initiating search for query: %s", query)
					m.query = query
					m.currentPage = 0
					m.loading = true
					m.textInput.Blur()
					return m, searchCardsCmd(query, 0)
				}
			} else if m.mode == ResultMode && len(m.results) > 0 {
				// View selected card
				// Calculate the actual index in the full results array
				actualIndex := m.currentPage*PageSize + m.selected
				if actualIndex >= 0 && actualIndex < len(m.results) {
					log.Info("Viewing card details for card ID: %d", m.results[actualIndex].ID)
					m.loading = true
					return m, getCardByIDCmd(m.results[actualIndex].ID)
				}
			} else if m.mode == CardMode {
				// Back to results
				log.Info("Returning to search results")
				m.mode = ResultMode
				m.card = nil
				m.loading = false
				return m, nil
			}

		case tea.KeyUp:
			if m.mode == ResultMode && len(m.getCurrentPageResults()) > 0 {
				m.selected--
				if m.selected < 0 {
					m.selected = len(m.getCurrentPageResults()) - 1
				}
				log.Debug("Selected item changed to index: %d", m.selected)
			}
			return m, nil

		case tea.KeyDown:
			if m.mode == ResultMode && len(m.getCurrentPageResults()) > 0 {
				m.selected++
				if m.selected >= len(m.getCurrentPageResults()) {
					m.selected = 0
				}
				log.Debug("Selected item changed to index: %d", m.selected)
			}
			return m, nil

		case tea.KeyRight:
			// Next page
			if m.mode == ResultMode && !m.loading {
				log.Info("Navigating to next page, current page=%d", m.currentPage)
				return m, m.nextPageCmd()
			}
			return m, nil

		case tea.KeyLeft:
			// Previous page
			if m.mode == ResultMode && !m.loading && m.currentPage > 0 {
				log.Info("Navigating to previous page, current page=%d", m.currentPage)
				return m, m.prevPageCmd()
			}
			return m, nil
		}

	case SearchResultMsg:
		log.Info("Received search results message, found %d results", len(msg.Results.Result))
		m.loading = false
		m.mode = ResultMode
		// Append new results to cached results
		m.results = append(m.results, msg.Results.Result...)
		// Update pagination info
		m.nextStart = msg.Results.Next
		m.totalPages = (len(m.results) + PageSize - 1) / PageSize
		// Reset selection
		m.selected = 0
		if len(m.results) == 0 {
			m.err = fmt.Errorf("未找到相关卡片")
			log.Warn("No results found for search")
		} else {
			// Check if we need to auto-fetch more results to fill the current page
			return m, m.autoFetchNextPageCmd()
		}
		return m, nil

	case SearchByIDResultMsg:
		log.Info("Received card by ID result message, card ID: %d", msg.Card.ID)
		m.loading = false
		m.mode = CardMode
		m.card = msg.Card
		return m, nil

	case CardResultMsg:
		log.Info("Received card result message, card ID: %d", msg.Card.ID)
		m.loading = false
		m.mode = CardMode
		m.card = msg.Card
		return m, nil

	case SearchErrorMsg:
		log.Error("Received search error message: %v", msg.Err)
		m.loading = false
		m.err = msg.Err
		m.textInput.Focus()
		return m, nil

	// Handle input changes
	case tea.WindowSizeMsg:
		log.Debug("Window size changed: width=%d, height=%d", msg.Width, msg.Height)
		// Handle window resizing if needed
	}

	// Update text input
	log.Debug("Updating text input")
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// getCurrentPageResults returns the results for the current page
func (m *Model) getCurrentPageResults() []api.Card {
	start := m.currentPage * PageSize
	end := start + PageSize
	
	// Ensure end doesn't exceed the total number of results
	if end > len(m.results) {
		end = len(m.results)
	}
	
	// If start is beyond the results, return empty slice
	if start >= len(m.results) {
		return []api.Card{}
	}
	
	return m.results[start:end]
}