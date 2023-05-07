package testutils

import (
	"bytes"
	"github.com/davido912-recipe-count-test-2020/internal/model"
	"io"
)

func MockDeliveryTime(deliveryTime string) *model.DeliveryTime {
	t, err := model.NewDeliveryTime(deliveryTime)
	if err != nil {
		panic(err)
	}
	return t
}

// MockData represents the same data found in MockRecipes
func MockData() io.Reader {
	data := `
  [{
    "postcode": "10245",
    "recipe": "Apple",
    "delivery": "Wednesday 1PM - 5PM"
  },
  {
    "postcode": "10245",
    "recipe": "Steak",
    "delivery": "Thursday 10AM - 2PM"
  },
  {
    "postcode": "10245",
    "recipe": "Salt",
    "delivery": "Thursday 12PM - 2PM"
  },
  {
    "postcode": "10342",
    "recipe": "Pear",
    "delivery": "Thursday 8PM - 11PM"
  },
  {
    "postcode": "10311",
    "recipe": "Honey",
    "delivery": "Thursday 3PM - 4PM"
  },
  {
    "postcode": "10311",
    "recipe": "Honey",
    "delivery": "Thursday 3PM - 4PM"
  }]
`
	return bytes.NewBufferString(data)
}

// MockRecipes represents the same data found in MockData
func MockRecipes() model.Recipes {
	return model.Recipes{
		{Recipe: "Apple", Postcode: "10245", Delivery: "Wednesday 1PM - 5PM", From: MockDeliveryTime("1PM"), To: MockDeliveryTime("5PM")},
		{Recipe: "Steak", Postcode: "10245", Delivery: "Thursday 10AM - 2PM", From: MockDeliveryTime("10AM"), To: MockDeliveryTime("2PM")},
		{Recipe: "Salt", Postcode: "10245", Delivery: "Thursday 12PM - 2PM", From: MockDeliveryTime("12PM"), To: MockDeliveryTime("2PM")},
		{Recipe: "Pear", Postcode: "10342", Delivery: "Thursday 8PM - 11PM", From: MockDeliveryTime("8PM"), To: MockDeliveryTime("11PM")},
		{Recipe: "Honey", Postcode: "10311", Delivery: "Thursday 3PM - 4PM", From: MockDeliveryTime("3PM"), To: MockDeliveryTime("4PM")},
		{Recipe: "Honey", Postcode: "10311", Delivery: "Thursday 3PM - 4PM", From: MockDeliveryTime("3PM"), To: MockDeliveryTime("4PM")},
	}
}
