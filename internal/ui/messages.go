package ui

import "ygocdb-tui/internal/api"

// Messages
type searchResultMsg struct {
	results *api.SearchResponse
	query   string
	start   int
}

type searchByIDResultMsg struct {
	card *api.GetCardResponse
}

type cardResultMsg struct {
	card *api.GetCardResponse
}

type searchErrorMsg struct {
	err error
}