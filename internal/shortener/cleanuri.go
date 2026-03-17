package shortener

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const cleanuriAPI = "https://cleanuri.com/api/v1/shorten"

type CleanURIClient struct {
	client *http.Client
}

type cleanuriRequest struct {
	URL string `json:"url"`
}

type cleanuriResponse struct {
	ResultURL string `json:"result_url"`
	Error     string `json:"error"`
}

func NewCleanURI() *CleanURIClient {
	return &CleanURIClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *CleanURIClient) Name() string {
	return "cleanuri"
}

func (c *CleanURIClient) Shorten(longURL, custom string) (string, error) {
	reqBody, _ := json.Marshal(cleanuriRequest{URL: longURL})
	
	resp, err := c.client.Post(cleanuriAPI, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to contact cleanuri: %w", err)
	}
	defer resp.Body.Close()
	
	var result cleanuriResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}
	
	if result.Error != "" {
		if containsAny(result.Error, "blocked", "blacklist", "invalid") {
			return "", ShortenError{
				Reason:  "DOMAIN_BLOCKED",
				Message: result.Error,
			}
		}
		return "", ShortenError{
			Reason:  "API_ERROR",
			Message: result.Error,
		}
	}
	
	if result.ResultURL == "" {
		return "", fmt.Errorf("empty result from cleanuri")
	}
	
	return result.ResultURL, nil
}

func containsAny(s string, substrs ...string) bool {
	for _, sub := range substrs {
		if contains(s, sub) {
			return true
		}
	}
	return false
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || containsAt(s, substr, 0) || containsRecursive(s, substr))
}

func containsAt(s, substr string, start int) bool {
	if start+len(substr) > len(s) {
		return false
	}
	for i := 0; i < len(substr); i++ {
		if s[start+i] != substr[i] {
			return false
		}
	}
	return true
}

func containsRecursive(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if containsAt(s, substr, i) {
			return true
		}
	}
	return false
}
