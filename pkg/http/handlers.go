package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/alexraileanu/thomas-appie/pkg/appie"
)

//type updateProductsRequest struct {
//	ApiName      string `json:"api_name"`
//	AppieId      int    `json:"appie_id"`
//	FriendlyName string `json:"friendly_name"`
//	RefererUrl   string `json:"referer_url"`
//}

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
