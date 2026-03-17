package shortener

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const tinyurlAPI = "https://tinyurl.com/api-create.php"

type TinyURLClient struct {
	client *http.Client
}

func NewTinyURL() *TinyURLClient {
	return &TinyURLClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *TinyURLClient) Name() string {
	return "tinyurl"
}

func (c *TinyURLClient) Shorten(longURL, custom string) (string, error) {
	apiURL := fmt.Sprintf("%s?url=%s", tinyurlAPI, url.QueryEscape(longURL))
	
	resp, err := c.client.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("failed to contact tinyurl: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}
	
	shortURL := strings.TrimSpace(string(body))
	
	if !strings.HasPrefix(shortURL, "http") {
		return "", ShortenError{
			Reason:  "API_ERROR",
			Message: shortURL,
		}
	}
	
	return shortURL, nil
}
