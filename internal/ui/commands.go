package ui

import (
	"strconv"
	"ygocdb-tui/internal/api"
	"ygocdb-tui/internal/log"
	tea "github.com/charmbracelet/bubbletea"
)

// Search cards command
func searchCards(query string, start int) tea.Cmd {
	log.Info("Initiating search command: query=%s, start=%d", query, start)
	
	return func() tea.Msg {
		log.Debug("Search command executing in background")
		
		// Check if query is a number (card ID) and we're on the first page
		if cardID, err := strconv.Atoi(query); err == nil && start == 0 {
			log.Info("Query identified as card ID: %d", cardID)
			
			// Query by card ID (only for first page)
			client := api.NewClient()
			log.Debug("Fetching card by ID: %d", cardID)
			card, err := client.GetCardByID(cardID)
			if err != nil {
				log.Error("Failed to fetch card by ID %d: %v", cardID, err)
				return searchErrorMsg{err}
			}
			
			log.Info("Successfully fetched card by ID: %d", card.ID)
			return searchByIDResultMsg{card}
		}
		
		// Search by name with pagination
		log.Info("Performing name search: query=%s, start=%d", query, start)
		client := api.NewClient()
		results, err := client.SearchCards(query, start)
		if err != nil {
			log.Error("Failed to perform search: %v", err)
			return searchErrorMsg{err}
		}
		
		log.Info("Search completed successfully, found %d results", len(results.Result))
		return searchResultMsg{results, query, start}
	}
}

// Get card by ID command
func getCardByID(id int) tea.Cmd {
	log.Info("Initiating get card by ID command: id=%d", id)
	
	return func() tea.Msg {
		log.Debug("Get card command executing in background")
		
		client := api.NewClient()
		log.Debug("Fetching card by ID: %d", id)
		card, err := client.GetCardByID(id)
		if err != nil {
			log.Error("Failed to fetch card by ID %d: %v", id, err)
			return searchErrorMsg{err}
		}
		
		log.Info("Successfully fetched card by ID: %d", card.ID)
		return cardResultMsg{card}
	}
}