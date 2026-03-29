package web

import (
	"net/http"
	"strings"
)

func (h *Handlers) handleLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := h.extractLogID(r.URL.Path)
	if id == "" {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	logs := h.manager.GetLogs(id)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(strings.Join(logs, "\n")))
}

func (h *Handlers) extractLogID(path string) string {
	path = strings.TrimPrefix(path, "/api/logs/")
	parts := strings.Split(path, "/")
	if len(parts) > 0 && parts[0] != "" {
		return parts[0]
	}
	return ""
}
