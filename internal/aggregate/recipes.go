package aggregate

import (
	"fmt"
	"github.com/davido912-recipe-count-test-2020/internal/model"
	"sort"
	"strings"
)

const (
	RecipeNameLenConstraint = 100
	DistinctRecipeCap       = 2000
)

type (

	// recipeMap represents recipe names and their counts
	recipeMap map[string]int

	// recipeMatcher case-sensitive term matching recipe names to terms passed by the user
	recipeMatcher struct {
		matches []string
		terms   []string
	}

	RecipeAggregator struct {
		recipeMap
		recipeMatcher
		sortedRecipeNames []string
	}
)

// aggregate aggregates primary data required for this component.
func (ra *RecipeAggregator) aggregate(recipe *model.Recipe) {
	_ = ra.add(recipe)
}

// postAggregate is used for additional aggregations that are done after the initial data has been aggregated
func (ra *RecipeAggregator) postAggregate() {
	ra.sortedRecipeNames = make([]string, 0, len(ra.recipeMap))

	for key := range ra.recipeMap {
		ra.sortedRecipeNames = append(ra.sortedRecipeNames, key)
	}

	sort.Strings(ra.sortedRecipeNames)
}

func (ra *RecipeAggregator) GetUniqueRecipeCount() int {
	return len(ra.recipeMap)
}

// GetRecipeCountsModel generates final model output used in the end report that details the counts of each recipe
// sorted by recipe name in ascending order
func (ra *RecipeAggregator) GetRecipeCountsModel() model.RecipeCounts {
	recipeCounts := make(model.RecipeCounts, 0, len(ra.recipeMap))

	for _, key := range ra.sortedRecipeNames {
		recipeCounts = append(recipeCounts, model.RecipeCount{
			Recipe:      key,
			RecipeCount: ra.recipeMap[key],
		})

	}

	return recipeCounts
}

func (ra *RecipeAggregator) GetRecipeMatches() model.RecipeMatches {
	for _, recipe := range ra.sortedRecipeNames {
		ra.match(recipe)
	}
	return ra.recipeMatcher.matches
}

func (r *recipeMatcher) append(recipeName string) {
	r.matches = append(r.matches, recipeName)
}

func (r *recipeMatcher) match(recipeName string) {
	for _, word := range r.terms {
		if strings.Contains(recipeName, word) {
			r.append(recipeName)
		}
	}
}

func (rm recipeMap) add(recipe *model.Recipe) error {
	if len(recipe.Recipe) > RecipeNameLenConstraint {
		return fmt.Errorf("recipe name: %s is longer than limit: %d", recipe.Recipe, RecipeNameLenConstraint)
	}
	rm[recipe.Recipe]++
	return nil
}
