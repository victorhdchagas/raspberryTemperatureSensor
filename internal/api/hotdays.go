package api

import (
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) getHotDays(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		var err error
		limit, err = parseLimit(limitStr)
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

func (h *Handler) getHotDaysHTML(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 6
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 || limit > 20 {
			limit = 6
		}
	}

	summaries, err := h.db.GetHotDays(limit)
	if err != nil {
		log.Printf("Error getting hot days: %v", err)
		respondHTML(w, http.StatusInternalServerError, "<p class=\"text-red-400\">Erro ao carregar dias quentes</p>")
		return
	}

	html := ""
	for _, day := range summaries {
		html += `
			<div class="bg-gray-700 rounded p-4">
				<p class="font-semibold text-yellow-400">` + day.Date.Format("02/01/2006") + `</p>
				<p class="text-2xl font-bold">` + formatFloat(day.AvgTemp) + `°C</p>
				<p class="text-sm text-gray-400">Média</p>
				<p class="text-sm text-gray-400">Max: ` + formatFloat(day.MaxTemp) + `°C | Min: ` + formatFloat(day.MinTemp) + `°C</p>
			</div>
		`
	}

	respondHTML(w, http.StatusOK, html)
}
