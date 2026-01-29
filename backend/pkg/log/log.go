// Package log provides a common logger interface for logging messages.
package log

// Logger interface for logging messages.
//
//go:generate mockgen -destination=../../mocking/log_mock.go -package mocking -source log.go Logger
type Logger interface {
	// Error logs an error message with the provided error and additional arguments.
	Error(msg string, err error, args ...any)
	// Info logs an info message with additional arguments.
	Info(msg string, args ...any)
	// Warning logs a warning message with additional arguments.
	Warning(msg string, args ...any)
	// Debug logs a debug message with additional arguments.
	Debug(msg string, args ...any)
	// NewGroup returns a new logger with the provided group.
	NewGroup(group string) Logger
	// With returns a new logger with the provided values.
	With(args ...any) Logger
}
