package accounts

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) queueStatus(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.Pool.Status())
}

func (h *Handler) setAccountPaused(w http.ResponseWriter, r *http.Request) {
	identifier := chi.URLParam(r, "identifier")
	if decoded, err := url.PathUnescape(identifier); err == nil {
		identifier = decoded
	}
	var req map[string]any
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"detail": "invalid json"})
		return
	}
	paused, ok := req["paused"].(bool)
	if !ok {
		writeJSON(w, http.StatusBadRequest, map[string]any{"detail": "paused must be boolean"})
		return
	}
	if _, ok := h.Store.FindAccount(identifier); !ok {
		writeJSON(w, http.StatusNotFound, map[string]any{"detail": "账号不存在"})
		return
	}
	if err := h.Store.UpdateAccountPaused(identifier, paused); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"detail": err.Error()})
		return
	}
	h.Pool.Reset()
	writeJSON(w, http.StatusOK, map[string]any{"success": true, "paused": paused})
}
