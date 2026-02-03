package api

import (
	"fmt"
	"log"
	"net/http"
)

type CurrentResponse struct {
	Temp      float64 `json:"temp"`
	Humidity  float64 `json:"humidity"`
	Timestamp string  `json:"timestamp"`
	Status    string  `json:"status"`
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
		Timestamp: metric.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		Status:    "ok",
	})
}

func (h *Handler) getCurrentHTML(w http.ResponseWriter, r *http.Request) {
	metric, err := h.db.GetLatestMetric()
	if err != nil {
		log.Printf("Error getting latest metric: %v", err)
		respondHTML(w, http.StatusInternalServerError, "<p class=\"text-red-400\">Erro ao carregar dados</p>")
		return
	}

	if metric == nil {
		respondHTML(w, http.StatusOK, `<p class="text-gray-400">Sem dados disponíveis</p>`)
		return
	}

	html := `
		<p class="text-6xl font-bold text-green-400 mb-1">` + formatFloat(metric.Temp) + `°C</p>
		<p class="text-3xl text-blue-400 mb-3">` + formatFloat(metric.Humidity) + `%</p>
		<p class="text-sm text-gray-400">` + metric.Timestamp.Format("02/01/2006 15:04:05") + `</p>
	`
	respondHTML(w, http.StatusOK, html)
}

func formatFloat(f float64) string {
	return fmt.Sprintf("%.1f", f)
}
