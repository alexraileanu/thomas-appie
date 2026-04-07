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

type ProductInfoResponse struct {
	ProductId   uint `json:"productId"`
	ProductCard struct {
		Title  string `json:"title"`
		Images []struct {
			URL string `json:"url"`
		} `json:"images"`
		IsBonus          bool    `json:"isBonus"`
		IsBonusPrice     bool    `json:"isBonusPrice"`
		BonusMechanism   string  `json:"bonusMechanism"`
		PriceBeforeBonus float32 `json:"priceBeforeBonus"`
	} `json:"productCard"`
}
