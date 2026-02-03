package api

import (
	"encoding/json"
	"log"
	"net/http"
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
	var req FeelingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Rating < 1 || req.Rating > 5 {
		respondError(w, http.StatusBadRequest, "Rating must be between 1 and 5")
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid date format (use YYYY-MM-DD)")
		return
	}

	userLog := db.UserLog{
		Date:       date,
		Rating:     req.Rating,
		Note:       req.Note,
		FeelingTag: req.FeelingTag,
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
