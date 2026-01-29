package trace

import (
	"context"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// SpanFromContext extracts a Span from the context.
// Returns nil if no span is found in the context or if the span is not recording.
// This function works with OpenTelemetry spans stored in the context.
func SpanFromContext(ctx context.Context) Span {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return nil
	}
	return &otelSpanFromContext{span: span}
}

// otelSpanFromContext wraps an OpenTelemetry span extracted from context.
type otelSpanFromContext struct {
	span trace.Span
}

// End completes the span.
func (s *otelSpanFromContext) End(opts ...SpanEndOption) {
	cfg := &spanEndConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.err != nil {
		s.span.RecordError(cfg.err)
		s.span.SetStatus(codes.Error, cfg.err.Error())
	}

	s.span.End()
}

// SetAttributes sets attributes on the span.
func (s *otelSpanFromContext) SetAttributes(attrs ...Attribute) {
	for _, attr := range attrs {
		s.span.SetAttributes(convertAttribute(attr))
	}
}

// AddEvent adds an event to the span.
func (s *otelSpanFromContext) AddEvent(name string, opts ...EventOption) {
	cfg := &eventConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	var otelOpts []trace.EventOption
	for _, attr := range cfg.attributes {
		otelOpts = append(otelOpts, trace.WithAttributes(convertAttribute(attr)))
	}

	s.span.AddEvent(name, otelOpts...)
}

// IsRecording returns whether the span is recording.
func (s *otelSpanFromContext) IsRecording() bool {
	return s.span.IsRecording()
}

// SpanContext returns the span context.
func (s *otelSpanFromContext) SpanContext() SpanContext {
	return &otelSpanContext{sc: s.span.SpanContext()}
}

// SetStatus sets the status of the span.
func (s *otelSpanFromContext) SetStatus(code StatusCode, description string) {
	var otelCode codes.Code
	switch code {
	case StatusOK:
		otelCode = codes.Ok
	case StatusError:
		otelCode = codes.Error
	case StatusUnset:
		otelCode = codes.Unset
	default:
		otelCode = codes.Unset
	}
	s.span.SetStatus(otelCode, description)
}

// RecordError records an error on the span.
func (s *otelSpanFromContext) RecordError(err error, opts ...ErrorOption) {
	cfg := &errorConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	var otelOpts []trace.EventOption
	for _, attr := range cfg.attributes {
		otelOpts = append(otelOpts, trace.WithAttributes(convertAttribute(attr)))
	}

	s.span.RecordError(err, otelOpts...)
}
