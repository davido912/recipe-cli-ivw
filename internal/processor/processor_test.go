package processor

import (
	"bytes"
	"github.com/davido912-recipe-count-test-2020/internal/aggregate"
	"github.com/davido912-recipe-count-test-2020/internal/model"
	"github.com/davido912-recipe-count-test-2020/internal/testutils"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestProcessor_Process(t *testing.T) {
	tcs := []struct {
		name    string
		data    io.Reader
		want    *model.ReportModel
		wantErr bool
	}{
		{
			name: "invalid event not included",
			data: bytes.NewBufferString(`[{"postcode": "10311","recipe": "Honey","delivery": "Thursday 3PM - 4PM"}
											,{"recipe": "Steak","delivery": "Thursday 3PM - 4PM"}]`),
			want: &model.ReportModel{
				UniqueRecipeCount: 1,
				CountPerRecipe: model.RecipeCounts{
					{Recipe: "Honey", RecipeCount: 1},
				},
				BusiestPostcode: model.PostcodeCount{
					Postcode:      "10311",
					DeliveryCount: 1,
				},
				CountPerPostcodeAndTime: model.PostcodeTimeCount{
					Postcode:      "10245",
					From:          "10AM",
					To:            "3PM",
					DeliveryCount: 0,
				},
				MatchByName: model.RecipeMatches{},
			},
			wantErr: false,
		},
		{
			name: "happy path",
			data: testutils.MockData(),
			want: &model.ReportModel{
				UniqueRecipeCount: 5,
				CountPerRecipe: model.RecipeCounts{
					{Recipe: "Apple", RecipeCount: 1},
					{Recipe: "Honey", RecipeCount: 2},
					{Recipe: "Pear", RecipeCount: 1},
					{Recipe: "Salt", RecipeCount: 1},
					{Recipe: "Steak", RecipeCount: 1},
				},
				BusiestPostcode: model.PostcodeCount{
					Postcode:      "10245",
					DeliveryCount: 3,
				},
				CountPerPostcodeAndTime: model.PostcodeTimeCount{
					Postcode:      "10245",
					From:          "10AM",
					To:            "3PM",
					DeliveryCount: 2,
				},
				MatchByName: model.RecipeMatches{"Pear", "Steak"},
			},
			wantErr: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			aggrInput := &aggregate.AggregatorInput{
				Postcode:     "10245",
				DeliveryFrom: testutils.MockDeliveryTime("10AM"),
				DeliveryTo:   testutils.MockDeliveryTime("3PM"),
				Terms:        []string{"ea"},
			}
			p := NewProcessor(1, aggrInput, nil)

			got, err := p.Process(tc.data)

			if tc.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tc.want, got)

		})
	}
}

func TestProcessor_unmarshalRecipeData(t *testing.T) {
	tcs := []struct {
		name    string
		data    io.Reader
		want    model.Recipes
		wantErr bool
	}{
		{
			name:    "happy path",
			data:    bytes.NewBufferString(`[{"postcode": "10311","recipe": "Honey","delivery": "Thursday 3PM - 4PM"}]`),
			want:    model.Recipes{&model.Recipe{Postcode: "10311", Recipe: "Honey", Delivery: "Thursday 3PM - 4PM"}},
			wantErr: false,
		},
		{
			name:    "invalid json passed",
			data:    bytes.NewBufferString(`fff`),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			p := Processor{}
			gotRecipes, err := p.unmarshalRecipeData(tc.data)
			if tc.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tc.want, gotRecipes)

		})
	}
}

func TestProcessor_validateRequiredFields(t *testing.T) {
	tcs := []struct {
		name    string
		recipe  *model.Recipe
		wantErr bool
	}{
		{
			name:    "valid recipe",
			recipe:  &model.Recipe{Postcode: "10311", Recipe: "Honey", Delivery: "Thursday 3PM - 4PM"},
			wantErr: false,
		},
		{
			name:    "missing recipe name",
			recipe:  &model.Recipe{Postcode: "10311", Delivery: "Thursday 3PM - 4PM"},
			wantErr: true,
		},
		{
			name:    "missing postcode",
			recipe:  &model.Recipe{Recipe: "Honey", Delivery: "Thursday 3PM - 4PM"},
			wantErr: true,
		},
		{
			name:    "missing delivery",
			recipe:  &model.Recipe{Postcode: "10311", Recipe: "Honey"},
			wantErr: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			p := Processor{}
			err := p.validateRequiredFields(tc.recipe)
			if tc.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

		})
	}
}

func TestProcessor_parseDelivery(t *testing.T) {
	tcs := []struct {
		name     string
		recipe   *model.Recipe
		wantErr  bool
		wantFrom *model.DeliveryTime
		wantTo   *model.DeliveryTime
	}{
		{
			name:     "valid delivery",
			recipe:   &model.Recipe{Postcode: "10311", Recipe: "Honey", Delivery: "Thursday 3PM - 4PM"},
			wantErr:  false,
			wantFrom: testutils.MockDeliveryTime("3PM"),
			wantTo:   testutils.MockDeliveryTime("4PM"),
		},
		{
			name:     "invalid delivery",
			recipe:   &model.Recipe{Postcode: "10311", Delivery: "Thursday 4PM"},
			wantErr:  true,
			wantFrom: nil,
			wantTo:   nil,
		},
		{
			name:     "missing delivery",
			recipe:   &model.Recipe{Recipe: "Honey"},
			wantErr:  true,
			wantFrom: nil,
			wantTo:   nil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			aggrInput := &aggregate.AggregatorInput{
				Postcode:     "10311",
				DeliveryFrom: testutils.MockDeliveryTime("8PM"),
				DeliveryTo:   testutils.MockDeliveryTime("11PM"),
				Terms:        []string{},
			}
			p := NewProcessor(0, aggrInput, nil)
			err := p.parseDelivery(tc.recipe)
			if tc.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tc.wantFrom, tc.recipe.From)
			assert.Equal(t, tc.wantTo, tc.recipe.To)
		})
	}
}

func Test_toChunks(t *testing.T) {
	tcs := []struct {
		name       string
		slice      []string
		chunksize  int
		wantChunks int
	}{
		{
			name:       "test chunks when even chunk size",
			slice:      make([]string, 100),
			chunksize:  20,
			wantChunks: 5,
		},
		{
			name:       "test chunks when even uneven chunk size",
			slice:      make([]string, 100),
			chunksize:  30,
			wantChunks: 4,
		},
		{
			name:       "test chunks when chunk size bigger than slice size",
			slice:      make([]string, 100),
			chunksize:  200,
			wantChunks: 1,
		},
		{
			name:       "zero chunk size - defaults to 1",
			slice:      make([]string, 100),
			chunksize:  0,
			wantChunks: 100,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			chunks := toChunks(tc.slice, tc.chunksize)
			assert.Equal(t, tc.wantChunks, len(chunks))

		})
	}
}
