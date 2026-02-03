package api

import (
	"net/http"

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

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/current", h.getCurrent)
	mux.HandleFunc("GET /api/current/html", h.getCurrentHTML)
	mux.HandleFunc("GET /api/current/stream", h.streamCurrent)
	mux.HandleFunc("GET /api/history", h.getHistory)
	mux.HandleFunc("GET /api/stats/hot-days", h.getHotDays)
	mux.HandleFunc("GET /api/stats/hot-days/html", h.getHotDaysHTML)
	mux.HandleFunc("POST /api/feeling", h.postFeeling)
}
