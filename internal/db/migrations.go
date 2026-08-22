package db

func GetMigrations() []string {
	return []string{
		`CREATE TABLE IF NOT EXISTS raw_metrics (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			temp REAL NOT NULL,
			humidity REAL NOT NULL,
			temp_ext REAL
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
}
