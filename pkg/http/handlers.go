package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) containersListHandler(c echo.Context) error {
	products, err := s.dbService.GetProducts()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, products)
}
