package api

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/wutachi/raspberryTemperatureSensor/internal/db"
)

type FeelingRequest struct {
	Date       string `json:"date"`
	Rating     int    `json:"rating"`
	Note       string `json:"note"`
	FeelingTag string `json:"feeling_tag"`
}

type FeelingResponse struct {
	Message string `json:"message"`
}

func (h *Handler) postFeeling(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	dateStr := r.FormValue("date")
	ratingStr := r.FormValue("rating")
	note := r.FormValue("note")
	feelingTag := r.FormValue("feeling_tag")

	if dateStr == "" || ratingStr == "" {
		respondError(w, http.StatusBadRequest, "date and rating are required")
		return
	}

	rating, err := strconv.Atoi(ratingStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid rating value")
		return
	}

	if rating < 1 || rating > 5 {
		respondError(w, http.StatusBadRequest, "Rating must be between 1 and 5")
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid date format (use YYYY-MM-DD)")
		return
	}

	if _, err := h.db.CreateSummaryForDateIfNotExists(date); err != nil {
		log.Printf("Error creating daily summary: %v", err)
	}

	userLog := db.UserLog{
		Date:       date,
		Rating:     rating,
		Note:       note,
		FeelingTag: feelingTag,
	}

	if err := h.db.InsertUserLog(userLog); err != nil {
		log.Printf("Error inserting user log: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to save feeling")
		return
	}

	respondJSON(w, http.StatusOK, FeelingResponse{
		Message: "Feeling saved successfully",
	})
}

func (h *Handler) getFeelingByDate(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		respondError(w, http.StatusBadRequest, "date parameter is required")
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid date format (use YYYY-MM-DD)")
		return
	}

	feeling, err := h.db.GetUserLogByDate(date)
	if err != nil {
		log.Printf("Error getting user log: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to get feeling")
		return
	}

	if feeling == nil {
		respondJSON(w, http.StatusOK, nil)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"rating": feeling.Rating,
		"note":   feeling.Note,
		"tag":    feeling.FeelingTag,
		"date":   feeling.Date.Format("2006-01-02"),
	})
}
