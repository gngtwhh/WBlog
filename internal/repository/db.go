package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dsn string) (*sql.DB, error) {
	// create parent dir if not exists
	dbPath := dsn
	if strings.HasPrefix(dsn, "file:") {
		dbPath = dsn[5:]
	}
	if idx := strings.Index(dbPath, "?"); idx != -1 {
		dbPath = dbPath[:idx]
	}
	dir := filepath.Dir(dbPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	db, err := sql.Open("sqlite3", dsn)
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
	const schema = `
	-- -----------------------------------------------------
	-- 1. Users
	-- -----------------------------------------------------
	CREATE TABLE IF NOT EXISTS users (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		username   TEXT NOT NULL UNIQUE,
		password   TEXT NOT NULL,
		nickname   TEXT DEFAULT '',
		avatar     TEXT DEFAULT '',
		role       INTEGER DEFAULT 1,   -- 1: User, 99: Admin
		status     INTEGER DEFAULT 1,   -- 1: Normal, 0: Banned
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	-- Auto-update updated_at on change
	CREATE TRIGGER IF NOT EXISTS trg_users_updated_at
	AFTER UPDATE ON users
	BEGIN
		UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
	END;

	-- -----------------------------------------------------
	-- 2. Articles
	-- -----------------------------------------------------
	CREATE TABLE IF NOT EXISTS articles (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		title      TEXT NOT NULL,
		author     TEXT NOT NULL,
		content    TEXT NOT NULL,
		abstract   TEXT DEFAULT '',
		view_count INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TRIGGER IF NOT EXISTS trg_articles_updated_at
	AFTER UPDATE ON articles
	BEGIN
		UPDATE articles SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
	END;

	-- -----------------------------------------------------
	-- 3. Comments
	-- -----------------------------------------------------
	CREATE TABLE IF NOT EXISTS comments (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id    INTEGER NOT NULL,
		article_id INTEGER NOT NULL,
		username   TEXT NOT NULL,     -- Denormalized for read performance
		content    TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TRIGGER IF NOT EXISTS trg_comments_updated_at
	AFTER UPDATE ON comments
	BEGIN
		UPDATE comments SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
	END;

	-- -----------------------------------------------------
	-- 4. Indices
	-- -----------------------------------------------------
	CREATE INDEX IF NOT EXISTS idx_comments_article_id ON comments(article_id);
	CREATE INDEX IF NOT EXISTS idx_articles_created_at ON articles(created_at DESC);
	`

	if _, err := db.Exec(schema); err != nil {
		log.Printf("Init database schema failed: %v", err)
		return err
	}

	return nil
}
