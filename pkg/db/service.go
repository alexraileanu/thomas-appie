package db

import (
	"fmt"
	"time"

	"github.com/alexraileanu/thomas-appie/pkg/logger"

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
	monday := fmt.Sprintf("%s 00:00:00.000", getMonday())
	now := time.Now()

	var products []*appie.Product
	result := s.db.handler.
		Preload("DiscountedProducts", "created_at BETWEEN ? AND ?", monday, now).
		Joins("JOIN discounted_products ON discounted_products.product_id = products.id").
		Where("discounted_products.created_at BETWEEN ? AND ?", monday, now).
		Find(&products)
	if result.Error != nil {
		s.loggerService.Error("Error fetching discounted products from the db", map[string]interface{}{"error": result.Error.Error()})
		return nil, result.Error
	}

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
		var rows int64
		result := s.db.handler.Table("products").Where(appie.Product{AppieId: product.AppieId}).Count(&rows)
		if result.Error != nil {
			s.loggerService.Error("Error saving product", map[string]interface{}{"error": result.Error})
			return result.Error
		}

		if rows == 0 {
			// If the product does not exist, create it
			if err := s.db.handler.Create(&product).Error; err != nil {
				s.loggerService.Error("Error creating product", map[string]interface{}{"error": err.Error()})
				return err
			}
		} else {
			// If the product exists, update it
			if err := s.db.handler.Model(&appie.Product{}).Where("appie_id = ?", product.AppieId).Updates(product).Error; err != nil {
				s.loggerService.Error("Error updating product", map[string]interface{}{"error": err.Error()})
				return err
			}
		}
	}
	s.db.handler.Where("appie_id NOT IN (?)", pluckIds(products)).Debug().Delete(&appie.Product{})

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

func pluckIds(products []appie.Product) []int {
	ids := make([]int, len(products))
	for i, product := range products {
		ids[i] = product.AppieId
	}
	return ids
}
