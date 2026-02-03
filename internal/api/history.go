package api

import (
	"log"
	"net/http"
	"time"
)

type HistoryRequest struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

func (h *Handler) getHistory(w http.ResponseWriter, r *http.Request) {
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	if startStr == "" || endStr == "" {
		respondError(w, http.StatusBadRequest, "start and end parameters are required")
		return
	}

	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid start date format")
		return
	}

	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid end date format")
		return
	}

	metrics, err := h.db.GetMetricsByDateRange(start, end)
	if err != nil {
		log.Printf("Error getting metrics: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to get history")
		return
	}

	respondJSON(w, http.StatusOK, metrics)
}
