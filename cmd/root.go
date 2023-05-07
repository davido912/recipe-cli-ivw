package cmd

import (
	"github.com/davido912-recipe-count-test-2020/internal/aggregate"
	"github.com/davido912-recipe-count-test-2020/internal/cli"
	"github.com/davido912-recipe-count-test-2020/internal/processor"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = cli.NewRootCmd(func(cmd *cobra.Command, args []string) error {
	log.Debug().Msgf("opening file in path: %s", cli.Filepath)
	dataFile, err := os.Open(cli.Filepath)
	if err != nil {
		return err
	}
	defer func() { _ = dataFile.Close() }()

	aggrInput := &aggregate.AggregatorInput{
		Postcode:     cli.Postcode,
		DeliveryFrom: cli.DeliveryFrom,
		DeliveryTo:   cli.DeliveryTo,
		Terms:        cli.MatchRecipeTerms,
	}

	proc := processor.NewProcessor(2024, aggrInput, nil)
	report, err := proc.Process(dataFile)
	if err != nil {
		return err
	}

	err = report.Dumps(cli.Output)
	if err != nil {
		return err
	}
	defer func() { _ = cli.Output.Close() }()

	return nil
})

func Run() {
	rootCmd.AddCommand(cli.NewVersionCmd())
	cli.MustCli(rootCmd)
}
