package appie

import (
	"time"
)

type Product struct {
	ID uint `gorm:"primarykey" json:"id"`

	ApiName      string `json:"api_name"`
	FriendlyName string `json:"friendly_name"`
	RefererUrl   string `json:"referer_url"`
	AppieId      int    `json:"appie_id"`

	DiscountedProducts []DiscountedProducts `json:"-"`
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

type ProductInfoResponse struct {
	Data struct {
		Product struct {
			Id         int    `json:"id"`
			Title      string `json:"title"`
			SmartLabel string `json:"smartLabel"`
			Price      struct {
				Now struct {
					Amount float64 `json:"amount"`
				} `json:"now"`
				Was struct {
					Amount float64
				} `json:"was"`
				UnitInfo struct {
					Price struct {
						Amount float64 `json:"amount"`
					} `json:"price"`
					Description string `json:"description"`
				} `json:"unitInfo"`
				Discount struct {
					SegmentId   int    `json:"segmentId"`
					Description string `json:"description"`
				} `json:"discount"`
			}
		} `json:"product"`
	} `json:"data"`
}
