package aggregate

import (
	"fmt"
	"github.com/davido912-recipe-count-test-2020/internal/model"
)

const (
	PostcodeLenConstraint = 10
	DistinctPostcodesCap  = 1_000_000
)

type postcodeMap map[string]int

type PostcodeAggregator struct {
	postcodeMap
	postcodeTimeCount model.PostcodeTimeCount
	aggrDeliveryFrom  *model.DeliveryTime
	aggrDeliveryTo    *model.DeliveryTime
}

// add adds postcode to map while incrementing its count
func (pm postcodeMap) add(recipe *model.Recipe) error {
	if len(recipe.Postcode) > PostcodeLenConstraint {
		return fmt.Errorf("postcode: %s is longer than limit: %d", recipe.Postcode, PostcodeLenConstraint)
	}
	pm[recipe.Postcode]++
	return nil
}

// aggregate aggregates all the relevant data required from recipes + performs checks
func (pa *PostcodeAggregator) aggregate(recipe *model.Recipe) {

	_ = pa.add(recipe)

	if pa.checkPostcodeEquals(recipe) && pa.checkDeliveryInTimespan(recipe) {
		pa.incrementPostcodeCount()
	}

}

// GetBusiestPostcode return a model.PostcodeCount model with the Postcode that has the most events
func (pa *PostcodeAggregator) GetBusiestPostcode() model.PostcodeCount {
	busiestPostcode := model.PostcodeCount{
		Postcode:      "n/a",
		DeliveryCount: 0,
	}

	for k, v := range pa.postcodeMap {
		if v > busiestPostcode.DeliveryCount {
			busiestPostcode = model.PostcodeCount{
				Postcode:      k,
				DeliveryCount: v,
			}
		}
	}
	return busiestPostcode
}

// GetPostcodeTimeCount return a model.PostcodeTimeCount that contains the count of all deliveries happening in the
// designated postcode during the designated delivery timespan (e.g. 4PM to 8PM)
func (pa *PostcodeAggregator) GetPostcodeTimeCount() model.PostcodeTimeCount {
	return pa.postcodeTimeCount
}

// checkPostcodeEquals checks whether the postcode passed to the aggregator is equal the postcode of a recipe
func (pa *PostcodeAggregator) checkPostcodeEquals(recipe *model.Recipe) bool {
	return pa.postcodeTimeCount.Postcode == recipe.Postcode
}

// checkDeliveryInTimespan check if recipe delivery timespan is in the timespan of the delivery time passed
// to the aggregator
func (pa *PostcodeAggregator) checkDeliveryInTimespan(recipe *model.Recipe) bool {
	return recipe.From.InclusiveBetween(pa.aggrDeliveryFrom, pa.aggrDeliveryTo) &&
		recipe.To.InclusiveBetween(pa.aggrDeliveryFrom, pa.aggrDeliveryTo)
}

func (pa *PostcodeAggregator) incrementPostcodeCount() {
	pa.postcodeTimeCount.DeliveryCount++
}
