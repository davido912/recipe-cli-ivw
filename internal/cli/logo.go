package cli

import "github.com/common-nighthawk/go-figure"

const (
	cliLogoPhrase = "Fresh-CLI"
)

var (
	CliLogo = figure.NewColorFigure(cliLogoPhrase,
		"epic",
		"purple",
		true)
)
