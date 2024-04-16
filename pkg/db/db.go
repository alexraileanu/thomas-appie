package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"github.com/alexraileanu/thomas-appie/pkg/appie"
)

type DB struct {
	d *sql.DB
}

func New() (*DB, error) {
	d, err := sql.Open("sqlite3", "./data/thomas.db?cache=shared")
	if err != nil {
		return nil, err
	}

	_, err = d.Exec("CREATE TABLE IF NOT EXISTS products (id INTEGER PRIMARY KEY, name TEXT, in_bonus TINYINT)")
	if err != nil {
		return nil, err
	}

	return &DB{
		d: d,
	}, nil
}

func (db *DB) Save(products []*appie.ProductToCheck) error {
	tx, err := db.d.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("DELETE FROM products")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	stmt, err = tx.Prepare("INSERT INTO products(name, in_bonus) VALUES(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, p := range products {
		name := fmt.Sprintf("%s (%s)", p.FriendlyName, p.ApiName)
		_, err = stmt.Exec(name, p.HasDiscount)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
