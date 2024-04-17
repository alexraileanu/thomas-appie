package appie

type ProductToCheck struct {
	ID           int    `json:"-"`
	ApiName      string `json:"apiName"`
	FriendlyName string `json:"friendlyName"`
	RefererUrl   string `json:"-"`
	InBonus      bool   `json:"inBonus"`

	BonusData ProductInfoResponse `json:"-"`
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
