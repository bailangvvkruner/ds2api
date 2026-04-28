package accounts

import "net/http"

func (h *Handler) statistics(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.Pool.GetStatistics())
}
