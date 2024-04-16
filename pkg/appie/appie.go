package appie

import (
	"embed"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

const URL = "https://www.ah.nl/gql"

//go:embed queryFormat.json
var queryFormatFile embed.FS

func PerformProductsCheck(productsToWatch []*ProductToCheck) ([]ProductToCheck, []ProductToCheck, error) {
	var productsInBonus []ProductToCheck
	var productsNotInBonus []ProductToCheck

	for _, product := range productsToWatch {
		err := makeGqlRequest(product)
		if err != nil {
			return nil, nil, err
		}
		hasDiscount := product.BonusData.Data.Product.Price.Discount.SegmentId != 0
		product.HasDiscount = hasDiscount
		if hasDiscount {
			productsInBonus = append(productsInBonus, ProductToCheck{
				FriendlyName: product.FriendlyName,
				ApiName:      product.ApiName,
				BonusData:    product.BonusData,
			})
		} else {
			productsNotInBonus = append(productsNotInBonus, ProductToCheck{
				FriendlyName: product.FriendlyName,
				ApiName:      product.ApiName,
			})
		}
	}

	return productsInBonus, productsNotInBonus, nil
}

func makeGqlRequest(product *ProductToCheck) error {
	gqlQuery, err := readQueryFile()
	if err != nil {
		return err
	}

	preparedRequest := fmt.Sprintf(gqlQuery, product.ID, time.Now().Format("2006-01-02"))

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
		return err
	}

	return json.Unmarshal(resp.Body(), &product.BonusData)
}

func readQueryFile() (string, error) {
	data, err := queryFormatFile.ReadFile("queryFormat.json")
	if err != nil {
		return "", err
	}

	return string(data), nil
}
