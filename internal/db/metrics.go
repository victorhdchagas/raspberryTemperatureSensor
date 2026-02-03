package db

import (
	"time"
)

func (d *Database) InsertMetric(temp, humidity float64) error {
	_, err := d.db.Exec(
		"INSERT INTO raw_metrics (temp, humidity) VALUES (?, ?)",
		temp, humidity,
	)
	return err
}

func (d *Database) GetLatestMetric() (*RawMetric, error) {
	row := d.db.QueryRow(
		"SELECT id, timestamp, temp, humidity FROM raw_metrics ORDER BY timestamp DESC LIMIT 1",
	)

	var metric RawMetric
	err := row.Scan(&metric.ID, &metric.Timestamp, &metric.Temp, &metric.Humidity)
	if err != nil {
		if err == nil {
			return nil, nil
		}
		return nil, err
	}

	return &metric, nil
}

func (d *Database) GetMetricsByDateRange(start, end time.Time) ([]RawMetric, error) {
	rows, err := d.db.Query(
		"SELECT id, timestamp, temp, humidity FROM raw_metrics WHERE timestamp BETWEEN ? AND ? ORDER BY timestamp ASC",
		start, end,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []RawMetric
	for rows.Next() {
		var metric RawMetric
		if err := rows.Scan(&metric.ID, &metric.Timestamp, &metric.Temp, &metric.Humidity); err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
	}

	return metrics, nil
}

func (d *Database) DeleteOldMetrics(days int) error {
	_, err := d.db.Exec(
		"DELETE FROM raw_metrics WHERE timestamp < datetime('now', '-' || ? || ' days')",
		days,
	)
	return err
}
