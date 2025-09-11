package ui

import "ygocdb-tui/internal/api"

// SearchResultMsg represents a message containing search results
type SearchResultMsg struct {
	Results *api.SearchResponse
	Query   string
	Start   int
}

// SearchByIDResultMsg represents a message containing a card fetched by ID
type SearchByIDResultMsg struct {
	Card *api.GetCardResponse
}

// CardResultMsg represents a message containing card details
type CardResultMsg struct {
	Card *api.GetCardResponse
}

// SearchErrorMsg represents a message containing a search error
type SearchErrorMsg struct {
	Err error
}