package web

import (
	"net"
	"net/http"
	"net/url"
	"strings"
)

// The dashboard API is unauthenticated: anything that can reach it can start,
// stop and delete tunnels. Binding to loopback keeps other machines out, but a
// browser on this machine reaches loopback just fine, so any page the user
// visits could drive the API from JavaScript. Checking Origin is what stops
// that, since browsers set it on cross-origin requests and scripts cannot
// forge it.

// isLoopbackHost reports whether host refers to this machine. Covers the
// literal IPs as well as "localhost" and the reserved .localhost subdomains.
func isLoopbackHost(host string) bool {
	if host == "localhost" || strings.HasSuffix(host, ".localhost") {
		return true
	}
	if ip := net.ParseIP(host); ip != nil {
		return ip.IsLoopback()
	}
	return false
}

// isWailsHost reports whether host belongs to the Wails desktop webview. The
// webview opens its WebSocket straight to 127.0.0.1, so those requests really
// are cross-origin and have to be allowed by name.
func isWailsHost(host string) bool {
	return host == "wails" || host == "wails.localhost" || strings.HasSuffix(host, ".wails.localhost")
}

// allowedOrigin reports whether the request may drive the local API, returning
// the origin to echo back in CORS headers.
//
// A missing Origin is allowed: browsers always send one for cross-origin
// requests, so a request without it is same-origin or comes from a non-browser
// client (curl, the desktop shell) that was never subject to the same-origin
// policy in the first place.
func allowedOrigin(r *http.Request) (origin string, ok bool) {
	origin = r.Header.Get("Origin")
	if origin == "" {
		return "", true
	}

	u, err := url.Parse(origin)
	if err != nil {
		return "", false
	}

	host := u.Hostname()
	if isLoopbackHost(host) || isWailsHost(host) {
		return origin, true
	}

	return "", false
}

// allowedHost reports whether the request's Host header names this machine.
//
// This closes DNS rebinding, which the Origin check alone does not: an attacker
// points evil.com at 127.0.0.1, so their page reaches the dashboard as a
// *same-origin* request and never sends an Origin header at all. The Host
// header still says "evil.com", which is what gives it away.
func allowedHost(r *http.Request) bool {
	host := r.Host
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}
	host = strings.Trim(host, "[]")

	if host == "" {
		return false
	}

	return isLoopbackHost(host) || isWailsHost(host)
}

// guardHost rejects requests whose Host header does not name this machine.
func guardHost(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !allowedHost(r) {
			http.Error(w, "forbidden host", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
