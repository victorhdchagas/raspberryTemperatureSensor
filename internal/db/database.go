package db

import (
	"database/sql"

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
	migrations := GetMigrations()

	for _, migration := range migrations {
		if _, err := d.db.Exec(migration); err != nil {
			return err
		}
	}

	return nil
}
