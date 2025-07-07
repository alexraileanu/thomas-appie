package product

import (
	"github.com/alexraileanu/thomas-appie/pkg/appie"
	"github.com/alexraileanu/thomas-appie/pkg/utl"
)

func New() *Service {
	return &Service{}
}

type Service struct {
}

func (s *Service) GetAllProducts() ([]appie.Product, error) {
	return utl.ParseProductsJson()
}

func (s *Service) UpdateProducts(products []appie.Product) error {
	err := utl.UpdateProductsJson(products)
	if err != nil {
		return err
	}

	return nil
}
