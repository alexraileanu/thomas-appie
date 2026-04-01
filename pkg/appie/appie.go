package appie

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/alexraileanu/thomas-appie/pkg/config"
	"github.com/alexraileanu/thomas-appie/pkg/logger"

	"github.com/go-resty/resty/v2"
)

const GQLUrl = "https://www.ah.nl/gql"
const RestURL = "https://api.ah.nl"

//go:embed queryFormat.json
var queryFormatFile embed.FS

type Appie struct {
	loggerService *logger.Service
	config        config.Appie

	authToken string
}

func New(loggerService *logger.Service, config config.Appie) *Appie {
	return &Appie{loggerService: loggerService, config: config}
}

func (a *Appie) PerformProductsCheck(productsToWatch []Product) ([]Product, []Product, error) {
	if a.config.V2 {
		a.loggerService.Info("V2 active", nil)
		err := a.getAnonAuthToken()
		if err != nil {
			return nil, nil, err
		}
	}

	var productsInBonus []Product
	var productsNotInBonus []Product

	for _, product := range productsToWatch {
		a.loggerService.Info(fmt.Sprintf("Checking product %s", product.FriendlyName), nil)

		if a.config.V2 {
			bonusInfo, err := a.makeRequest(product)
			if err != nil {
				return nil, nil, err
			}

			hasBonus := bonusInfo.ProductCard.IsBonus
			productBonusInfo := DiscountedProducts{
				ProductID:   product.ID,
				InBonus:     bonusInfo.ProductCard.IsBonus,
				Description: bonusInfo.ProductCard.BonusMechanism,
			}
			product.DiscountedProducts = append(product.DiscountedProducts, productBonusInfo)
			if len(bonusInfo.ProductCard.Images) > 0 {
				product.Image = bonusInfo.ProductCard.Images[0].URL
			}
			if hasBonus {
				productsInBonus = append(productsInBonus, product)
			} else {
				productsNotInBonus = append(productsNotInBonus, product)
			}
		} else {
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
	}

	return productsInBonus, productsNotInBonus, nil
}

func (a *Appie) makeBaseRequest(hasAuth bool, r *http.Request) error {
	r.Header.Add("User-Agent", a.config.UserAgent)
	r.Header.Add("Client-Version", a.config.ClientVersion)
	r.Header.Add("X-Application", a.config.XApplication)
	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/json")
	if hasAuth {
		r.Header.Add("Authorization", "Bearer "+a.authToken)
	}

	return nil
}

func (a *Appie) getAnonAuthToken() error {
	a.loggerService.Info("Getting Anon Auth Token", nil)

	body := fmt.Sprintf(`{"clientId": "%s"}`, a.config.ClientName)

	r, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", RestURL, "mobile-auth/v1/auth/token/anonymous"), io.NopCloser(bytes.NewBufferString(body)))
	if err != nil {
		return err
	}
	err = a.makeBaseRequest(false, r)
	if err != nil {
		return err
	}
	c := http.Client{}

	resp, err := c.Do(r)
	if err != nil {
		a.loggerService.Error("Anon auth token request failed", map[string]interface{}{"error": err.Error()})
		return err
	}
	defer resp.Body.Close()

	var authResponse AnonAuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResponse)
	if err != nil {
		a.loggerService.Error("Anon auth token decoding failed", map[string]interface{}{"error": err.Error()})
		return err
	}

	a.loggerService.Info("Got anon auth token", nil)
	a.loggerService.Debug("Anon auth token", map[string]interface{}{"token": authResponse.AccessToken})
	a.authToken = authResponse.AccessToken

	return nil
}

func (a *Appie) makeRequest(product Product) (ProductInfoResponse, error) {
	url := fmt.Sprintf("%s/%s/%d", RestURL, "mobile-services/product/detail/v4/fir", product.AppieId)
	a.loggerService.Info("Fetching product info", map[string]interface{}{"product": product.FriendlyName})
	a.loggerService.Debug("Product request URL", map[string]interface{}{"product": product.FriendlyName, "url": url})
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ProductInfoResponse{}, err
	}

	err = a.makeBaseRequest(true, r)
	if err != nil {
		a.loggerService.Error("Failed to build product request", map[string]interface{}{"product": product.FriendlyName, "error": err.Error()})
		return ProductInfoResponse{}, err
	}
	c := http.Client{}

	resp, err := c.Do(r)
	if err != nil {
		a.loggerService.Error("Product info request failed", map[string]interface{}{"product": product.FriendlyName, "error": err.Error()})
		return ProductInfoResponse{}, err
	}
	defer resp.Body.Close()

	var productInfo ProductInfoResponse
	err = json.NewDecoder(resp.Body).Decode(&productInfo)
	if err != nil {
		a.loggerService.Error("Failed to decode product info response", map[string]interface{}{"product": product.FriendlyName, "error": err.Error()})
		return ProductInfoResponse{}, err
	}

	a.loggerService.Debug("Product info response", map[string]interface{}{"product": product.FriendlyName, "is_bonus": productInfo.ProductCard.IsBonus, "bonus_mechanism": productInfo.ProductCard.BonusMechanism})
	return productInfo, nil
}

func (a *Appie) makeGqlRequest(product Product) (GQLProductInfoResponse, error) {
	gqlQuery, err := a.readQueryFile()
	if err != nil {
		a.loggerService.Error("Error reading query file", map[string]interface{}{"error": err.Error()})
		return GQLProductInfoResponse{}, err
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

	a.loggerService.Info("Fetching product info", map[string]interface{}{"product": product.FriendlyName})
	a.loggerService.Debug("GQL request body", map[string]interface{}{"product": product.FriendlyName, "body": preparedRequest})

	resp, err := r.SetBody(preparedRequest).Post(GQLUrl)
	if err != nil {
		a.loggerService.Error("Product info request failed", map[string]interface{}{"product": product.FriendlyName, "error": err.Error()})
		return GQLProductInfoResponse{}, err
	}

	if resp.StatusCode() != 200 {
		a.loggerService.Error("Unexpected response from Appie", map[string]interface{}{"product": product.FriendlyName, "status_code": resp.StatusCode()})
		a.loggerService.Debug("Error response body", map[string]interface{}{"product": product.FriendlyName, "body": string(resp.Body())})
		return GQLProductInfoResponse{}, err
	}

	productInfo := new(GQLProductInfoResponse)
	err = json.Unmarshal(resp.Body(), productInfo)
	if err != nil {
		a.loggerService.Error("Failed to decode product info response", map[string]interface{}{"product": product.FriendlyName, "error": err.Error()})
		return GQLProductInfoResponse{}, err

	}

	a.loggerService.Debug("Product info response", map[string]interface{}{"product": product.FriendlyName, "has_discount": productInfo.Data.Product.Price.Discount.SegmentId != 0, "discount_description": productInfo.Data.Product.Price.Discount.Description})
	return *productInfo, nil
}

func (a *Appie) readQueryFile() (string, error) {
	data, err := queryFormatFile.ReadFile("queryFormat.json")
	if err != nil {
		return "", err
	}

	return string(data), nil
}
