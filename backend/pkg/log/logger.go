package log

import (
	"github.com/rs/zerolog"
)

const (
	productionEnv = "prod"
	callersToSkip = 3
)

type logger struct {
	zlog *zerolog.Logger
}

type loggerOptions struct {
	callersToSkip int
}

// Option configures logger creation.
type Option func(*loggerOptions)

// WithCallersToSkip sets the number of caller frames to skip.
// This is useful for wrapped loggers (like tracelog) that add an extra layer of indirection.
func WithCallersToSkip(skip int) Option {
	return func(opts *loggerOptions) {
		opts.callersToSkip = skip
	}
}

func NewLogger(environment string, opts ...Option) Logger {
	options := &loggerOptions{
		callersToSkip: callersToSkip,
	}

	for _, opt := range opts {
		opt(options)
	}

	var zlog *zerolog.Logger
	if environment == productionEnv {
		zlog = setupProdLogger()
	} else {
		zlog = setupDevLoggerWithSkip(options.callersToSkip)
	}

	return &logger{
		zlog: zlog,
	}
}

func (*logger) handleArgs(event *zerolog.Event, args ...any) {
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			key, okKey := args[i].(string)
			if okKey {
				event.Interface(key, args[i+1])
			} else {
				event = event.Interface("invalid_key_type", args[i])
			}
		} else {
			event = event.Interface("dangling_arg", args[i])
		}
	}
}

func (l *logger) Error(msg string, err error, args ...any) {
	event := l.zlog.Error().Stack().Err(err)
	l.handleArgs(event, args...)
	event.Msg(msg)
}

func (l *logger) Info(msg string, args ...any) {
	event := l.zlog.Info()
	l.handleArgs(event, args...)
	event.Msg(msg)
}

func (l *logger) Warning(msg string, args ...any) {
	event := l.zlog.Warn()
	l.handleArgs(event, args...)
	event.Msg(msg)
}

func (l *logger) Debug(msg string, args ...any) {
	event := l.zlog.Debug()
	l.handleArgs(event, args...)
	event.Msg(msg)
}

func (l *logger) NewGroup(group string) Logger {
	// Create a new logger with the "group" field added to the context
	zlog := l.zlog.With().Str("group", group).Logger()

	// Return a new instance of the logger with the scoped context
	return &logger{
		zlog: &zlog,
	}
}

func (l *logger) With(args ...any) Logger {
	nl := l.zlog.With()
	event := zerolog.Dict()
	l.handleArgs(event, args...)

	nl = nl.Dict("additional_info", event)

	zlog := nl.Logger()
	return &logger{
		zlog: &zlog,
	}
}
