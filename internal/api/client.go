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