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
