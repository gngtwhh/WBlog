package repository

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(datasourcePath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", datasourcePath)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	if err := createTable(db); err != nil {
		return nil, err
	}

	return db, nil
}

func createTable(db *sql.DB) error {
	const sqlStmt = `
		-- 1. Create articles Table
		CREATE TABLE IF NOT EXISTS articles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			author TEXT NOT NULL,
			content TEXT NOT NULL,
			abstract TEXT,
			view_count INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		-- 2. Create Trigger
		CREATE TRIGGER IF NOT EXISTS update_articles_updated_at
		AFTER UPDATE ON articles
		BEGIN
			UPDATE articles
			SET updated_at = CURRENT_TIMESTAMP
			WHERE id = NEW.id;
		END;
		`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("Error creating database tables: %v", err)
		return err
	}

	return nil
}
