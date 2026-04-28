package accounts

import (
	"net/http"

	"ds2api/internal/config"
)

func (h *Handler) queueStatus(w http.ResponseWriter, _ *http.Request) {
	status := h.Pool.Status()
	config.Logger.Debug("[queue_status] returning status", "today_requests", status["today_requests"])
	writeJSON(w, http.StatusOK, status)
}
