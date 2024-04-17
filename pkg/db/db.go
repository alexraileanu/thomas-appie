package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	handler *sql.DB
}

func New() (*DB, error) {
	handler, err := sql.Open("sqlite3", "./data/thomas.db?cache=shared")
	if err != nil {
		return nil, err
	}

	_, err = handler.Exec("CREATE TABLE IF NOT EXISTS products (id INTEGER PRIMARY KEY, api_name TEXT, friendly_name TEXT, in_bonus TINYINT)")
	if err != nil {
		return nil, err
	}

	return &DB{
		handler: handler,
	}, nil
}
