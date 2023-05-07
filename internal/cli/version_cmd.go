package cli

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewVersionCmd() *cobra.Command {

	return &cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, args []string) {
			CliLogo.Print()
			fmt.Println()
			printSection("Version", version)
		},
	}
}

func printSection(section, value string) {
	colorize := color.New(color.FgMagenta).SprintfFunc()
	fmt.Printf("%-20s %s\n", colorize(section+":"), value)
}
