package maintenance

import (
	"context"
	"log"
	"time"

	"github.com/wutachi/raspberryTemperatureSensor/internal/db"
)

type Worker struct {
	db *db.Database
}

func NewWorker(database *db.Database) *Worker {
	return &Worker{
		db: database,
	}
}

func (w *Worker) Start(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Printf("Maintenance worker started: running every %v", interval)

	for {
		select {
		case <-ctx.Done():
			log.Println("Maintenance worker stopped")
			return
		case <-ticker.C:
			if err := w.runMaintenance(); err != nil {
				log.Printf("Maintenance error: %v", err)
			}
		}
	}
}

func (w *Worker) runMaintenance() error {
	log.Println("Running daily maintenance...")

	yesterday := time.Now().AddDate(0, 0, -1)
	yesterday = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, time.UTC)

	summary, err := w.db.CalculateDailySummary(yesterday)
	if err != nil {
		log.Printf("Error calculating daily summary for %v: %v", yesterday, err)
		return err
	}

	if summary != nil {
		if err := w.db.InsertDailySummary(*summary); err != nil {
			log.Printf("Error inserting daily summary: %v", err)
			return err
		}
		log.Printf("Daily summary created for %v: AvgTemp=%.2f°C, AvgHumidity=%.2f%%",
			yesterday.Format("2006-01-02"), summary.AvgTemp, summary.AvgHumidity)
	}

	if err := w.db.DeleteOldMetrics(30); err != nil {
		log.Printf("Error deleting old metrics: %v", err)
		return err
	}

	log.Println("Maintenance completed")
	return nil
}

func (w *Worker) CreateSummaryForDate(date time.Time) error {
	summary, err := w.db.CalculateDailySummary(date)
	if err != nil {
		log.Printf("Error calculating daily summary for %v: %v", date, err)
		return err
	}

	if summary != nil {
		if err := w.db.InsertDailySummary(*summary); err != nil {
			log.Printf("Error inserting daily summary: %v", err)
			return err
		}
		log.Printf("Daily summary created for %v: AvgTemp=%.2f°C, AvgHumidity=%.2f%%",
			date.Format("2006-01-02"), summary.AvgTemp, summary.AvgHumidity)
	}

	return nil
}
