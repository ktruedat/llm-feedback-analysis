package log

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

func setupDevLogger() *zerolog.Logger {
	return setupDevLoggerWithSkip(callersToSkip)
}

func setupDevLoggerWithSkip(skip int) *zerolog.Logger {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	logger := zerolog.New(os.Stdout).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		CallerWithSkipFrameCount(skip).
		Logger()

	return &logger
}

func setupProdLogger() *zerolog.Logger {
	logger := zerolog.New(os.Stdout).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Logger()

	return &logger
}
