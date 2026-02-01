package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/wutachi/raspberryTemperatureSensor/internal/db"
)

type Handler struct {
	db *db.Database
}

func NewHandler(database *db.Database) *Handler {
	return &Handler{
		db: database,
	}
}

type CurrentResponse struct {
	Temp      float64   `json:"temp"`
	Humidity  float64   `json:"humidity"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}

type HistoryRequest struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/current", h.getCurrent)
	mux.HandleFunc("GET /api/history", h.getHistory)
	mux.HandleFunc("GET /api/stats/hot-days", h.getHotDays)
	mux.HandleFunc("POST /api/feeling", h.postFeeling)
}

func (h *Handler) getCurrent(w http.ResponseWriter, r *http.Request) {
	metric, err := h.db.GetLatestMetric()
	if err != nil {
		log.Printf("Error getting latest metric: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to get current reading")
		return
	}

	if metric == nil {
		respondJSON(w, http.StatusOK, CurrentResponse{
			Status: "no_data",
		})
		return
	}

	respondJSON(w, http.StatusOK, CurrentResponse{
		Temp:      metric.Temp,
		Humidity:  metric.Humidity,
		Timestamp: metric.Timestamp,
		Status:    "ok",
	})
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

func (h *Handler) getHotDays(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		var err error
		_, err = parseLimit(limitStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid limit parameter")
			return
		}
	}

	summaries, err := h.db.GetHotDays(limit)
	if err != nil {
		log.Printf("Error getting hot days: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to get hot days")
		return
	}

	respondJSON(w, http.StatusOK, summaries)
}

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

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

func parseLimit(s string) (int, error) {
	var limit int
	_, err := fmt.Sscanf(s, "%d", &limit)
	return limit, err
}
