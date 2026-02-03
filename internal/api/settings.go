package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type SettingsRequest struct {
	SensorInterval int `json:"sensorInterval"`
}

type SettingsResponse struct {
	SensorIntervalMinutes int `json:"sensorIntervalMinutes"`
}

func (h *Handler) getSettings(w http.ResponseWriter, r *http.Request) {
	interval := h.config.GetSensorReadingInterval()
	minutes := int(interval.Minutes())

	respondJSON(w, http.StatusOK, SettingsResponse{
		SensorIntervalMinutes: minutes,
	})
}

func (h *Handler) updateSettings(w http.ResponseWriter, r *http.Request) {
	var req SettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.SensorInterval < 1 || req.SensorInterval > 1440 {
		respondError(w, http.StatusBadRequest, "Sensor interval must be between 1 and 1440 minutes")
		return
	}

	h.config.SetSensorReadingInterval(time.Duration(req.SensorInterval) * time.Minute)
	log.Printf("Settings updated: sensor interval = %d minutes", req.SensorInterval)

	respondJSON(w, http.StatusOK, SettingsResponse{
		SensorIntervalMinutes: req.SensorInterval,
	})
}
