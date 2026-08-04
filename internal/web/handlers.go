package web

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/sthbryan/ftm/internal/config"
	"github.com/sthbryan/ftm/internal/process"
)

type Handlers struct {
	manager *process.Manager
	config  *config.Config
	server  *Server
}

func NewHandlers(manager *process.Manager, cfg *config.Config, server *Server) *Handlers {
	return &Handlers{
		manager: manager,
		config:  cfg,
		server:  server,
	}
}

func (h *Handlers) Route(w http.ResponseWriter, r *http.Request) {
	if !h.setCORS(w, r) {
		http.Error(w, "forbidden origin", http.StatusForbidden)
		return
	}

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	path := strings.TrimSuffix(r.URL.Path, "/")
	if path == "" {
		path = "/"
	}

	switch {
	case path == "/api/tunnels":
		h.handleTunnels(w, r)
	case strings.HasPrefix(path, "/api/logs/"):
		h.handleLogs(w, r)
	case path == "/api/status":
		h.handleStatus(w)
	case path == "/api/settings":
		h.handleSettings(w, r)
	case path == "/api/providers":
		h.handleProviders(w)
	case path == "/api/i18n":
		h.handleI18n(w, r)
	case path == "/api/i18n/current":
		h.handleI18nCurrent(w, r)
	case path == "/api/update":
		h.handleUpdate(w, r)
	case path == "/api/detect-port":
		h.handleDetectPort(w)
	case strings.HasPrefix(path, "/api/tunnels/"):
		h.handleTunnelActions(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *Handlers) setCORS(w http.ResponseWriter, r *http.Request) bool {
	origin, ok := allowedOrigin(r)
	if !ok {
		return false
	}

	if origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Add("Vary", "Origin")
	}
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	return true
}

func (h *Handlers) tunnelToMap(t config.TunnelConfig) map[string]interface{} {
	item := map[string]interface{}{
		"id":             t.ID,
		"name":           t.Name,
		"provider":       string(t.Provider),
		"port":           t.LocalPort,
		"state":          "stopped",
		"publicUrl":      "",
		"errorMessage":   "",
		"expiresAt":      int64(0),
		"startedAt":      int64(0),
		"visitors":       0,
		"activeSessions": 0,
	}

	if status, ok := h.manager.GetStatus(t.ID); ok {
		item["publicUrl"] = status.PublicURL
		item["errorMessage"] = status.ErrorMessage
		item["state"] = string(status.State)
		item["expiresAt"] = status.ExpiresAt
		item["startedAt"] = status.StartedAt
		item["visitors"] = status.Visitors
		item["activeSessions"] = status.ActiveSessions
	}

	if item["state"] == "stopped" {
		if needsInstall, canInstall := h.manager.CheckInstallation(t.Provider); needsInstall && canInstall {
			item["state"] = "need_installing"
		}
	}

	return item
}

func (h *Handlers) writeTunnelJSON(w http.ResponseWriter, t config.TunnelConfig) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.tunnelToMap(t))
}

func MarshalJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
