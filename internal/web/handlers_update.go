package web

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/sthbryan/ftm/internal/version"
)

func (h *Handlers) handleUpdate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getUpdate(w)
	case http.MethodPost:
		h.postUpdate(w)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handlers) getUpdate(w http.ResponseWriter) {
	info := h.server.updateSvc.Info()
	w.Header().Set("Content-Type", "application/json")
	if info == nil {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"current":   version.Version,
			"hasUpdate": false,
		})
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"current":    version.Version,
		"latest":     info.LatestVersion,
		"tag":        info.Tag,
		"assetName":  info.AssetName,
		"releaseUrl": info.ReleaseURL,
		"hasUpdate":  info.HasUpdate,
	})
}

func (h *Handlers) postUpdate(w http.ResponseWriter) {
	info := h.server.updateSvc.Info()
	if info == nil || !info.HasUpdate {
		http.Error(w, "no update available", http.StatusBadRequest)
		return
	}
	if err := h.server.updateSvc.Apply(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"ok": true})
	go func() {
		time.Sleep(500 * time.Millisecond)
		os.Exit(0)
	}()
}
