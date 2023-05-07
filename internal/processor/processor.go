package processor

import (
	"fmt"
	"github.com/davido912-recipe-count-test-2020/internal/aggregate"
	"github.com/davido912-recipe-count-test-2020/internal/model"
	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"
	"io"
	"regexp"
	"sync"
)

type Processor struct {
	*aggregate.Aggregator
	deliveryRegex *regexp.Regexp
	chunkSize     int
	dlq           chan *model.Recipe
}

func NewProcessor(chunkSize int, aggrinput *aggregate.AggregatorInput, dlq chan *model.Recipe) *Processor {
	rgx, err := regexp.Compile("(?:1[0-2]|[1-9])[AP]M")
	if err != nil {

		panic(err)
	}
	return &Processor{
		Aggregator:    aggregate.NewAggregator(aggrinput),
		deliveryRegex: rgx,
		chunkSize:     chunkSize,
		dlq:           dlq,
	}
}

// Process main entrypoint of this component - reads the data, breaks it into chunks for faster processing.
// if event is invalid it is discarded or forwarded to dlq channel (if present). Aggregates are finally calculated and
// end report model is generated
func (p *Processor) Process(data io.Reader) (*model.ReportModel, error) {
	recipes, err := p.unmarshalRecipeData(data)
	if err != nil {
		return nil, fmt.Errorf("failed parsing JSON input file: %w", err)
	}

	chunks := toChunks(recipes, p.chunkSize)
	log.Debug().Msgf("chunk size of %d generated %d chunks", p.chunkSize, len(chunks))

	var processorsWg sync.WaitGroup
	processorsWg.Add(len(chunks))

	// recipeChan is used to allow concurrent aggregating while processing is still taking place
	// buffered channel is used to not block
	recipeChan := make(chan *model.Recipe, len(recipes))

	for _, chunk := range chunks {

		go func(recipes model.Recipes) {
			defer processorsWg.Done()
			for _, recipe := range recipes {
				err := p.processRecipe(recipe)
				if err != nil {
					log.Error().Err(err).Msgf("failed processing recipe: %T", recipe)
					if p.dlq != nil {
						p.dlq <- recipe
					}

				} else {
					recipeChan <- recipe
				}
			}
		}(chunk)

	}

	var aggregatorsWg sync.WaitGroup
	aggregatorsWg.Add(1)

	go func() {
		defer aggregatorsWg.Done()
		p.Aggregate(recipeChan)
	}()

	// wait for processors to finish transmitting all events to recipe channel
	processorsWg.Wait()
	close(recipeChan)
	aggregatorsWg.Wait()

	return p.generateReport(), nil
}

// generateReport outputs the final model used for the reporting
func (p *Processor) generateReport() *model.ReportModel {
	reportModel := model.NewReportModel()
	reportModel.SetUniqueRecipeCount(p.GetUniqueRecipeCount())
	reportModel.SetCountPerRecipe(p.GetRecipeCountsModel())
	reportModel.SetMatchByName(p.GetRecipeMatches())
	reportModel.SetCountPerPostcodeAndTime(p.GetPostcodeTimeCount())
	reportModel.SetBusiestPostcode(p.GetBusiestPostcode())
	return reportModel
}

// processRecipe validates field + parses event
func (p *Processor) processRecipe(recipe *model.Recipe) error {
	if err := p.validateRequiredFields(recipe); err != nil {
		return err
	}

	err := p.parseDelivery(recipe)
	if err != nil {
		return err
	}

	return nil
}

// unmarshalRecipeData read data from a buffer/file and deserialize into []model.Recipe
func (p *Processor) unmarshalRecipeData(data io.Reader) (model.Recipes, error) {
	bs, _ := io.ReadAll(data)

	var recipes model.Recipes

	err := json.Unmarshal(bs, &recipes)
	if err != nil {
		return nil, err
	}

	return recipes, nil
}

// validateRequiredFields ensures all the fields are present in the JSON events
func (p *Processor) validateRequiredFields(recipe *model.Recipe) error {
	if recipe.Recipe == "" || recipe.Delivery == "" || recipe.Postcode == "" {
		return fmt.Errorf("one of required fields [postcode, delivery, recipe] is missing or blank")
	}
	return nil
}

// parseDelivery parses the delivery field in the JSON events into a deserialized object
func (p *Processor) parseDelivery(recipe *model.Recipe) error {
	found := p.deliveryRegex.FindAll([]byte(recipe.Delivery), -1)

	if len(found) < 2 {
		return fmt.Errorf("invalid delivery time: %s", recipe.Delivery)
	}

	from, to := string(found[0]), string(found[1])

	parsedFrom, err := model.NewDeliveryTime(from)
	if err != nil {
		return err
	}
	parsedTo, err := model.NewDeliveryTime(to)
	if err != nil {
		return err
	}
	recipe.From, recipe.To = parsedFrom, parsedTo
	return nil
}

// toChunks helper function to chunk big arrays resulted from JSON unmarshalling to allow for concurrency and faster
// processing. Smallest chunkSize possible is 1
func toChunks[T any](slice []T, chunkSize int) [][]T {
	if chunkSize < 1 {
		chunkSize = 1
	}

	var chunks [][]T
	for {
		if len(slice) == 0 {
			break
		}

		if len(slice) < chunkSize {
			chunkSize = len(slice)
		}

		chunks = append(chunks, slice[:chunkSize])
		slice = slice[chunkSize:]
	}

	return chunks
}
