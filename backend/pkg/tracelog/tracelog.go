// Package tracelog provides a hybrid logger that combines standard logging
// with distributed tracing capabilities. It integrates with OpenTelemetry
// to add context to spans while maintaining the Logger interface contract.
package tracelog

import (
	"context"

	"github.com/ktruedat/llm-feedback-analysis/pkg/log"
	"github.com/ktruedat/llm-feedback-analysis/pkg/trace"
)

// TraceLogger is a hybrid logger that combines standard logging with distributed tracing.
// It implements the Logger interface while also providing tracing capabilities.
type TraceLogger interface {
	Error(msg string, err error, args ...any)
	Info(msg string, args ...any)
	Warning(msg string, args ...any)
	Debug(msg string, args ...any)

	// NewGroup returns a new TraceLogger with the provided group.
	// This overrides log.Logger.NewGroup to return TraceLogger instead of log.Logger.
	NewGroup(group string) TraceLogger

	// With returns a new TraceLogger with the provided values.
	// This overrides log.Logger.With to return TraceLogger instead of log.Logger.
	With(args ...any) TraceLogger

	// StartSpan creates a new span and returns a context-aware logger.
	// The logger will automatically add trace context to log messages.
	StartSpan(ctx context.Context, name string, opts ...trace.SpanOption) (context.Context, TraceLogger, trace.Span)

	// WithSpan returns a logger that operates within an existing span context.
	// If no span exists in the context, it falls back to standard logging.
	WithSpan(ctx context.Context) TraceLogger

	// SetSpanAttributes sets attributes on the current span in the context.
	SetSpanAttributes(ctx context.Context, attrs ...trace.Attribute)

	// AddSpanEvent adds an event to the current span in the context.
	AddSpanEvent(ctx context.Context, name string, opts ...trace.EventOption)

	// RecordSpanError records an error on the current span in the context.
	RecordSpanError(ctx context.Context, err error, opts ...trace.ErrorOption)
}

// traceLogger implements TraceLogger.
type traceLogger struct {
	logger log.Logger
	tracer trace.Tracer
}

// NewTraceLogger creates a new TraceLogger that combines a Logger with a Tracer.
func NewTraceLogger(logger log.Logger, tracer trace.Tracer) TraceLogger {
	return &traceLogger{
		logger: logger,
		tracer: tracer,
	}
}

// StartSpan creates a new span and returns a context-aware logger.
func (tl *traceLogger) StartSpan(ctx context.Context, name string, opts ...trace.SpanOption) (
	context.Context,
	TraceLogger,
	trace.Span,
) {
	ctx, span := tl.tracer.Start(ctx, name, opts...)
	return ctx, &traceLoggerWithSpan{logger: tl.logger, tracer: tl.tracer, span: span}, span
}

// WithSpan returns a logger that operates within an existing span context.
func (tl *traceLogger) WithSpan(ctx context.Context) TraceLogger {
	span := trace.SpanFromContext(ctx)
	if span == nil {
		// No span in context, return standard logger wrapper
		return &traceLoggerWithSpan{logger: tl.logger, tracer: tl.tracer, span: nil}
	}
	return &traceLoggerWithSpan{logger: tl.logger, tracer: tl.tracer, span: span}
}

// SetSpanAttributes sets attributes on the current span in the context.
func (*traceLogger) SetSpanAttributes(ctx context.Context, attrs ...trace.Attribute) {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		span.SetAttributes(attrs...)
	}
}

// AddSpanEvent adds an event to the current span in the context.
func (*traceLogger) AddSpanEvent(ctx context.Context, name string, opts ...trace.EventOption) {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		span.AddEvent(name, opts...)
	}
}

// RecordSpanError records an error on the current span in the context.
func (*traceLogger) RecordSpanError(ctx context.Context, err error, opts ...trace.ErrorOption) {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		span.RecordError(err, opts...)
	}
}

// Logger interface implementation - delegate to underlying logger

func (tl *traceLogger) Error(msg string, err error, args ...any) {
	tl.logger.Error(msg, err, args...)
}

func (tl *traceLogger) Info(msg string, args ...any) {
	tl.logger.Info(msg, args...)
}

func (tl *traceLogger) Warning(msg string, args ...any) {
	tl.logger.Warning(msg, args...)
}

func (tl *traceLogger) Debug(msg string, args ...any) {
	tl.logger.Debug(msg, args...)
}

func (tl *traceLogger) NewGroup(group string) TraceLogger {
	return &traceLogger{
		logger: tl.logger.NewGroup(group),
		tracer: tl.tracer,
	}
}

func (tl *traceLogger) With(args ...any) TraceLogger {
	return &traceLogger{
		logger: tl.logger.With(args...),
		tracer: tl.tracer,
	}
}

// traceLoggerWithSpan is a logger that has an associated span.
type traceLoggerWithSpan struct {
	logger log.Logger
	tracer trace.Tracer
	span   trace.Span
}

// StartSpan creates a new child span.
func (tl *traceLoggerWithSpan) StartSpan(ctx context.Context, name string, opts ...trace.SpanOption) (
	context.Context,
	TraceLogger,
	trace.Span,
) {
	if tl.span != nil {
		// Add parent span context to options
		if sc := tl.span.SpanContext(); sc.IsValid() {
			opts = append(
				opts, trace.WithAttributes(
					trace.Attribute{Key: "parent.trace_id", Value: sc.TraceID()},
					trace.Attribute{Key: "parent.span_id", Value: sc.SpanID()},
				),
			)
		}
	}
	ctx, span := tl.tracer.Start(ctx, name, opts...)
	return ctx, &traceLoggerWithSpan{logger: tl.logger, tracer: tl.tracer, span: span}, span
}

