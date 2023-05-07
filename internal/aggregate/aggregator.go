package aggregate

import (
	"github.com/davido912-recipe-count-test-2020/internal/model"
	"github.com/rs/zerolog/log"
)

type AggregatorInput struct {

	// Postcode, DeliveryFrom, DeliveryTo are used for functional requirement 4
	Postcode                 string
	DeliveryFrom, DeliveryTo *model.DeliveryTime

	// Terms used for functional requirement 5 - matching recipes
	Terms []string
}

type Aggregator struct {
	*PostcodeAggregator
	*RecipeAggregator
}

// NewAggregator returns an instance that calculates postcode and recipe metrics. The parameters passed to
// the aggregator are used to collect additional metrics
func NewAggregator(aggrInput *AggregatorInput) *Aggregator {
	recipeAggregator := &RecipeAggregator{
		recipeMap: make(recipeMap, DistinctRecipeCap),
		recipeMatcher: recipeMatcher{
			matches: make([]string, 0, DistinctRecipeCap),
			terms:   aggrInput.Terms,
		},
	}
	postcodeAggregator := &PostcodeAggregator{
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

	return &Aggregator{
		PostcodeAggregator: postcodeAggregator,
		RecipeAggregator:   recipeAggregator,
	}
}

func (a *Aggregator) Aggregate(recipeChan chan *model.Recipe) {

	// consume all events from channel
	a.listen(recipeChan, func(recipe *model.Recipe) {
		a.PostcodeAggregator.aggregate(recipe)
		a.RecipeAggregator.aggregate(recipe)
	})

	a.RecipeAggregator.postAggregate()
}

func (a *Aggregator) listen(recipeChan chan *model.Recipe, processFunc func(*model.Recipe)) {
	for {
		select {
		case recipe, ok := <-recipeChan:
			if ok == false {
				log.Debug().Msg("recipe channel is empty and closed. done processing.")
				return
			}
			processFunc(recipe)
		}
	}
}
