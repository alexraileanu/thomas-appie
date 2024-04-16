package appie

type ProductToCheck struct {
	ID           int
	ApiName      string
	FriendlyName string
	RefererUrl   string
	HasDiscount  bool

	BonusData ProductInfoResponse
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