// WithSpan returns itself if it already has a span, otherwise creates a new one from context.
func (tl *traceLoggerWithSpan) WithSpan(ctx context.Context) TraceLogger {
	if tl.span != nil {
		return tl
	}
	span := trace.SpanFromContext(ctx)
	return &traceLoggerWithSpan{logger: tl.logger, tracer: tl.tracer, span: span}
}

// SetSpanAttributes sets attributes on the current span.
func (tl *traceLoggerWithSpan) SetSpanAttributes(ctx context.Context, attrs ...trace.Attribute) {
	if tl.span != nil {
		tl.span.SetAttributes(attrs...)
	} else {
		span := trace.SpanFromContext(ctx)
		if span != nil {
			span.SetAttributes(attrs...)
		}
	}
}

// AddSpanEvent adds an event to the current span.
func (tl *traceLoggerWithSpan) AddSpanEvent(ctx context.Context, name string, opts ...trace.EventOption) {
	if tl.span != nil {
		tl.span.AddEvent(name, opts...)
	} else {
		span := trace.SpanFromContext(ctx)
		if span != nil {
			span.AddEvent(name, opts...)
		}
	}
}

// RecordSpanError records an error on the current span.
func (tl *traceLoggerWithSpan) RecordSpanError(ctx context.Context, err error, opts ...trace.ErrorOption) {
	if tl.span != nil {
		tl.span.RecordError(err, opts...)
	} else {
		span := trace.SpanFromContext(ctx)
		if span != nil {
			span.RecordError(err, opts...)
		}
	}
}

// Logger interface implementation with span context

func (tl *traceLoggerWithSpan) Error(msg string, err error, args ...any) {
	// Add trace context to log message
	if tl.span != nil {
		sc := tl.span.SpanContext()
		if sc.IsValid() {
			args = append(args, "trace_id", sc.TraceID(), "span_id", sc.SpanID())
		}
		// Record error on span
		tl.span.RecordError(err)
		tl.span.SetStatus(trace.StatusError, err.Error())
	}
	tl.logger.Error(msg, err, args...)
}

func (tl *traceLoggerWithSpan) Info(msg string, args ...any) {
	// Add trace context to log message
	if tl.span != nil {
		sc := tl.span.SpanContext()
		if sc.IsValid() {
			args = append(args, "trace_id", sc.TraceID(), "span_id", sc.SpanID())
		}
		// Add as event to span
		tl.span.AddEvent(
			msg, trace.WithEventAttributes(
				trace.Attribute{Key: "log.level", Value: "info"},
			),
		)
	}
	tl.logger.Info(msg, args...)
}

func (tl *traceLoggerWithSpan) Warning(msg string, args ...any) {
	// Add trace context to log message
	if tl.span != nil {
		sc := tl.span.SpanContext()
		if sc.IsValid() {
			args = append(args, "trace_id", sc.TraceID(), "span_id", sc.SpanID())
		}
		// Add as event to span
		tl.span.AddEvent(
			msg, trace.WithEventAttributes(
				trace.Attribute{Key: "log.level", Value: "warning"},
			),
		)
	}
	tl.logger.Warning(msg, args...)
}

func (tl *traceLoggerWithSpan) Debug(msg string, args ...any) {
	// Add trace context to log message
	if tl.span != nil {
		sc := tl.span.SpanContext()
		if sc.IsValid() {
			args = append(args, "trace_id", sc.TraceID(), "span_id", sc.SpanID())
		}
		// Add as event to span
		tl.span.AddEvent(
			msg, trace.WithEventAttributes(
				trace.Attribute{Key: "log.level", Value: "debug"},
			),
		)
	}
	tl.logger.Debug(msg, args...)
}

func (tl *traceLoggerWithSpan) NewGroup(group string) TraceLogger {
	return &traceLoggerWithSpan{
		logger: tl.logger.NewGroup(group),
		tracer: tl.tracer,
		span:   tl.span,
	}
}

func (tl *traceLoggerWithSpan) With(args ...any) TraceLogger {
	return &traceLoggerWithSpan{
		logger: tl.logger.With(args...),
		tracer: tl.tracer,
		span:   tl.span,
	}
}

// AsLogger wraps a TraceLogger as a log.Logger.
// This allows TraceLogger to be used where log.Logger is expected.
// Note: NewGroup() and With() on the returned Logger will return log.Logger, not TraceLogger.
func AsLogger(tl TraceLogger) log.Logger {
	return &loggerWrapper{traceLogger: tl}
}

// loggerWrapper wraps TraceLogger to implement log.Logger interface.
type loggerWrapper struct {
	traceLogger TraceLogger
}

func (lw *loggerWrapper) Error(msg string, err error, args ...any) {
	lw.traceLogger.Error(msg, err, args...)
}

func (lw *loggerWrapper) Info(msg string, args ...any) {
	lw.traceLogger.Info(msg, args...)
}

func (lw *loggerWrapper) Warning(msg string, args ...any) {
	lw.traceLogger.Warning(msg, args...)
}

func (lw *loggerWrapper) Debug(msg string, args ...any) {
	lw.traceLogger.Debug(msg, args...)
}

func (lw *loggerWrapper) NewGroup(group string) log.Logger {
	// Convert TraceLogger back to log.Logger
	return AsLogger(lw.traceLogger.NewGroup(group))
}

func (lw *loggerWrapper) With(args ...any) log.Logger {
	// Convert TraceLogger back to log.Logger
	return AsLogger(lw.traceLogger.With(args...))
}
