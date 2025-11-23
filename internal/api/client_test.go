package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := NewClient()
	if c == nil {
		t.Fatal("NewClient returned nil")
	}
	if c.httpClient == nil {
		t.Error("HTTP client not initialized")
	}
}

func TestWithAuthToken(t *testing.T) {
	token := "test-token"
	c := NewClient(WithAuthToken(token))
	if c.authToken != token {
		t.Errorf("expected token %q, got %q", token, c.authToken)
	}
}

func TestDoRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	c := NewClient(WithBaseURL(server.URL), WithAuthToken("test-token"))
	resp, err := c.doRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("doRequest failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestListRecordingsSinglePage(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		// Verify auth header
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected Bearer token, got %q", auth)
		}

		// Verify URL path
		if r.URL.Path != "/clubs/test-club/recordings/" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{"identifier": "id1", "slug": "slug1", "title": "Recording 1", "duration": 100, "created": "2025-01-01T00:00:00Z"},
			{"identifier": "id2", "slug": "slug2", "title": "Recording 2", "duration": 200, "created": "2025-01-02T00:00:00Z"}
		]`))
	}))
	defer server.Close()

	c := NewClient(WithBaseURL(server.URL), WithAuthToken("test-token"))
	result, err := c.ListRecordings("test-club", nil)
	if err != nil {
		t.Fatalf("ListRecordings failed: %v", err)
	}

	if len(result.Recordings) != 2 {
		t.Errorf("expected 2 recordings, got %d", len(result.Recordings))
	}

	if result.Recordings[0].Title != "Recording 1" {
		t.Errorf("expected title 'Recording 1', got %q", result.Recordings[0].Title)
	}

	if requestCount != 1 {
		t.Errorf("expected 1 request, got %d", requestCount)
	}
}

func TestListRecordingsWithPagination(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		page := r.URL.Query().Get("page")

		if page == "" || page == "1" {
			// First page - include Link header with next
			w.Header().Set("Link", `<http://example.com?page=2>; rel="next"`)
			w.Header().Set("x-veo-total-count", "3")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{"identifier": "id1", "slug": "slug1", "title": "Recording 1", "duration": 100, "created": "2025-01-01T00:00:00Z"},
				{"identifier": "id2", "slug": "slug2", "title": "Recording 2", "duration": 200, "created": "2025-01-02T00:00:00Z"}
			]`))
		} else if page == "2" {
			// Second page - no Link header (last page)
			w.Header().Set("x-veo-total-count", "3")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[
				{"identifier": "id3", "slug": "slug3", "title": "Recording 3", "duration": 300, "created": "2025-01-03T00:00:00Z"}
			]`))
		}
	}))
	defer server.Close()

	c := NewClient(WithBaseURL(server.URL), WithAuthToken("test-token"))

	opts := &ListRecordingsOptions{FetchAll: true}
	result, err := c.ListRecordings("test-club", opts)
	if err != nil {
		t.Fatalf("ListRecordings failed: %v", err)
	}

	if len(result.Recordings) != 3 {
		t.Errorf("expected 3 recordings, got %d", len(result.Recordings))
	}

	if requestCount != 2 {
		t.Errorf("expected 2 requests, got %d", requestCount)
	}

	if result.Recordings[2].Title != "Recording 3" {
		t.Errorf("expected third recording title 'Recording 3', got %q", result.Recordings[2].Title)
	}

	if result.TotalCount != 3 {
		t.Errorf("expected total count 3, got %d", result.TotalCount)
	}
}

func TestListRecordingsSpecificPage(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		page := r.URL.Query().Get("page")

		if page != "2" {
			t.Errorf("expected page=2, got page=%s", page)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[
			{"identifier": "id3", "slug": "slug3", "title": "Recording 3", "duration": 300, "created": "2025-01-03T00:00:00Z"}
		]`))
	}))
	defer server.Close()

	c := NewClient(WithBaseURL(server.URL), WithAuthToken("test-token"))

	opts := &ListRecordingsOptions{Page: 2}
	result, err := c.ListRecordings("test-club", opts)
	if err != nil {
		t.Fatalf("ListRecordings failed: %v", err)
	}

	if len(result.Recordings) != 1 {
		t.Errorf("expected 1 recording, got %d", len(result.Recordings))
	}

	if requestCount != 1 {
		t.Errorf("expected 1 request for specific page, got %d", requestCount)
	}
}

func TestHasNextPage(t *testing.T) {
	tests := []struct {
		name       string
		linkHeader string
		expected   bool
	}{
		{
			name:       "has next page",
			linkHeader: `<http://example.com?page=2>; rel="next"`,
			expected:   true,
		},
		{
			name:       "no next page",
			linkHeader: `<http://example.com?page=1>; rel="previous"`,
			expected:   false,
		},
		{
			name:       "empty header",
			linkHeader: "",
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasNextPage(tt.linkHeader)
			if result != tt.expected {
				t.Errorf("hasNextPage(%q) = %v, expected %v", tt.linkHeader, result, tt.expected)
			}
		})
	}
}
