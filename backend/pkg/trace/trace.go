// Package trace provides an abstract tracing interface for distributed tracing.
// It supports pluggable implementations, with OpenTelemetry as the default implementation.
package trace

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// Tracer is an abstract interface for distributed tracing.
// Implementations can use different tracing backends (OpenTelemetry, Jaeger, etc.)
// while maintaining a consistent API contract.
//
//go:generate mockgen -destination=../../mocking/trace_mock.go -package=mocking -source=trace.go Tracer
type Tracer interface {
	// Start creates a new span from the context with the given name.
	// Returns a new context containing the span and the span itself.
	Start(ctx context.Context, name string, opts ...SpanOption) (context.Context, Span)

	// Shutdown gracefully shuts down the tracer and flushes any pending spans.
	Shutdown(ctx context.Context) error
}

// Span represents a single operation in a trace.
type Span interface {
	// End completes the span and records its end time.
	End(opts ...SpanEndOption)

	// SetAttributes sets attributes on the span.
	SetAttributes(attrs ...Attribute)

	// AddEvent adds an event to the span.
	AddEvent(name string, opts ...EventOption)

	// IsRecording returns whether the span is recording events.
	IsRecording() bool

	// SpanContext returns the span context.
	SpanContext() SpanContext

	// SetStatus sets the status of the span.
	SetStatus(code StatusCode, description string)

	// RecordError records an error as an exception event on the span.
	RecordError(err error, opts ...ErrorOption)
}

// SpanContext contains the trace and span IDs.
type SpanContext interface {
	// TraceID returns the trace ID.
	TraceID() string

	// SpanID returns the span ID.
	SpanID() string

	// IsValid returns whether the span context is valid.
	IsValid() bool
}

// StatusCode represents the status code of a span.
type StatusCode int

const (
	// StatusUnset indicates the status is unset.
	StatusUnset StatusCode = iota
	// StatusOK indicates the operation completed successfully.
	StatusOK
	// StatusError indicates the operation ended with an error.
	StatusError
)

// Attribute represents a key-value pair attribute.
type Attribute struct {
	Key   string
	Value interface{}
}

// SpanOption configures a span.
type SpanOption func(*spanConfig)

// SpanEndOption configures span end behavior.
type SpanEndOption func(*spanEndConfig)

// EventOption configures an event.
type EventOption func(*eventConfig)

// ErrorOption configures error recording.
type ErrorOption func(*errorConfig)

type spanConfig struct {
	kind       trace.SpanKind
	attributes []Attribute
}

type spanEndConfig struct {
	err error
}

type eventConfig struct {
	attributes []Attribute
}

type errorConfig struct {
	attributes []Attribute
}

// WithSpanKind sets the span kind.
func WithSpanKind(kind trace.SpanKind) SpanOption {
	return func(cfg *spanConfig) {
		cfg.kind = kind
	}
}

// WithAttributes sets attributes on the span at creation time.
func WithAttributes(attrs ...Attribute) SpanOption {
	return func(cfg *spanConfig) {
		cfg.attributes = append(cfg.attributes, attrs...)
	}
}

// WithError sets an error when ending the span.
func WithError(err error) SpanEndOption {
	return func(cfg *spanEndConfig) {
		cfg.err = err
	}
}

// WithEventAttributes sets attributes on an event.
func WithEventAttributes(attrs ...Attribute) EventOption {
	return func(cfg *eventConfig) {
		cfg.attributes = append(cfg.attributes, attrs...)
	}
}

// WithErrorAttributes sets attributes when recording an error.
func WithErrorAttributes(attrs ...Attribute) ErrorOption {
	return func(cfg *errorConfig) {
		cfg.attributes = append(cfg.attributes, attrs...)
	}
}
