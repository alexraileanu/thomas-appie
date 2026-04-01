package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/alexraileanu/thomas-appie/pkg/appie"
	"github.com/alexraileanu/thomas-appie/pkg/logger"
)

type DB struct {
	handler *gorm.DB
}

func New(user string, password string, host string, port string, dbName string, loggerService *logger.Service) (*DB, error) {
	loggerService.Debug("Connecting to database", map[string]interface{}{"host": host, "port": port, "db": dbName})
	dbDsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True", user, password, host, port, dbName)
	handler, err := gorm.Open(mysql.Open(dbDsn))
	if err != nil {
		loggerService.Error("Failed to connect to database", map[string]interface{}{"host": host, "port": port, "error": err.Error()})
		return nil, err
	}

	loggerService.Debug("Running auto-migrations", nil)
	err = handler.AutoMigrate(&appie.Product{}, &appie.DiscountedProducts{})
	if err != nil {
		loggerService.Error("Auto-migration failed", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	loggerService.Info("Connected to database", map[string]interface{}{"host": host, "port": port, "db": dbName})
	return &DB{
		handler: handler,
	}, nil
}
