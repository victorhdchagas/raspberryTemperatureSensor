package db

import (
	"database/sql"
	"log"
	"time"
)

func (d *Database) InsertDailySummary(summary DailySummary) error {
	_, err := d.db.Exec(
		"INSERT OR REPLACE INTO daily_summaries (date, avg_temp, avg_humidity, max_temp, min_temp) VALUES (?, ?, ?, ?, ?)",
		summary.Date, summary.AvgTemp, summary.AvgHumidity, summary.MaxTemp, summary.MinTemp,
	)
	return err
}

func (d *Database) GetHotDays(limit int) ([]DailySummary, error) {
	rows, err := d.db.Query(
		"SELECT id, date, avg_temp, avg_humidity, max_temp, min_temp FROM daily_summaries ORDER BY avg_temp DESC LIMIT ?",
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []DailySummary
	for rows.Next() {
		var summary DailySummary
		if err := rows.Scan(&summary.ID, &summary.Date, &summary.AvgTemp, &summary.AvgHumidity, &summary.MaxTemp, &summary.MinTemp); err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}

	return summaries, nil
}

func (d *Database) GetDailySummaryByDate(date time.Time) (*DailySummary, error) {
	row := d.db.QueryRow(
		"SELECT id, date, avg_temp, avg_humidity, max_temp, min_temp FROM daily_summaries WHERE date = ?",
		date.Format("2006-01-02"),
	)

	var summary DailySummary
	err := row.Scan(&summary.ID, &summary.Date, &summary.AvgTemp, &summary.AvgHumidity, &summary.MaxTemp, &summary.MinTemp)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &summary, nil
}

func (d *Database) GetDailySummariesByRange(start, end time.Time) ([]DailySummary, error) {
	rows, err := d.db.Query(
		"SELECT id, date, avg_temp, avg_humidity, max_temp, min_temp FROM daily_summaries WHERE date >= ? AND date <= ? ORDER BY date ASC",
		start, end,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []DailySummary
	for rows.Next() {
		var summary DailySummary
		if err := rows.Scan(&summary.ID, &summary.Date, &summary.AvgTemp, &summary.AvgHumidity, &summary.MaxTemp, &summary.MinTemp); err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}

	return summaries, nil
}

func (d *Database) CalculateDailySummary(date time.Time) (*DailySummary, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	end := start.Add(24 * time.Hour)

	row := d.db.QueryRow(
		`SELECT 
			AVG(temp) as avg_temp,
			AVG(humidity) as avg_humidity,
			MAX(temp) as max_temp,
			MIN(temp) as min_temp
		FROM raw_metrics 
		WHERE timestamp >= ? AND timestamp < ?`,
		start, end,
	)

	var summary DailySummary
	err := row.Scan(&summary.AvgTemp, &summary.AvgHumidity, &summary.MaxTemp, &summary.MinTemp)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	summary.Date = date
	return &summary, nil
}

func (d *Database) CreateSummaryForDateIfNotExists(date time.Time) error {
	summary, err := d.GetDailySummaryByDate(date)
	if err != nil {
		return err
	}

	if summary != nil {
		return nil
	}

	calculatedSummary, err := d.CalculateDailySummary(date)
	if err != nil {
		return err
	}

	if calculatedSummary != nil {
		if err := d.InsertDailySummary(*calculatedSummary); err != nil {
			return err
		}
		log.Printf("Daily summary created for %v: AvgTemp=%.2f°C, AvgHumidity=%.2f%%",
			date.Format("2006-01-02"), calculatedSummary.AvgTemp, calculatedSummary.AvgHumidity)
	}

	return nil
}

func (d *Database) CreateSummariesForRange(start, end time.Time) error {
	for date := start; !date.After(end); date = date.AddDate(0, 0, 1) {
		if err := d.CreateSummaryForDateIfNotExists(date); err != nil {
			log.Printf("Error creating summary for %v: %v", date, err)
		}
	}
	return nil
}
