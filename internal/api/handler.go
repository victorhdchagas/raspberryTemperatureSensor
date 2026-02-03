package api

import (
	"net/http"

	"github.com/wutachi/raspberryTemperatureSensor/internal/app"
	"github.com/wutachi/raspberryTemperatureSensor/internal/db"
)

type Handler struct {
	db     *db.Database
	config *app.Config
}

func NewHandler(database *db.Database, config *app.Config) *Handler {
	return &Handler{
		db:     database,
		config: config,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/current", h.getCurrent)
	mux.HandleFunc("GET /api/current/html", h.getCurrentHTML)
	mux.HandleFunc("GET /api/current/stream", h.streamCurrent)
	mux.HandleFunc("GET /api/history", h.getHistory)
	mux.HandleFunc("GET /api/stats/hot-days", h.getHotDays)
	mux.HandleFunc("GET /api/stats/hot-days/html", h.getHotDaysHTML)
	mux.HandleFunc("GET /api/contribution", h.getContributionGraph)
	mux.HandleFunc("GET /api/day/{date}", h.getDayDetails)
	mux.HandleFunc("POST /api/feeling", h.postFeeling)
	mux.HandleFunc("GET /api/feeling", h.getFeelingByDate)
	mux.HandleFunc("POST /api/admin/generate-summaries", h.generateSummaries)
	mux.HandleFunc("GET /api/settings", h.getSettings)
	mux.HandleFunc("POST /api/settings", h.updateSettings)
}
