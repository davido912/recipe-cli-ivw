package aggregate

import (
	"github.com/davido912-recipe-count-test-2020/internal/model"
	"github.com/davido912-recipe-count-test-2020/internal/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRecipeAggregator_Aggregate(t *testing.T) {
	recipes := testutils.MockRecipes()

	aggr := RecipeAggregator{
		recipeMap: make(recipeMap),
		recipeMatcher: recipeMatcher{
			terms: []string{"ea"},
		},
	}
	for _, recipe := range recipes {
		aggr.aggregate(recipe)
	}

	aggr.postAggregate()

	wantMatches := model.RecipeMatches{"Pear", "Steak"}
	wantMap := recipeMap{
		"Apple": 1,
		"Steak": 1,
		"Salt":  1,
		"Pear":  1,
		"Honey": 2,
	}

	assert.Equal(t, wantMatches, aggr.GetRecipeMatches())
	assert.Equal(t, wantMap, aggr.recipeMap)
}

func TestRecipeAggregator_GetUniqueRecipeCount(t *testing.T) {
	recipes := testutils.MockRecipes()
	aggr := RecipeAggregator{
		recipeMap: make(recipeMap),
		recipeMatcher: recipeMatcher{
			terms: []string{},
		},
	}
	for _, recipe := range recipes {
		aggr.aggregate(recipe)
	}
	aggr.postAggregate()

	got := aggr.GetUniqueRecipeCount()

	assert.Equal(t, 5, got)
}

func TestRecipeAggregator_GetRecipeCountsModel(t *testing.T) {
	recipes := testutils.MockRecipes()
	aggr := RecipeAggregator{
		recipeMap: make(recipeMap),
		recipeMatcher: recipeMatcher{
			terms: []string{},
		},
	}
	for _, recipe := range recipes {
		aggr.aggregate(recipe)
	}
	aggr.postAggregate()
	got := aggr.GetRecipeCountsModel()

	want := model.RecipeCounts{
		{Recipe: "Apple", RecipeCount: 1},
		{Recipe: "Honey", RecipeCount: 2},
		{Recipe: "Pear", RecipeCount: 1},
		{Recipe: "Salt", RecipeCount: 1},
		{Recipe: "Steak", RecipeCount: 1},
	}

	assert.Equal(t, want, got)
}

func TestRecipeMatcher_match(t *testing.T) {
	rm := recipeMatcher{
		matches: []string{},
		terms:   []string{"Apple"},
	}

	rm.match("Honey Apple Pie")

	want := []string{"Honey Apple Pie"}
	assert.Equal(t, want, rm.matches)
}
