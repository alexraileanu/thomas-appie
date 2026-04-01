package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/alexraileanu/thomas-appie/pkg/appie"
)

func (s *Server) refreshProducts(c echo.Context) error {
	s.loggerService.Info("Refreshing products", nil)

	p, err := s.dbService.GetProducts()
	if err != nil {
		s.loggerService.Error("Failed to fetch products for refresh", map[string]interface{}{"error": err.Error()})
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	a := appie.New(s.loggerService, s.conf)
	inBonus, notInBonus, err := a.PerformProductsCheck(p)
	if err != nil {
		s.loggerService.Error("Product check failed during refresh", map[string]interface{}{"error": err.Error()})
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	_ = s.dbService.SaveProduct(append(inBonus, notInBonus...))

	discounts, err := s.dbService.GetDiscountedProductsThisWeek()
	if err != nil {
		s.loggerService.Error("Failed to fetch discounts after refresh", map[string]interface{}{"error": err.Error()})
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	s.loggerService.Info("Products refreshed", map[string]interface{}{"in_bonus": len(inBonus), "not_in_bonus": len(notInBonus)})
	return c.JSON(http.StatusOK, discounts)
}

func (s *Server) getDiscountedProducts(c echo.Context) error {
	s.loggerService.Debug("Fetching discounted products this week", nil)
	products, err := s.dbService.GetDiscountedProductsThisWeek()
	if err != nil {
		s.loggerService.Error("Failed to fetch discounted products", map[string]interface{}{"error": err.Error()})
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	s.loggerService.Debug("Fetched discounted products", map[string]interface{}{"count": len(products)})
	return c.JSON(http.StatusOK, products)
}

func (s *Server) getAllProducts(c echo.Context) error {
	s.loggerService.Debug("Fetching all products from db", nil)
	products, err := s.dbService.GetProducts()
	if err != nil {
		s.loggerService.Error("Failed to fetch all products", map[string]interface{}{"error": err.Error()})
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	s.loggerService.Debug("Fetched all products", map[string]interface{}{"count": len(products)})
	return c.JSON(http.StatusOK, products)
}

func (s *Server) updateProducts(c echo.Context) error {
	r := new([]appie.Product)
	if err := c.Bind(r); err != nil {
		s.loggerService.Error("Invalid request body for product update", map[string]interface{}{"error": err.Error()})
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	s.loggerService.Info("Updating products", map[string]interface{}{"count": len(*r)})
	if err := s.dbService.SaveProduct(*r); err != nil {
		s.loggerService.Error("Failed to save products", map[string]interface{}{"error": err.Error()})
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
