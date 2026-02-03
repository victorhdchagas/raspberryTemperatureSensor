package api

import (
	"log"
	"net/http"
	"time"
)

func (h *Handler) generateSummaries(w http.ResponseWriter, r *http.Request) {
	today := time.Now()
	oneYearAgo := today.AddDate(-1, 0, 0)

	start := time.Date(oneYearAgo.Year(), oneYearAgo.Month(), oneYearAgo.Day(), 0, 0, 0, 0, time.UTC)
	end := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)

	log.Printf("Generating summaries from %s to %s", start.Format("2006-01-02"), end.Format("2006-01-02"))

	if err := h.db.CreateSummariesForRange(start, end); err != nil {
		log.Printf("Error generating summaries: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to generate summaries")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "Summaries generated successfully",
		"start":   start.Format("2006-01-02"),
		"end":     end.Format("2006-01-02"),
	})
}
