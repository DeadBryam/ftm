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
	resp := map[string]interface{}{
		"current":    version.Version,
		"latest":     "",
		"tag":        "",
		"assetName":  "",
		"releaseUrl": "",
		"hasUpdate":  false,
	}
	if info := h.server.updateSvc.Info(); info != nil {
		resp["latest"] = info.LatestVersion
		resp["tag"] = info.Tag
		resp["assetName"] = info.AssetName
		resp["releaseUrl"] = info.ReleaseURL
		resp["hasUpdate"] = info.HasUpdate
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
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
