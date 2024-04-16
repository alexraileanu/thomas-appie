package utl

import (
	"encoding/json"
	"os"

	"github.com/alexraileanu/thomas-appie/pkg/appie"
)

func ParseProductsJson() ([]*appie.ProductToCheck, error) {
	fileContents, err := os.ReadFile(os.Getenv("PRODUCTS_JSON_FILE_PATH"))
	if err != nil {
		return nil, err
	}

	products := new([]*appie.ProductToCheck)
	err = json.Unmarshal(fileContents, products)
	if err != nil {
		return nil, err
	}

	return *products, nil
}
