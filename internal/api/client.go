package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"ygocdb-tui/internal/log"
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
	log.Debug("Creating new API client")
	
	client := &Client{
		baseURL:    BaseURL,
		httpClient: &http.Client{},
	}
	
	log.Debug("API client created with baseURL: %s", client.baseURL)
	return client
}

// SearchCards searches for cards by query with pagination
func (c *Client) SearchCards(query string, start int) (*SearchResponse, error) {
	log.Info("Searching cards with query: %s, start: %d", query, start)
	
	// URL encode the query
	encodedQuery := url.QueryEscape(query)
	log.Debug("Encoded query: %s", encodedQuery)

	// Construct the URL with start parameter for pagination
	apiURL := fmt.Sprintf("%s/api/v0/?search=%s&start=%d", c.baseURL, encodedQuery, start)
	log.Debug("API URL: %s", apiURL)

	// Make the request
	log.Debug("Making HTTP request to API")
	resp, err := c.httpClient.Get(apiURL)
	if err != nil {
		log.Error("Failed to make HTTP request: %v", err)
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	
	log.Debug("Received HTTP response with status code: %d", resp.StatusCode)

	// Check status code
	if resp.StatusCode != http.StatusOK {
		log.Error("API returned non-OK status code: %d", resp.StatusCode)
		return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	// Read response body
	log.Debug("Reading response body")
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Failed to read response body: %v", err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	log.Debug("Response body read, size: %d bytes", len(body))

	// Parse JSON response
	log.Debug("Parsing JSON response")
	var searchResp SearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		log.Error("Failed to parse JSON response: %v", err)
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}
	
	log.Info("Search completed successfully, found %d results", len(searchResp.Result))
	return &searchResp, nil
}

// GetCardByID gets a card by its ID
func (c *Client) GetCardByID(cardID int) (*GetCardResponse, error) {
	log.Info("Getting card by ID: %d", cardID)
	
	// Construct the URL
	apiURL := fmt.Sprintf("%s/api/v0/card/%d", c.baseURL, cardID)
	log.Debug("API URL: %s", apiURL)

	// Make the request
	log.Debug("Making HTTP request to API")
	resp, err := c.httpClient.Get(apiURL)
	if err != nil {
		log.Error("Failed to make HTTP request: %v", err)
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	
	log.Debug("Received HTTP response with status code: %d", resp.StatusCode)

	// Check status code
	if resp.StatusCode != http.StatusOK {
		log.Error("API returned non-OK status code: %d", resp.StatusCode)
		return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	// Read response body
	log.Debug("Reading response body")
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Failed to read response body: %v", err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	log.Debug("Response body read, size: %d bytes", len(body))

	// Parse JSON response
	log.Debug("Parsing JSON response")
	var cardResp GetCardResponse
	if err := json.Unmarshal(body, &cardResp); err != nil {
		log.Error("Failed to parse JSON response: %v", err)
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}
	
	log.Info("Card retrieval completed successfully, card ID: %d", cardResp.ID)
	return &cardResp, nil
}