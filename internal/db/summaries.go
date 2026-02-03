package db

import (
	"database/sql"
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
		"SELECT id, date, avg_temp, avg_humidity, max_temp, min_temp FROM daily_summaries WHERE DATE(date) = ?",
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
		"SELECT id, date, avg_temp, avg_humidity, max_temp, min_temp FROM daily_summaries WHERE DATE(date) >= ? AND DATE(date) <= ? ORDER BY date ASC",
		start.Format("2006-01-02"), end.Format("2006-01-02"),
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
			COALESCE(AVG(temp), 0) as avg_temp,
			COALESCE(AVG(humidity), 0) as avg_humidity,
			COALESCE(MAX(temp), 0) as max_temp,
			COALESCE(MIN(temp), 0) as min_temp
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

func (d *Database) CreateSummaryForDateIfNotExists(date time.Time) (bool, error) {
	summary, err := d.GetDailySummaryByDate(date)
	if err != nil {
		return false, err
	}

	if summary != nil {
		return false, nil
	}

	calculatedSummary, err := d.CalculateDailySummary(date)
	if err != nil {
		return false, err
	}

	if calculatedSummary != nil && calculatedSummary.AvgTemp > 0 {
		if err := d.InsertDailySummary(*calculatedSummary); err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}

func (d *Database) CreateSummariesForRange(start, end time.Time) ([]string, error) {
	var created []string
	for date := start; !date.After(end); date = date.AddDate(0, 0, 1) {
		wasCreated, err := d.CreateSummaryForDateIfNotExists(date)
		if err != nil {
			return created, err
		}
		if wasCreated {
			created = append(created, date.Format("2006-01-02"))
		}
	}
	return created, nil
}
