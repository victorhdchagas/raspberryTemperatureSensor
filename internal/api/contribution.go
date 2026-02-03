package api

import (
	"log"
	"net/http"
	"time"

	"github.com/wutachi/raspberryTemperatureSensor/internal/db"
)

type ContributionResponse struct {
	Date        string  `json:"date"`
	AvgTemp     float64 `json:"avg_temp"`
	AvgHumidity float64 `json:"avg_humidity"`
	MaxTemp     float64 `json:"max_temp"`
	MinTemp     float64 `json:"min_temp"`
}

func (h *Handler) getContributionGraph(w http.ResponseWriter, r *http.Request) {
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	if startStr == "" || endStr == "" {
		respondError(w, http.StatusBadRequest, "start and end parameters are required")
		return
	}

	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid start date format (use YYYY-MM-DD)")
		return
	}

	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid end date format (use YYYY-MM-DD)")
		return
	}

	summaries, err := h.db.GetDailySummariesByRange(start, end)
	if err != nil {
		log.Printf("Error getting daily summaries: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to get contribution data")
		return
	}

	response := make([]ContributionResponse, len(summaries))
	for i := range summaries {
		response[i] = dailySummaryToResponse(&summaries[i])
	}

	respondJSON(w, http.StatusOK, response)
}

func (h *Handler) getDayDetails(w http.ResponseWriter, r *http.Request) {
	dateStr := r.PathValue("date")
	if dateStr == "" {
		respondError(w, http.StatusBadRequest, "date parameter is required")
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid date format (use YYYY-MM-DD)")
		return
	}

	summary, err := h.db.GetDailySummaryByDate(date)
	if err != nil {
		log.Printf("Error getting daily summary: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to get day details")
		return
	}

	if summary == nil {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"date":    dateStr,
			"message": "No data available for this date",
		})
		return
	}

	response := dailySummaryToResponse(summary)
	respondJSON(w, http.StatusOK, response)
}

func dailySummaryToResponse(summary *db.DailySummary) ContributionResponse {
	return ContributionResponse{
		Date:        summary.Date.Format("2006-01-02"),
		AvgTemp:     summary.AvgTemp,
		AvgHumidity: summary.AvgHumidity,
		MaxTemp:     summary.MaxTemp,
		MinTemp:     summary.MinTemp,
	}
}
