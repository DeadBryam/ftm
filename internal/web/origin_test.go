package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAllowedOrigin(t *testing.T) {
	tests := []struct {
		name   string
		origin string
		want   bool
	}{
		{"missing origin is same-origin or non-browser", "", true},
		{"loopback ipv4", "http://127.0.0.1:40500", true},
		{"loopback ipv4 other port", "http://127.0.0.1:3000", true},
		{"loopback ipv6", "http://[::1]:40500", true},
		{"localhost", "http://localhost:40500", true},
		{"vite dev server", "http://localhost:3000", true},
		{"wails scheme", "wails://wails", true},
		{"wails localhost", "http://wails.localhost", true},

		{"attacker site", "https://evil.com", false},
		{"attacker subdomain of localhost lookalike", "https://localhost.evil.com", false},
		{"attacker with loopback in the path", "https://evil.com/127.0.0.1", false},
		{"public ip", "http://203.0.113.10", false},
		{"lan ip", "http://192.168.1.50:40500", false},
		{"null origin from sandboxed iframe", "null", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/api/tunnels", nil)
			if tt.origin != "" {
				r.Header.Set("Origin", tt.origin)
			}

			got, ok := allowedOrigin(r)
			if ok != tt.want {
				t.Fatalf("allowedOrigin(%q) = %v, want %v", tt.origin, ok, tt.want)
			}
			if ok && tt.origin != "" && got != tt.origin {
				t.Fatalf("allowedOrigin(%q) echoed %q, want the request origin", tt.origin, got)
			}
			if !ok && got != "" {
				t.Fatalf("allowedOrigin(%q) echoed %q on rejection, want empty", tt.origin, got)
			}
		})
	}
}

func TestIsLoopbackHost(t *testing.T) {
	allowed := []string{"localhost", "app.localhost", "127.0.0.1", "127.0.0.53", "::1"}
	for _, host := range allowed {
		if !isLoopbackHost(host) {
			t.Errorf("isLoopbackHost(%q) = false, want true", host)
		}
	}

	denied := []string{"evil.com", "localhost.evil.com", "192.168.1.1", "10.0.0.1", "0.0.0.0", ""}
	for _, host := range denied {
		if isLoopbackHost(host) {
			t.Errorf("isLoopbackHost(%q) = true, want false", host)
		}
	}
}
