package appie

import (
	"embed"
	"encoding/json"
	"fmt"
	"time"

	"github.com/alexraileanu/thomas-appie/pkg/config"
	"github.com/alexraileanu/thomas-appie/pkg/logger"

	"github.com/go-resty/resty/v2"
)

const URL = "https://www.ah.nl/gql"

//go:embed queryFormat.json
var queryFormatFile embed.FS

type Appie struct {
	loggerService *logger.Service
	config        config.Appie
}

func New(loggerService *logger.Service, config config.Appie) *Appie {
	return &Appie{loggerService: loggerService, config: config}
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
			Label:       bonusInfo.Data.Product.Title,
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
	r.Header.Add("client-name", a.config.ClientName)
	r.Header.Add("client-version", a.config.ClientVersion)
	r.Header.Add("x-client-platform-type", a.config.ClientPlatformType)
	r.Header.Add("User-Agent", a.config.UserAgent)
	r.Header.Add("Referer", product.RefererUrl)
	r.Header.Add("Content-Type", "application/json")

	a.loggerService.Info("Making request to Appie", map[string]interface{}{"request": preparedRequest})

	resp, err := r.SetBody(preparedRequest).Post(URL)
	if err != nil {
		a.loggerService.Error("Error making request to Appie", map[string]interface{}{"error": err.Error()})
		return ProductInfoResponse{}, err
	}

	if resp.StatusCode() != 200 {
		a.loggerService.Error("Got err response from Appie", map[string]interface{}{"response_headers": resp.Header(), "status_code": resp.Header()})
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
