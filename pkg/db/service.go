package db

import (
	"github.com/alexraileanu/thomas-appie/pkg/appie"
)

type Service struct {
	db *DB
}

func NewDBService(db *DB) *Service {
	return &Service{db: db}
}

func (s *Service) GetProducts() ([]appie.ProductToCheck, error) {
	var products []appie.ProductToCheck
	rows, err := s.db.handler.Query("SELECT friendly_name, api_name, in_bonus FROM products")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var p appie.ProductToCheck
		err = rows.Scan(&p.FriendlyName, &p.ApiName, &p.InBonus)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (db *DB) Save(products []*appie.ProductToCheck) error {
	tx, err := db.handler.Begin()
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

	stmt, err = tx.Prepare("INSERT INTO products(api_name, friendly_name, in_bonus) VALUES(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, p := range products {
		_, err = stmt.Exec(p.ApiName, p.FriendlyName, p.InBonus)
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
