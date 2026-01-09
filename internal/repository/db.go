package repository

import (
	"database/sql"

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
	return db, nil
}
