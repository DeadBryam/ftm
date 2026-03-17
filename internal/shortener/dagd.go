package shortener

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const dagdAPI = "https://da.gd/shorten"

type DagdClient struct {
	client *http.Client
}

func NewDagd() *DagdClient {
	return &DagdClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *DagdClient) Name() string {
	return "dagd"
}

func (c *DagdClient) Shorten(longURL, custom string) (string, error) {
	apiURL := fmt.Sprintf("%s?url=%s", dagdAPI, url.QueryEscape(longURL))
	
	resp, err := c.client.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("failed to contact da.gd: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}
	
	shortURL := strings.TrimSpace(string(body))
	
	if strings.Contains(shortURL, "Blacklisted") || strings.Contains(shortURL, "blacklist") {
		return "", ShortenError{
			Reason:  "DOMAIN_BLOCKED",
			Message: "This domain is blocked by da.gd",
		}
	}
	
	if !strings.HasPrefix(shortURL, "http") {
		return "", ShortenError{
			Reason:  "API_ERROR",
			Message: shortURL,
		}
	}
	
	return shortURL, nil
}
