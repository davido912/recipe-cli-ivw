package log

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func InitLogging() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

func SilenceLogging() {
	zerolog.SetGlobalLevel(zerolog.Disabled)
}
