package shared

import (
	"os"

	"github.com/rs/zerolog"
)

type CtxKey string

const RadioKey CtxKey = "radio"

var Logger = zerolog.New(os.Stdout).With().Timestamp().Logger().Output(zerolog.ConsoleWriter{Out: os.Stdout})
var CWD = Must(os.Getwd())

func Must[T any](t T, err error) T {
	if err != nil {
		Logger.Fatal().Err(err).Msg("Must failed")
	}

	return t
}
