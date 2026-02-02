package db

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

func New(path string) (*Database, error) {
	sqlDB, err := sql.Open("sqlite3", path+"?_foreign_keys=on&_journal_mode=WAL")
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	return &Database{db: sqlDB}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) Migrate() error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS raw_metrics (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			temp REAL NOT NULL,
			humidity REAL NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS daily_summaries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date DATE NOT NULL UNIQUE,
			avg_temp REAL,
			avg_humidity REAL,
			max_temp REAL,
			min_temp REAL
		)`,
		`CREATE TABLE IF NOT EXISTS user_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date DATE NOT NULL UNIQUE,
			rating INTEGER CHECK(rating >= 1 AND rating <= 5),
			note TEXT,
			feeling_tag TEXT
		)`,
		`CREATE INDEX IF NOT EXISTS idx_raw_metrics_timestamp ON raw_metrics(timestamp)`,
		`CREATE INDEX IF NOT EXISTS idx_daily_summaries_date ON daily_summaries(date)`,
		`CREATE INDEX IF NOT EXISTS idx_user_logs_date ON user_logs(date)`,
	}

	for _, migration := range migrations {
		if _, err := d.db.Exec(migration); err != nil {
			return err
		}
	}

	return nil
}

type RawMetric struct {
	ID        int64     `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Temp      float64   `json:"temp"`
	Humidity  float64   `json:"humidity"`
}

type DailySummary struct {
	ID          int64     `json:"id"`
	Date        time.Time `json:"date"`
	AvgTemp     float64   `json:"avg_temp"`
	AvgHumidity float64   `json:"avg_humidity"`
	MaxTemp     float64   `json:"max_temp"`
	MinTemp     float64   `json:"min_temp"`
}

type UserLog struct {
	ID         int64     `json:"id"`
	Date       time.Time `json:"date"`
	Rating     int       `json:"rating"`
	Note       string    `json:"note"`
	FeelingTag string    `json:"feeling_tag"`
}

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
		if err == sql.ErrNoRows {
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

func (d *Database) InsertUserLog(log UserLog) error {
	_, err := d.db.Exec(
		"INSERT OR REPLACE INTO user_logs (date, rating, note, feeling_tag) VALUES (?, ?, ?, ?)",
		log.Date, log.Rating, log.Note, log.FeelingTag,
	)
	return err
}

func (d *Database) GetUserLogByDate(date time.Time) (*UserLog, error) {
	row := d.db.QueryRow(
		"SELECT id, date, rating, note, feeling_tag FROM user_logs WHERE date = ?",
		date,
	)

	var log UserLog
	err := row.Scan(&log.ID, &log.Date, &log.Rating, &log.Note, &log.FeelingTag)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &log, nil
}

func (d *Database) DeleteOldMetrics(days int) error {
	_, err := d.db.Exec(
		"DELETE FROM raw_metrics WHERE timestamp < datetime('now', '-' || ? || ' days')",
		days,
	)
	return err
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
