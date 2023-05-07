package aggregate

import (
	"github.com/davido912-recipe-count-test-2020/internal/model"
	"github.com/davido912-recipe-count-test-2020/internal/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAggregator_listen(t *testing.T) {
	recipes := testutils.MockRecipes()
	recipeChan := make(chan *model.Recipe, len(recipes))

	for _, recipe := range recipes {
		recipeChan <- recipe
	}
	close(recipeChan)

	var processed int
	aggr := Aggregator{}

	aggr.listen(recipeChan, func(_ *model.Recipe) {
		processed++
	})

	assert.Equal(t, len(recipes), processed)
}

// all sub-functionalities inside this function are already tested directly in other tests
// this test just tests that it runs without panicing
func TestAggregator_Aggregate(t *testing.T) {
	recipes := testutils.MockRecipes()
	recipeChan := make(chan *model.Recipe, len(recipes))

	for _, recipe := range recipes {
		recipeChan <- recipe
	}
	close(recipeChan)

	aggrInput := &AggregatorInput{
		Postcode:     "10245",
		DeliveryFrom: testutils.MockDeliveryTime("10AM"),
		DeliveryTo:   testutils.MockDeliveryTime("3PM"),
		Terms:        []string{"ea"},
	}
	aggr := NewAggregator(aggrInput)

	aggr.Aggregate(recipeChan)
}
