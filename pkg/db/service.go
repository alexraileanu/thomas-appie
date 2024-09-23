package db

import (
	"fmt"
	"github.com/alexraileanu/thomas-appie/pkg/logger"
	"time"

	"github.com/alexraileanu/thomas-appie/pkg/appie"
)

type Service struct {
	db            *DB
	loggerService *logger.Service
}

func NewDBService(db *DB, loggerService *logger.Service) *Service {
	return &Service{db: db, loggerService: loggerService}
}

func (s *Service) GetProducts() ([]appie.Product, error) {
	var products []appie.Product
	result := s.db.handler.Table("products").Find(&products)
	if result.Error != nil {
		s.loggerService.Error("Error fetching products from the db", map[string]interface{}{"error": result.Error.Error()})
		return nil, result.Error
	}

	return products, nil
}

func (s *Service) GetDiscountedProductsThisWeek() ([]*appie.Product, error) {
	// getMonday returns only the date but created_at contains the time as well.
	// for this we can consider the start of the day
	monday := fmt.Sprintf("%s 00:00:00.000", getMonday())

	var products []*appie.Product
	s.db.handler.Joins("JOIN discounted_products ON products.id = discounted_products.product_id").
		Where("discounted_products.created_at > ?", monday).
		Preload("DiscountedProducts").
		Find(&products)

	for _, product := range products {
		product.Discount = product.DiscountedProducts[0]
	}

	return products, nil
}

func (s *Service) SaveDiscountedProducts(products []appie.Product) error {
	monday := fmt.Sprintf("%s 00:00:00.000", getMonday())
	now := time.Now()
	for _, product := range products {
		s.db.handler.Where("product_id = ? AND created_at BETWEEN ? AND ?", product.ID, monday, now).FirstOrCreate(&product.DiscountedProducts[0])
	}

	return nil
}

func (s *Service) SaveProduct(products []appie.Product) error {
	for _, product := range products {
		result := s.db.handler.Where(appie.Product{AppieId: product.AppieId}).FirstOrCreate(&product)
		if result.Error != nil {
			s.loggerService.Error("Error saving product", map[string]interface{}{"error": result.Error.Error()})
			return result.Error
		}
	}

	return nil
}

func getMonday() string {
	// Get the current date
	today := time.Now()

	// Calculate the number of days to subtract to get to Monday (considering Sunday as the first day of the week)
	daysToSubtract := int(today.Weekday() - time.Monday)
	if daysToSubtract < 0 {
		daysToSubtract += 7 // Add 7 to loop back to Monday of the previous week
	}

	return today.AddDate(0, 0, -daysToSubtract).Format("2006-01-02")
}
