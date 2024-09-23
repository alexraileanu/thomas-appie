package appie

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/alexraileanu/thomas-appie/pkg/logger"
	"time"

	"github.com/go-resty/resty/v2"
)

const URL = "https://www.ah.nl/gql"

//go:embed queryFormat.json
var queryFormatFile embed.FS

type Appie struct {
	loggerService *logger.Service
}

func New(loggerService *logger.Service) *Appie {
	return &Appie{loggerService: loggerService}
}

func (a *Appie) PerformProductsCheck(productsToWatch []Product) ([]Product, []Product, error) {
	var productsInBonus []Product
	var productsNotInBonus []Product

	for _, product := range productsToWatch {
		a.loggerService.Info(fmt.Sprintf("Checking product %s", product.FriendlyName), nil)
		bonusInfo, err := a.makeGqlRequest(product)
		if err != nil {
			return nil, nil, err
		}
		hasDiscount := bonusInfo.Data.Product.Price.Discount.SegmentId != 0
		productBonusInfo := DiscountedProducts{
			ProductID:   product.ID,
			InBonus:     hasDiscount,
			Description: bonusInfo.Data.Product.Price.Discount.Description,
			Label:       bonusInfo.Data.Product.SmartLabel,
		}

		productBonusInfo.InBonus = hasDiscount
		product.DiscountedProducts = append(product.DiscountedProducts, productBonusInfo)

		if hasDiscount {
			productsInBonus = append(productsInBonus, product)
		} else {
			productsNotInBonus = append(productsNotInBonus, product)
		}
	}

	return productsInBonus, productsNotInBonus, nil
}

func (a *Appie) makeGqlRequest(product Product) (ProductInfoResponse, error) {
	gqlQuery, err := a.readQueryFile()
	if err != nil {
		a.loggerService.Error("Error reading query file", map[string]interface{}{"error": err.Error()})
		return ProductInfoResponse{}, err
	}

	preparedRequest := fmt.Sprintf(gqlQuery, product.AppieId, time.Now().Format("2006-01-02"))

	client := resty.New()
	r := client.R()

	// we pretend we're a valid browser request
	r.Header.Add("client-name", "ah-products")
	r.Header.Add("client-version", "6.500.0")
	r.Header.Add("Referer", product.RefererUrl)
	r.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/117.0")
	r.Header.Add("Content-Type", "application/json")

	resp, err := r.SetBody(preparedRequest).Post(URL)
	if err != nil {
		a.loggerService.Error("Error making request to Appie", map[string]interface{}{"error": err.Error()})
		return ProductInfoResponse{}, err
	}

	productInfo := new(ProductInfoResponse)
	err = json.Unmarshal(resp.Body(), productInfo)
	if err != nil {
		a.loggerService.Error("Error unmarshalling response from Appie", map[string]interface{}{"error": err.Error()})
		return ProductInfoResponse{}, err

	}

	return *productInfo, nil
}

func (a *Appie) readQueryFile() (string, error) {
	data, err := queryFormatFile.ReadFile("queryFormat.json")
	if err != nil {
		return "", err
	}

	return string(data), nil
}
