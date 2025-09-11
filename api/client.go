package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	BaseURL = "https://ygocdb.com"
)

// Card represents a Yu-Gi-Oh! card
type Card struct {
	CID    int    `json:"cid"`
	ID     int    `json:"id"`
	CnName string `json:"cn_name"`
	ScName string `json:"sc_name"`
	MdName string `json:"md_name"`
	NwbbsN string `json:"nwbbs_n"`
	CnocgN string `json:"cnocg_n"`
	JpRuby string `json:"jp_ruby"`
	JpName string `json:"jp_name"`
	EnName string `json:"en_name"`
	Text   Text   `json:"text"`
	Data   Data   `json:"data"`
}

// Text represents card text information
type Text struct {
	Name  string `json:"name"`
	Types string `json:"types"`
	PDesc string `json:"pdesc"`
	Desc  string `json:"desc"`
}

// Data represents card data information
type Data struct {
	OT      int `json:"ot"`
	Setcode int `json:"setcode"`
	Type    int `json:"type"`
	Atk     int `json:"atk"`
	Def     int `json:"def"`
	Level   int `json:"level"`
	Race    int `json:"race"`
	Attrib  int `json:"attribute"`
}

// SearchResponse represents the response from search API
type SearchResponse struct {
	Result []Card `json:"result"`
	Next   int    `json:"next"`
}

// GetCardResponse represents the response from get card API
type GetCardResponse struct {
	ID   int  `json:"id"`
	Data Data `json:"data"`
	Text Text `json:"text"`
}

// Client represents the API client
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new API client
func NewClient() *Client {
	return &Client{
		baseURL:    BaseURL,
		httpClient: &http.Client{},
	}
}

// SearchCards searches for cards by query with pagination
func (c *Client) SearchCards(query string, start int) (*SearchResponse, error) {
	// URL encode the query
	encodedQuery := url.QueryEscape(query)

	// Construct the URL with start parameter for pagination
	apiURL := fmt.Sprintf("%s/api/v0/?search=%s&start=%d", c.baseURL, encodedQuery, start)

	// Make the request
	resp, err := c.httpClient.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse JSON response
	var searchResp SearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &searchResp, nil
}

// GetCardByID gets a card by its ID
func (c *Client) GetCardByID(cardID int) (*GetCardResponse, error) {
	// Construct the URL
	apiURL := fmt.Sprintf("%s/api/v0/card/%d", c.baseURL, cardID)

	// Make the request
	resp, err := c.httpClient.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse JSON response
	var cardResp GetCardResponse
	if err := json.Unmarshal(body, &cardResp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &cardResp, nil
}
