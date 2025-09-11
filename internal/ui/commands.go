package ui

import (
	"strconv"
	"ygocdb-tui/internal/api"
	"ygocdb-tui/internal/log"
	tea "github.com/charmbracelet/bubbletea"
)

// searchCardsCmd creates a command to search for cards
func searchCardsCmd(query string, start int) tea.Cmd {
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
				return SearchErrorMsg{Err: err}
			}
			
			log.Info("Successfully fetched card by ID: %d", card.ID)
			return SearchByIDResultMsg{Card: card}
		}
		
		// Search by name with pagination
		log.Info("Performing name search: query=%s, start=%d", query, start)
		client := api.NewClient()
		results, err := client.SearchCards(query, start)
		if err != nil {
			log.Error("Failed to perform search: %v", err)
			return SearchErrorMsg{Err: err}
		}
		
		log.Info("Search completed successfully, found %d results", len(results.Result))
		return SearchResultMsg{Results: results, Query: query, Start: start}
	}
}

// getCardByIDCmd creates a command to get a card by ID
func getCardByIDCmd(id int) tea.Cmd {
	log.Info("Initiating get card by ID command: id=%d", id)
	
	return func() tea.Msg {
		log.Debug("Get card command executing in background")
		
		client := api.NewClient()
		log.Debug("Fetching card by ID: %d", id)
		card, err := client.GetCardByID(id)
		if err != nil {
			log.Error("Failed to fetch card by ID %d: %v", id, err)
			return SearchErrorMsg{Err: err}
		}
		
		log.Info("Successfully fetched card by ID: %d", card.ID)
		return CardResultMsg{Card: card}
	}
}

// nextPageCmd handles navigation to the next page
func (m *Model) nextPageCmd() tea.Cmd {
	log.Info("Checking if next page needs to be fetched")
	
	// Calculate the start index for the next page
	nextPageStart := (m.currentPage + 1) * PageSize
	
	// If we have enough cached results, just update the page
	if nextPageStart < len(m.results) {
		log.Debug("Next page is already cached, moving to page %d", m.currentPage+1)
		m.currentPage++
		m.selected = 0
		// After moving to next page, check if we need to auto-fetch more results to fill it
		return m.autoFetchNextPageCmd()
	}
	
	// If we don't have enough results, check if there are more pages available
	if m.nextStart > 0 {
		log.Info("Fetching next page from API, start=%d", m.nextStart)
		m.loading = true
		return searchCardsCmd(m.query, m.nextStart)
	}
	
	log.Debug("No more pages available")
	return nil
}

// prevPageCmd handles navigation to the previous page
func (m *Model) prevPageCmd() tea.Cmd {
	log.Info("Navigating to previous page, current page=%d", m.currentPage)
	
	if m.currentPage > 0 {
		m.currentPage--
		m.selected = 0
	}
	
	return nil
}

// autoFetchNextPageCmd automatically fetches more results to fill the current page
func (m *Model) autoFetchNextPageCmd() tea.Cmd {
	log.Info("Checking if we need to auto-fetch more results to fill current page")
	
	// Get current page results
	currentPageResults := m.getCurrentPageResults()
	
	// Calculate if current page is the last page based on current results
	expectedTotalPages := (len(m.results) + PageSize - 1) / PageSize
	isLastPage := m.currentPage == expectedTotalPages-1
	
	// If current page has less than PageSize items and there are more results available
	// and we're on the last page, then auto-fetch more results
	if len(currentPageResults) < PageSize && m.nextStart > 0 && isLastPage {
		log.Info("Current page has %d items (less than %d), is the last page, and more results are available. Auto-fetching next page.", 
			len(currentPageResults), PageSize)
		m.loading = true
		return searchCardsCmd(m.query, m.nextStart)
	}
	
	return nil
}