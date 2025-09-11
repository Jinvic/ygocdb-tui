package ui

import (
	"strconv"
	"ygocdb-tui/internal/api"
	tea "github.com/charmbracelet/bubbletea"
)

// Search cards command
func searchCards(query string, start int) tea.Cmd {
	return func() tea.Msg {
		// Check if query is a number (card ID) and we're on the first page
		if cardID, err := strconv.Atoi(query); err == nil && start == 0 {
			// Query by card ID (only for first page)
			client := api.NewClient()
			card, err := client.GetCardByID(cardID)
			if err != nil {
				return searchErrorMsg{err}
			}
			return searchByIDResultMsg{card}
		}
		
		// Search by name with pagination
		client := api.NewClient()
		results, err := client.SearchCards(query, start)
		if err != nil {
			return searchErrorMsg{err}
		}
		return searchResultMsg{results, query, start}
	}
}

// Get card by ID command
func getCardByID(id int) tea.Cmd {
	return func() tea.Msg {
		client := api.NewClient()
		card, err := client.GetCardByID(id)
		if err != nil {
			return searchErrorMsg{err}
		}
		return cardResultMsg{card}
	}
}