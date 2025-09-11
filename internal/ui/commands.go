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

// nextPage checks if we need to fetch more results from API
func (m *model) nextPage() tea.Cmd {
	log.Info("Checking if next page needs to be fetched")
	
	// Calculate the start index for the next page
	nextPageStart := (m.currentPage + 1) * PageSize
	
	// If we have enough cached results, just update the page
	if nextPageStart < len(m.results) {
		log.Debug("Next page is already cached, moving to page %d", m.currentPage+1)
		m.currentPage++
		m.selected = 0
		// After moving to next page, check if we need to auto-fetch more results to fill it
		return m.autoFetchNextPage()
	}
	
	// If we don't have enough results, check if there are more pages available
	if m.nextStart > 0 {
		log.Info("Fetching next page from API, start=%d", m.nextStart)
		m.loading = true
		return searchCards(m.query, m.nextStart)
	}
	
	log.Debug("No more pages available")
	return nil
}

// prevPage navigates to the previous page
func (m *model) prevPage() tea.Cmd {
	log.Info("Navigating to previous page, current page=%d", m.currentPage)
	
	if m.currentPage > 0 {
		m.currentPage--
		m.selected = 0
	}
	
	return nil
}

// autoFetchNextPage checks if we need to automatically fetch more results 
// to fill the current page when it has fewer than PageSize items
func (m *model) autoFetchNextPage() tea.Cmd {
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
		return searchCards(m.query, m.nextStart)
	}
	
	return nil
}