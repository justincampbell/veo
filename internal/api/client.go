package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/justincampbell/veo/internal/models"
)

// Client represents a Veo API client
type Client struct {
	baseURL    string
	httpClient *http.Client
	authToken  string
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client)

// WithAuthToken sets the authentication token
func WithAuthToken(token string) ClientOption {
	return func(c *Client) {
		c.authToken = token
	}
}

// WithBaseURL sets a custom base URL
func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// NewClient creates a new Veo API client
func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		baseURL: "https://app.veo.co/api/app",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// doRequest performs an HTTP request with authentication
func (c *Client) doRequest(method, path string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(method, c.baseURL+path, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	if c.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.authToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

// decodeResponse decodes a JSON response into the target interface
func decodeResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	if target != nil {
		if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// ListRecordingsOptions contains options for listing recordings
type ListRecordingsOptions struct {
	Page     int  // Page number (1-indexed, 0 means first page)
	FetchAll bool // If true, fetch all pages
}

// ListRecordingsResult contains recordings and metadata
type ListRecordingsResult struct {
	Recordings []models.Recording
	TotalCount int // Total count from API (x-veo-total-count header)
}

// ListRecordings lists recordings for a club with pagination support
func (c *Client) ListRecordings(clubSlug string, opts *ListRecordingsOptions) (*ListRecordingsResult, error) {
	if opts == nil {
		opts = &ListRecordingsOptions{Page: 1}
	}

	// Build query parameters with fields we want
	params := url.Values{}
	params.Set("filter", "own")
	params.Set("fields", "camera")
	params.Add("fields", "created")
	params.Add("fields", "start")
	params.Add("fields", "duration")
	params.Add("fields", "identifier")
	params.Add("fields", "slug")
	params.Add("fields", "title")
	params.Add("fields", "url")
	params.Add("fields", "thumbnail")
	params.Add("fields", "reel_url")
	params.Add("fields", "team")
	params.Add("fields", "privacy")
	params.Add("fields", "permissions")
	params.Add("fields", "is_accessible")

	var allRecordings []models.Recording
	var totalCount int
	page := opts.Page
	if page == 0 {
		page = 1
	}

	for {
		// Add page parameter if not first page
		if page > 1 {
			params.Set("page", fmt.Sprintf("%d", page))
		}

		path := fmt.Sprintf("/clubs/%s/recordings/?%s", clubSlug, params.Encode())

		resp, err := c.doRequest("GET", path, nil)
		if err != nil {
			return nil, err
		}

		// Get total count from header (only need to do this once)
		if totalCount == 0 {
			if countStr := resp.Header.Get("x-veo-total-count"); countStr != "" {
				fmt.Sscanf(countStr, "%d", &totalCount)
			}
		}

		var recordings []models.Recording
		if err := decodeResponse(resp, &recordings); err != nil {
			return nil, err
		}

		allRecordings = append(allRecordings, recordings...)

		// If not fetching all, or no more pages, stop
		if !opts.FetchAll {
			break
		}

		// Check for next page in Link header
		linkHeader := resp.Header.Get("Link")
		if !hasNextPage(linkHeader) {
			break
		}

		page++
	}

	return &ListRecordingsResult{
		Recordings: allRecordings,
		TotalCount: totalCount,
	}, nil
}

// hasNextPage checks if the Link header contains a "next" relation
func hasNextPage(linkHeader string) bool {
	// Check if Link header contains rel="next"
	for i := 0; i <= len(linkHeader)-10; i++ {
		if linkHeader[i:i+10] == `rel="next"` {
			return true
		}
	}
	return false
}

// RecordingDetails represents detailed information about a recording/match
type RecordingDetails struct {
	ID                   string                 `json:"id"`
	Identifier           string                 `json:"identifier"`
	Slug                 string                 `json:"slug"`
	Title                string                 `json:"title"`
	Created              time.Time              `json:"created"`
	Start                time.Time              `json:"start"`
	End                  time.Time              `json:"end"`
	Duration             int                    `json:"duration"`
	Type                 string                 `json:"type"`
	OwnTeamHomeOrAway    string                 `json:"own_team_home_or_away"`
	OpponentTeamName     string                 `json:"opponent_team_name"`
	OpponentClubName     string                 `json:"opponent_club_name"`
	OpponentTeamColor    string                 `json:"opponent_team_color"`
	OpponentShortName    string                 `json:"opponent_short_name"`
	OwnTeamColor         string                 `json:"own_team_color"`
	OwnTeamFormation     string                 `json:"own_team_formation"`
	OpponentTeamFormation string                `json:"opponent_team_formation"`
	Team                 interface{}            `json:"team"` // Can be string, object, or null
	Info                 map[string]interface{} `json:"info"`
	Permissions          map[string]interface{} `json:"permissions"`
}

// GetRecording retrieves detailed information about a specific recording/match
func (c *Client) GetRecording(identifier string) (*RecordingDetails, error) {
	path := fmt.Sprintf("/matches/%s/", identifier)

	resp, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var details RecordingDetails
	if err := decodeResponse(resp, &details); err != nil {
		return nil, err
	}

	return &details, nil
}
