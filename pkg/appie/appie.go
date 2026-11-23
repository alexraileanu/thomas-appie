package appie

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/alexraileanu/thomas-appie/pkg/config"
	"github.com/alexraileanu/thomas-appie/pkg/logger"
)

const (
	restURL    = "https://api.ah.nl"
	graphQLURL = restURL + "/graphql"
)

// productDetailsQuery is a trimmed copy of the Appie app's FetchProductDetails
// query (see gwillem/appie-go products.go). It replaces
// /mobile-services/product/detail/v4/fir/{id}, which returns 404 since 2026-11.
// source: https://github.com/gwillem/appie-go/blob/main/products.go#L75-L77
const productDetailsQuery = `query FetchProductDetails($productId: Int!) {
  product(id: $productId) {
    id
    title
    priceV2 {
      discount { description }
    }
    imagePack {
      large { url }
      medium { url }
      small { url }
    }
  }
}`

type Appie struct {
	loggerService *logger.Service
	config        config.Appie

	authToken string
}

func New(loggerService *logger.Service, config config.Appie) *Appie {
	return &Appie{loggerService: loggerService, config: config}
}

func (a *Appie) PerformProductsCheck(products []Product) ([]Product, []Product, error) {
	if err := a.getAnonAuthToken(); err != nil {
		return nil, nil, err
	}

	var inBonus, notInBonus []Product

	for _, product := range products {
		a.loggerService.Info(fmt.Sprintf("Checking product %s", product.FriendlyName), nil)

		info, err := a.makeRequest(product)
		if err != nil {
			return nil, nil, err
		}

		isBonus := info.IsBonus()
		product.DiscountedProducts = append(product.DiscountedProducts, DiscountedProducts{
			ProductID:   product.ID,
			InBonus:     isBonus,
			Description: info.BonusMechanism(),
			Label:       info.Title,
		})
		if image := info.ImageURL(); image != "" {
			product.Image = image
		}

		if isBonus {
			inBonus = append(inBonus, product)
		} else {
			notInBonus = append(notInBonus, product)
		}
		a.loggerService.Info("Fetched product info", map[string]any{
			"isBonus":        isBonus,
			"bonusMechanism": info.BonusMechanism(),
			"image":          info.ImageURL(),
			"title":          product.FriendlyName,
		})
	}

	return inBonus, notInBonus, nil
}

func (a *Appie) addBaseHeaders(hasAuth bool, r *http.Request) {
	r.Header.Add("User-Agent", a.config.UserAgent)
	r.Header.Add("Client-Version", a.config.ClientVersion)
	r.Header.Add("X-Application", a.config.XApplication)
	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/json")
	if hasAuth {
		r.Header.Add("Authorization", "Bearer "+a.authToken)
	}
}

func (a *Appie) getAnonAuthToken() error {
	a.loggerService.Info("Getting anon auth token", nil)

	body := fmt.Sprintf(`{"clientId": "%s"}`, a.config.ClientName)
	url := fmt.Sprintf("%s/mobile-auth/v1/auth/token/anonymous", restURL)

	r, err := http.NewRequest("POST", url, io.NopCloser(bytes.NewBufferString(body)))
	if err != nil {
		return err
	}
	a.addBaseHeaders(false, r)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		a.loggerService.Error("Anon auth token request failed", map[string]any{"error": err.Error()})
		return err
	}
	defer resp.Body.Close()

	var authResponse AnonAuthResponse
	if err = json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		a.loggerService.Error("Anon auth token decoding failed", map[string]any{"error": err.Error()})
		return err
	}

	a.loggerService.Info("Got anon auth token", nil)
	a.loggerService.Debug("Anon auth token", map[string]any{"token": authResponse.AccessToken})
	a.authToken = authResponse.AccessToken

	return nil
}

func (a *Appie) makeRequest(product Product) (*ProductDetails, error) {
	a.loggerService.Info("Fetching product info", map[string]any{"product": product.FriendlyName})
	a.loggerService.Debug("Product request", map[string]any{"product": product.FriendlyName, "url": graphQLURL, "appie_id": product.AppieId})

	body, err := json.Marshal(graphQLRequest{
		Query:     productDetailsQuery,
		Variables: map[string]any{"productId": product.AppieId},
	})
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest("POST", graphQLURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	a.addBaseHeaders(true, r)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		a.loggerService.Error("Product info request failed", map[string]any{"product": product.FriendlyName, "error": err.Error()})
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		a.loggerService.Error("Failed to read product info response", map[string]any{"product": product.FriendlyName, "error": err.Error()})
		return nil, err
	}

	if resp.StatusCode >= 400 {
		err = fmt.Errorf("product info request for %s returned %d: %s", product.FriendlyName, resp.StatusCode, string(respBody))
		a.loggerService.Error("Product info request returned an error status", map[string]any{"product": product.FriendlyName, "status": resp.StatusCode, "body": string(respBody)})
		return nil, err
	}

	var info ProductDetailsResponse
	if err = json.Unmarshal(respBody, &info); err != nil {
		a.loggerService.Error("Failed to decode product info response", map[string]any{"product": product.FriendlyName, "error": err.Error()})
		return nil, err
	}

	if len(info.Errors) > 0 {
		err = fmt.Errorf("graphql error for %s: %s", product.FriendlyName, info.Errors[0].Message)
		a.loggerService.Error("Product info GraphQL error", map[string]any{"product": product.FriendlyName, "error": info.Errors[0].Message})
		return nil, err
	}

	if info.Data.Product == nil {
		err = fmt.Errorf("product %s (appie id %d) not found", product.FriendlyName, product.AppieId)
		a.loggerService.Error("Product not found", map[string]any{"product": product.FriendlyName, "appie_id": product.AppieId})
		return nil, err
	}

	details := info.Data.Product
	a.loggerService.Debug("Product info response", map[string]any{
		"product":         product.FriendlyName,
		"is_bonus":        details.IsBonus(),
		"bonus_mechanism": details.BonusMechanism(),
	})
	return details, nil
}
