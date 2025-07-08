package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/alexraileanu/thomas-appie/pkg/appie"
)

func (s *Server) refreshProducts(c echo.Context) error {
	// Fetch products from the Appie

	p, err := s.productService.GetAllProducts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	a := appie.New(s.loggerService, s.conf)
	inBonus, notInBonus, err := a.PerformProductsCheck(p)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	_ = s.dbService.SaveDiscountedProducts(append(inBonus, notInBonus...))

	discounts, err := s.dbService.GetDiscountedProductsThisWeek()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, discounts)
}

func (s *Server) getDiscountedProducts(c echo.Context) error {
	products, err := s.dbService.GetDiscountedProductsThisWeek()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, products)
}

func (s *Server) getAllProducts(c echo.Context) error {
	products, err := s.dbService.GetProducts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, products)
}

func (s *Server) updateProducts(c echo.Context) error {
	r := new([]appie.Product)
	if err := c.Bind(r); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Save the products to the database
	if err := s.dbService.SaveProduct(*r); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
