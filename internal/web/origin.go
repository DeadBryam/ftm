package web

import (
	"net"
	"net/http"
	"net/url"
	"strings"
)

func isLoopbackHost(host string) bool {
	if host == "localhost" || strings.HasSuffix(host, ".localhost") {
		return true
	}
	if ip := net.ParseIP(host); ip != nil {
		return ip.IsLoopback()
	}
	return false
}

func isWailsHost(host string) bool {
	return host == "wails" || host == "wails.localhost" || strings.HasSuffix(host, ".wails.localhost")
}

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

func guardHost(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !allowedHost(r) {
			http.Error(w, "forbidden host", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
