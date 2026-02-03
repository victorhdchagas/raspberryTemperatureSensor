package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Event struct {
	Event string
	Data  string
}

func (h *Handler) streamCurrent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	sendMetric := func() {
		metric, err := h.db.GetLatestMetric()
		if err != nil {
			log.Printf("Error getting metric for SSE: %v", err)
			return
		}

		if metric == nil {
			data := map[string]interface{}{
				"status": "no_data",
			}
			jsonData, _ := json.Marshal(data)
			fmt.Fprintf(w, "data: %s\n\n", jsonData)
			flusher.Flush()
			return
		}

		data := map[string]interface{}{
			"temp":      metric.Temp,
			"humidity":  metric.Humidity,
			"timestamp": metric.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
			"status":    "ok",
		}
		jsonData, _ := json.Marshal(data)
		fmt.Fprintf(w, "data: %s\n\n", jsonData)
		flusher.Flush()
	}

	sendMetric()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			sendMetric()
		}
	}
}
