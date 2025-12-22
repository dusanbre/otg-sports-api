package goalserve

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// Client represents the Goalserve API client
type Client struct {
	BaseURL     string
	APIKey      string
	HTTPClient  *http.Client
	rateLimiter *time.Ticker
}

// NewClient creates a new Goalserve API client
func NewClient() *Client {
	return &Client{
		BaseURL: getEnv("GOALSERVE_URL", "https://www.goalserve.com"),
		APIKey:  getEnv("GOALSERVE_API_KEY", ""),
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		rateLimiter: time.NewTicker(1 * time.Second), // 1 request per second
	}
}

// Close closes the client and releases resources
func (c *Client) Close() {
	if c.rateLimiter != nil {
		c.rateLimiter.Stop()
	}
}

// FetchTodayMatches fetches today's soccer matches from Goalserve API
func (c *Client) FetchTodayMatches() (*GoalServeScores, error) {
	// Wait for rate limiter
	<-c.rateLimiter.C

	url := fmt.Sprintf("%s/getfeed/%s/soccernew/home?json=1", c.BaseURL, c.APIKey)

	log.Printf("Fetching matches from GoalServe: %s", url)

	return c.fetchMatchesFromURL(url)
}

// FetchMatchesPast7Days fetches soccer matches for the past 7 days
func (c *Client) FetchMatchesPast7Days() (*GoalServeScores, error) {
	var allScores GoalServeScores
	allScores.Categories = make([]GoalServeScoreCategory, 0)

	// Fetch matches for past 7 days (d-1 to d-7)
	for day := 1; day <= 7; day++ {
		// Wait for rate limiter
		<-c.rateLimiter.C

		url := fmt.Sprintf("%s/getfeed/%s/soccernew/d-%d?json=1", c.BaseURL, c.APIKey, day)

		log.Printf("Fetching past matches from GoalServe (day %d): %s", day, url)

		scores, err := c.fetchMatchesFromURL(url)
		if err != nil {
			log.Printf("Failed to fetch matches for past day %d: %v", day, err)
			continue // Continue with other days even if one fails
		}

		// Merge categories
		allScores.Categories = append(allScores.Categories, scores.Categories...)
	}

	// Count total matches for logging
	var totalMatches int
	for _, category := range allScores.Categories {
		totalMatches += len(category.Matches.Match)
	}

	log.Printf("Successfully fetched past 7 days matches: %d total", totalMatches)
	return &allScores, nil
}

// FetchMatchesFuture7Days fetches soccer matches for the next 7 days
func (c *Client) FetchMatchesFuture7Days() (*GoalServeScores, error) {
	var allScores GoalServeScores
	allScores.Categories = make([]GoalServeScoreCategory, 0)

	// Fetch matches for next 7 days (d1 to d7)
	for day := 1; day <= 7; day++ {
		// Wait for rate limiter
		<-c.rateLimiter.C

		url := fmt.Sprintf("%s/getfeed/%s/soccernew/d%d?json=1", c.BaseURL, c.APIKey, day)

		log.Printf("Fetching future matches from GoalServe (day %d): %s", day, url)

		scores, err := c.fetchMatchesFromURL(url)
		if err != nil {
			log.Printf("Failed to fetch matches for future day %d: %v", day, err)
			continue // Continue with other days even if one fails
		}

		// Merge categories
		allScores.Categories = append(allScores.Categories, scores.Categories...)
	}

	// Count total matches for logging
	var totalMatches int
	for _, category := range allScores.Categories {
		totalMatches += len(category.Matches.Match)
	}

	log.Printf("Successfully fetched future 7 days matches: %d total", totalMatches)
	return &allScores, nil
}

// fetchMatchesFromURL is a helper function to fetch matches from a specific URL
func (c *Client) fetchMatchesFromURL(url string) (*GoalServeScores, error) {
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch matches: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse JSON response - the API response structure has scores at root level
	var jsonResponse map[string]interface{}
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Extract the scores object from the response
	scoresData, ok := jsonResponse["scores"]
	if !ok {
		return nil, fmt.Errorf("no scores field found in response")
	}

	// Convert scores data back to JSON for proper unmarshaling
	scoresJSON, err := json.Marshal(scoresData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal scores data: %w", err)
	}

	// Unmarshal directly into GoalServeScores
	var scores GoalServeScores
	if err := json.Unmarshal(scoresJSON, &scores); err != nil {
		return nil, fmt.Errorf("failed to parse scores JSON: %w", err)
	}

	// Count total matches for logging
	var totalMatches int
	for _, category := range scores.Categories {
		totalMatches += len(category.Matches.Match)
	}

	log.Printf("Successfully fetched matches: %d total", totalMatches)
	return &scores, nil
}

// getEnv gets an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
