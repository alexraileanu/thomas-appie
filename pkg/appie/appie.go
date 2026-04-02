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

const restURL = "https://api.ah.nl"

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

		product.DiscountedProducts = append(product.DiscountedProducts, DiscountedProducts{
			ProductID:   product.ID,
			InBonus:     info.ProductCard.IsBonus,
			Description: info.ProductCard.BonusMechanism,
		})
		if len(info.ProductCard.Images) > 0 {
			product.Image = info.ProductCard.Images[0].URL
		}

		if info.ProductCard.IsBonus {
			inBonus = append(inBonus, product)
		} else {
			notInBonus = append(notInBonus, product)
		}
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
		a.loggerService.Error("Anon auth token request failed", map[string]interface{}{"error": err.Error()})
		return err
	}
	defer resp.Body.Close()

	var authResponse AnonAuthResponse
	if err = json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		a.loggerService.Error("Anon auth token decoding failed", map[string]interface{}{"error": err.Error()})
		return err
	}

	a.loggerService.Info("Got anon auth token", nil)
	a.loggerService.Debug("Anon auth token", map[string]interface{}{"token": authResponse.AccessToken})
	a.authToken = authResponse.AccessToken

	return nil
}

func (a *Appie) makeRequest(product Product) (ProductInfoResponse, error) {
	url := fmt.Sprintf("%s/mobile-services/product/detail/v4/fir/%d", restURL, product.AppieId)
	a.loggerService.Info("Fetching product info", map[string]interface{}{"product": product.FriendlyName})
	a.loggerService.Debug("Product request URL", map[string]interface{}{"product": product.FriendlyName, "url": url})

	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ProductInfoResponse{}, err
	}
	a.addBaseHeaders(true, r)

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		a.loggerService.Error("Product info request failed", map[string]interface{}{"product": product.FriendlyName, "error": err.Error()})
		return ProductInfoResponse{}, err
	}
	defer resp.Body.Close()

	var info ProductInfoResponse
	if err = json.NewDecoder(resp.Body).Decode(&info); err != nil {
		a.loggerService.Error("Failed to decode product info response", map[string]interface{}{"product": product.FriendlyName, "error": err.Error()})
		return ProductInfoResponse{}, err
	}

	a.loggerService.Debug("Product info response", map[string]interface{}{
		"product":         product.FriendlyName,
		"is_bonus":        info.ProductCard.IsBonus,
		"bonus_mechanism": info.ProductCard.BonusMechanism,
	})
	return info, nil
}
