package main

type GQLQuery struct {
	Product struct {
		Id         int
		Title      string
		SmartLabel string
		Price      struct {
			Now struct {
				Amount float64
			}
			Was struct {
				Amount float64
			}
			UnitInfo struct {
				Price struct {
					Amount float64
				}
				Description string
			}
			Discount struct {
				SegmentId   int
				Description string
			}
		}
	} `graphql:"product(id: $id, date: $date)"`
}

type Query struct {
	Query GQLQuery `graphql:"product($id: Int!, $date: String)"`
}
