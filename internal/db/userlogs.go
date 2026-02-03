package db

import (
	"database/sql"
	"time"
)

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
