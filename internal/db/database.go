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

	return d.ensureTempExtColumn()
}

// ensureTempExtColumn adiciona a coluna temp_ext em bancos que já existiam
// antes dessa feature. Bancos novos já a criam no CREATE TABLE.
func (d *Database) ensureTempExtColumn() error {
	rows, err := d.db.Query(`PRAGMA table_info(raw_metrics)`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dflt sql.NullString
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			return err
		}
		if name == "temp_ext" {
			return nil // já existe
		}
	}

	_, err = d.db.Exec(`ALTER TABLE raw_metrics ADD COLUMN temp_ext REAL`)
	return err
}
