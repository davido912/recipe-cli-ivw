package aggregate

import (
	"github.com/davido912-recipe-count-test-2020/internal/model"
	"github.com/davido912-recipe-count-test-2020/internal/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func mockPostcodeAggr(aggrInput AggregatorInput) *PostcodeAggregator {
	return &PostcodeAggregator{
		postcodeMap: make(postcodeMap, DistinctPostcodesCap),
		postcodeTimeCount: model.PostcodeTimeCount{
			Postcode:      aggrInput.Postcode,
			From:          aggrInput.DeliveryFrom.Raw(),
			To:            aggrInput.DeliveryTo.Raw(),
			DeliveryCount: 0,
		},
		aggrDeliveryTo:   aggrInput.DeliveryTo,
		aggrDeliveryFrom: aggrInput.DeliveryFrom,
	}
}

func TestPostcodeAggregator_Aggregate(t *testing.T) {
	aggrInput := AggregatorInput{
		Postcode:     "10245",
		DeliveryFrom: testutils.MockDeliveryTime("12PM"),
		DeliveryTo:   testutils.MockDeliveryTime("5PM"),
	}
	aggr := mockPostcodeAggr(aggrInput)
	recipes := testutils.MockRecipes()
	for _, recipe := range recipes {
		aggr.aggregate(recipe)
	}

	wantMap := postcodeMap{
		"10245": 3,
		"10342": 1,
		"10311": 2,
	}
	assert.Equal(t, wantMap, aggr.postcodeMap)
	assert.Equal(t, 2, aggr.postcodeTimeCount.DeliveryCount)
}

func TestPostcodeAggregator_GetBusiestPostcode(t *testing.T) {
	aggrInput := AggregatorInput{
		Postcode:     "10245",
		DeliveryFrom: testutils.MockDeliveryTime("12PM"),
		DeliveryTo:   testutils.MockDeliveryTime("5PM"),
	}
	aggr := mockPostcodeAggr(aggrInput)
	recipes := testutils.MockRecipes()
	for _, recipe := range recipes {
		aggr.aggregate(recipe)
	}

	want := model.PostcodeCount{
		Postcode:      "10245",
		DeliveryCount: 3,
	}
	assert.Equal(t, want, aggr.GetBusiestPostcode())
}

func TestPostcodeAggregator_GetPostcodeTimeCount(t *testing.T) {
	aggrInput := AggregatorInput{
		Postcode:     "10245",
		DeliveryFrom: testutils.MockDeliveryTime("12PM"),
		DeliveryTo:   testutils.MockDeliveryTime("5PM"),
	}
	aggr := mockPostcodeAggr(aggrInput)
	recipes := testutils.MockRecipes()
	for _, recipe := range recipes {
		aggr.aggregate(recipe)
	}

	want := model.PostcodeTimeCount{
		Postcode:      "10245",
		From:          aggrInput.DeliveryFrom.Raw(),
		To:            aggrInput.DeliveryTo.Raw(),
		DeliveryCount: 2,
	}
	assert.Equal(t, want, aggr.GetPostcodeTimeCount())
}

func TestPostcodeMap_Add(t *testing.T) {

	tcs := []struct {
		name      string
		postcodes model.Recipes
		want      postcodeMap
		wantErr   bool
	}{
		{
			name: "happy path",
			postcodes: []*model.Recipe{
				{Postcode: "1234"},
				{Postcode: "1234"},
				{Postcode: "1555"},
			},
			want:    postcodeMap{"1234": 2, "1555": 1},
			wantErr: false,
		},
		{
			name: "invalid postcode",
			postcodes: []*model.Recipe{
				{Postcode: "1234567891011"},
			},
			want:    postcodeMap{},
			wantErr: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			pm := postcodeMap{}
			for _, pc := range tc.postcodes {
				err := pm.add(pc)
				if tc.wantErr {
					assert.NotNil(t, err)
				}
			}
			assert.Equal(t, tc.want, pm)

		})
	}
}

func TestCheckDeliveryInTimespan(t *testing.T) {
	aggrInput := AggregatorInput{
		Postcode:     "10245",
		DeliveryFrom: testutils.MockDeliveryTime("10AM"),
		DeliveryTo:   testutils.MockDeliveryTime("3PM"),
	}
	aggr := mockPostcodeAggr(aggrInput)

	tcs := []struct {
		name       string
		aggr       *PostcodeAggregator
		recipeFrom *model.DeliveryTime
		recipeTo   *model.DeliveryTime
		want       bool
	}{
		{
			name:       "delivery in timespan",
			aggr:       aggr,
			recipeFrom: testutils.MockDeliveryTime("11AM"),
			recipeTo:   testutils.MockDeliveryTime("1PM"),
			want:       true,
		},
		{
			name:       "delivery not in timespan",
			aggr:       aggr,
			recipeFrom: testutils.MockDeliveryTime("9AM"),
			recipeTo:   testutils.MockDeliveryTime("1PM"),
			want:       false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			recipe := &model.Recipe{
				From: tc.recipeFrom,
				To:   tc.recipeTo,
			}
			got := aggr.checkDeliveryInTimespan(recipe)
			assert.Equal(t, tc.want, got)

		})
	}
}

func TestCheckPostcodeEquals(t *testing.T) {
	postcode := "15231"

	aggrInput := AggregatorInput{
		Postcode:     postcode,
		DeliveryFrom: testutils.MockDeliveryTime("10AM"),
		DeliveryTo:   testutils.MockDeliveryTime("3PM"),
	}
	aggr := mockPostcodeAggr(aggrInput)

	tcs := []struct {
		name          string
		aggr          *PostcodeAggregator
		checkPostcode string
		want          bool
	}{
		{
			name:          "postcode matches",
			aggr:          aggr,
			checkPostcode: postcode,
			want:          true,
		},
		{
			name:          "postcode does not match",
			aggr:          aggr,
			checkPostcode: "22",
			want:          false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			recipe := &model.Recipe{Postcode: tc.checkPostcode}
			got := tc.aggr.checkPostcodeEquals(recipe)
			assert.Equal(t, tc.want, got)

		})
	}
}
