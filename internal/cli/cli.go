package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/davido912-recipe-count-test-2020/internal/log"
	"github.com/davido912-recipe-count-test-2020/internal/model"
	"github.com/spf13/cobra"
)

const (
	appName = "ivwcli"
	version = "0.0.1"
)

type CobraRunFunc func(cmd *cobra.Command, args []string) error

// Flags
var (
	LogEnabled       bool
	Filepath         string
	Output           *os.File
	MatchRecipeTerms []string
	Postcode         string
	DeliveryFrom     *model.DeliveryTime
	_deliveryFrom    string
	DeliveryTo       *model.DeliveryTime
	_deliveryTo      string
)

// Flag names
const (
	logEnableFlag    = "log"
	filepathFlag     = "file"
	outputFlag       = "output"
	postcodeFlag     = "count-postcode"
	deliveryToFlag   = "to"
	deliveryFromFlag = "from"
)

func NewRootCmd(entrypointFunc CobraRunFunc) *cobra.Command {
	cmd := &cobra.Command{
		Use:   appName,
		Short: "CLI implementation for processing recipe JSON files",
		Long:  "ivwCLI is a CLI tool enabling the processing or JSON data files and producing an aggregate report",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			logEnabled, _ := cmd.Flags().GetBool(logEnableFlag)
			if logEnabled {
				log.InitLogging()
			} else {
				log.SilenceLogging()
			}

			return validateFlags(cmd)
		},
		RunE:         entrypointFunc,
		SilenceUsage: true,
		Example: "ivwcli --file /tmp/file.json --match-recipes 'Speedy Steak Fajitas,Tex-Mex Tilapia' " +
			"-o stdout -p 10120 --from 11AM --to 3PM",
	}

	cmd.Flags().BoolVarP(&LogEnabled, logEnableFlag, "l", false, "Enable logs")
	cmd.Flags().StringVarP(&Filepath, filepathFlag, "f", "", "JSON file to process")
	cmd.Flags().StringP(outputFlag, "o", "stdout", "Output path for result (file/STDOUT)")

	cmd.Flags().StringVarP(&Postcode, postcodeFlag, "p", "10120", "specific postcode to count")
	cmd.Flags().StringVar(&_deliveryFrom, deliveryFromFlag, "10AM", "set delivery start time for postcode count (inclusive)")
	cmd.Flags().StringVar(&_deliveryTo, deliveryToFlag, "3PM", "set delivery end time for postcode count (inclusive)")

	cmd.Flags().StringSliceVarP(
		&MatchRecipeTerms,
		"match-recipes",
		"m",
		[]string{"Potato", "Veggie", "Mushroom"},
		"Match recipe names (comma separated)`",
	)

	cmd.MarkFlagRequired(filepathFlag)

	return cmd
}

// MustCli instantiates the CLI with the given entrypoint function passed to it
func MustCli(rootCmd *cobra.Command) {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

// validateFlags validates values passed to flags
func validateFlags(cmd *cobra.Command) error {
	err := validateOutputFlag(cmd)
	if err != nil {
		return err
	}
	return validateDeliveryFlags(cmd)
}

// validateDeliveryFlags validates that the delivery flags are passed in correct format. additionally, the timespan
// passed must occur in the same 24hour period, for example 3AM to 1PM, NOT 8PM to 2AM
func validateDeliveryFlags(cmd *cobra.Command) (err error) {
	from, _ := cmd.Flags().GetString(deliveryFromFlag)
	to, _ := cmd.Flags().GetString(deliveryToFlag)

	DeliveryFrom, err = model.NewDeliveryTime(from)
	if err != nil {
		return err
	}
	DeliveryTo, err = model.NewDeliveryTime(to)
	if err != nil {
		return err
	}

	if DeliveryTo.Before(DeliveryFrom.Time) {
		return fmt.Errorf("invalid delivery time (%s - %s) delivery times are not date scoped,"+
			" 'to' must occur before 'from' time", from, to)
	}

	return nil
}

// validateOutputFlag validates that output is either stdout or a valid file path
func validateOutputFlag(cmd *cobra.Command) error {
	val, _ := cmd.Flags().GetString(outputFlag)

	if strings.ToLower(val) == "stdout" {
		Output = os.Stdout
	} else {
		f, err := os.OpenFile(val, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		Output = f
	}

	return nil
}
