package utl

import (
	"encoding/json"
	"os"

	"github.com/alexraileanu/thomas-appie/pkg/appie"
	"github.com/alexraileanu/thomas-appie/pkg/logger"
)

func ParseProductsJson(loggerService *logger.Service) ([]appie.Product, error) {
	path := os.Getenv("PRODUCTS_JSON_FILE_PATH")
	loggerService.Debug("Reading products JSON", map[string]interface{}{"path": path})

	fileContents, err := os.ReadFile(path)
	if err != nil {
		loggerService.Error("Failed to read products JSON", map[string]interface{}{"path": path, "error": err.Error()})
		return nil, err
	}

	products := new([]appie.Product)
	err = json.Unmarshal(fileContents, products)
	if err != nil {
		loggerService.Error("Failed to parse products JSON", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	loggerService.Debug("Parsed products JSON", map[string]interface{}{"count": len(*products)})
	return *products, nil
}

func UpdateProductsJson(products []appie.Product, loggerService *logger.Service) error {
	path := os.Getenv("PRODUCTS_JSON_FILE_PATH")
	loggerService.Debug("Writing products JSON", map[string]interface{}{"path": path, "count": len(products)})

	fileContents, err := json.MarshalIndent(products, "", "  ")
	if err != nil {
		loggerService.Error("Failed to marshal products JSON", map[string]interface{}{"error": err.Error()})
		return err
	}

	err = os.WriteFile(path, fileContents, 0644)
	if err != nil {
		loggerService.Error("Failed to write products JSON", map[string]interface{}{"path": path, "error": err.Error()})
		return err
	}

	return nil
}
