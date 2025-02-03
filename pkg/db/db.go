package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/alexraileanu/thomas-appie/pkg/appie"
)

type DB struct {
	handler *gorm.DB
}

func New(user string, password string, host string, port string, dbName string) (*DB, error) {
	dbDsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True", user, password, host, port, dbName)
	handler, err := gorm.Open(mysql.Open(dbDsn))
	if err != nil {
		return nil, err
	}

	err = handler.AutoMigrate(&appie.Product{}, &appie.DiscountedProducts{})
	if err != nil {
		return nil, err
	}

	return &DB{
		handler: handler,
	}, nil
}
