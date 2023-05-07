package model

type (
	Recipes []*Recipe

	Recipe struct {
		Recipe   string        `json:"recipe"`
		Postcode string        `json:"postcode"`
		Delivery string        `json:"delivery"`
		From     *DeliveryTime `json:"-"`
		To       *DeliveryTime `json:"-"`
	}
)
