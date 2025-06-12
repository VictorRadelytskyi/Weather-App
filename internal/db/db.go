package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./feedback.db")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS feedback(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			satisfaction TEXT DEFAULT "",
			feedback TEXT DEFAULT "",
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)

	if err != nil {
		return nil, err
	}

	return db, nil
}
