package db

import (
	"fmt"
	"time"

	"github.com/alexraileanu/thomas-appie/pkg/config"
	"github.com/alexraileanu/thomas-appie/pkg/logger"

	"github.com/alexraileanu/thomas-appie/pkg/appie"
)

type Service struct {
	db            *DB
	loggerService *logger.Service
	cfg           config.Appie
}

func NewDBService(db *DB, loggerService *logger.Service, cfg config.Appie) *Service {
	return &Service{db: db, loggerService: loggerService, cfg: cfg}
}

func (s *Service) GetProducts() ([]appie.Product, error) {
	s.loggerService.Debug("Querying all products", nil)
	var products []appie.Product
	result := s.db.handler.Table("products").Find(&products)
	if result.Error != nil {
		s.loggerService.Error("Error fetching products from the db", map[string]interface{}{"error": result.Error.Error()})
		return nil, result.Error
	}

	s.loggerService.Info("Fetched products from the db", map[string]interface{}{"count": len(products)})
	return products, nil
}

func (s *Service) GetDiscountedProductsThisWeek() ([]*appie.Product, error) {
	monday := fmt.Sprintf("%s 00:00:00.000", getBonusDay(s.cfg.BonusDay))
	now := time.Now()

	s.loggerService.Debug("Querying discounted products", map[string]interface{}{"from": monday, "to": now.Format("2006-01-02 15:04:05")})
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

	s.loggerService.Debug("Fetched discounted products", map[string]interface{}{"count": len(products)})
	for _, product := range products {
		product.Discount = product.DiscountedProducts[0]
	}

	return products, nil
}

func (s *Service) SaveDiscountedProducts(products []appie.Product) error {
	monday := fmt.Sprintf("%s 00:00:00.000", getBonusDay(s.cfg.BonusDay))
	now := time.Now()

	s.loggerService.Debug("Saving discounted products", map[string]interface{}{"count": len(products)})
	for _, product := range products {
		s.loggerService.Debug("Upserting discounted product", map[string]interface{}{"product": product.FriendlyName, "in_bonus": product.DiscountedProducts[0].InBonus})
		s.db.handler.Where("product_id = ? AND created_at BETWEEN ? AND ?", product.ID, monday, now).FirstOrCreate(&product.DiscountedProducts[0])
	}

	return nil
}

func (s *Service) SaveProduct(products []appie.Product) error {
	s.loggerService.Debug("Saving products", map[string]interface{}{"count": len(products)})
	for _, product := range products {
		var rows int64
		result := s.db.handler.Table("products").Where(appie.Product{AppieId: product.AppieId}).Count(&rows)
		if result.Error != nil {
			s.loggerService.Error("Error saving product", map[string]interface{}{"error": result.Error})
			return result.Error
		}

		if rows == 0 {
			s.loggerService.Debug("Creating product", map[string]interface{}{"product": product.FriendlyName, "appie_id": product.AppieId})
			if err := s.db.handler.Create(&product).Error; err != nil {
				s.loggerService.Error("Error creating product", map[string]interface{}{"error": err.Error()})
				return err
			}
		} else {
			s.loggerService.Debug("Updating product", map[string]interface{}{"product": product.FriendlyName, "appie_id": product.AppieId})
			if err := s.db.handler.Model(&appie.Product{}).Where("appie_id = ?", product.AppieId).Updates(product).Error; err != nil {
				s.loggerService.Error("Error updating product", map[string]interface{}{"error": err.Error()})
				return err
			}
		}
	}
	s.loggerService.Debug("Deleting stale products", map[string]interface{}{"keeping_ids": pluckIds(products)})
	s.db.handler.Where("appie_id NOT IN (?)", pluckIds(products)).Debug().Delete(&appie.Product{})

	return nil
}

func getBonusDay(bonusDay int) string {
	// Get the current date
	today := time.Now()
	bonusWeekDay := time.Weekday(bonusDay)

	// Calculate the number of days to subtract to get to bonus day (normally monday but can be configured to be something else)
	daysToSubtract := int(today.Weekday() - bonusWeekDay)
	if daysToSubtract < 0 {
		daysToSubtract += 7 // Add 7 to loop back to bonus day of the previous week
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
