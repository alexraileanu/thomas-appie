package appie

import (
	"time"
)

type Product struct {
	ID uint `gorm:"primarykey" json:"-"`

	ApiName      string `json:"api_name"`
	FriendlyName string `json:"friendly_name"`
	RefererUrl   string `json:"referer_url"`
	AppieId      int    `json:"appie_id"`

	Image string `json:"image"`

	DiscountedProducts []DiscountedProducts `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	Discount           DiscountedProducts   `json:"discount" gorm:"-"`
	LastBonus          *time.Time           `json:"last_bonus" gorm:"-"`

	CreatedAt time.Time `gorm:"column:created_at" json:"-"`
}

type DiscountedProducts struct {
	ID uint `gorm:"primarykey" json:"-"`

	ProductID   uint     `json:"-"`
	InBonus     bool     `json:"in_bonus"`
	Description string   `json:"description"`
	Label       string   `json:"label"`
	Product     *Product `json:"-"`

	CreatedAt time.Time `gorm:"column:created_at" json:"-"`
}

type AnonAuthResponse struct {
	AccessToken string `json:"access_token"`
}

type graphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type imageRendition struct {
	URL string `json:"url"`
}

// ProductDetails is the subset of AH's FetchProductDetails GraphQL product we use.
type ProductDetails struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	PriceV2 struct {
		Discount *struct {
			Description string `json:"description"`
		} `json:"discount"`
	} `json:"priceV2"`
	ImagePack []struct {
		Large  *imageRendition `json:"large"`
		Medium *imageRendition `json:"medium"`
		Small  *imageRendition `json:"small"`
	} `json:"imagePack"`
}

// IsBonus reports whether the product currently has a discount.
func (p *ProductDetails) IsBonus() bool {
	return p.PriceV2.Discount != nil
}

// BonusMechanism returns the discount description, or "" when not in bonus.
func (p *ProductDetails) BonusMechanism() string {
	if p.PriceV2.Discount == nil {
		return ""
	}
	return p.PriceV2.Discount.Description
}

// ImageURL returns the medium image URL, falling back to large, then small.
func (p *ProductDetails) ImageURL() string {
	if len(p.ImagePack) == 0 {
		return ""
	}
	pack := p.ImagePack[0]
	for _, r := range []*imageRendition{pack.Medium, pack.Large, pack.Small} {
		if r != nil && r.URL != "" {
			return r.URL
		}
	}
	return ""
}

type ProductDetailsResponse struct {
	Data struct {
		Product *ProductDetails `json:"product"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}
