package trace

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// Config holds configuration for the OpenTelemetry tracer.
type Config struct {
	// ServiceName is the name of the service being traced.
	ServiceName string

	// ServiceVersion is the version of the service.
	ServiceVersion string

	// Endpoint is the OTLP endpoint URL (e.g., http://tempo:4318).
	// If empty, uses stdout exporter for development.
	Endpoint string

	// Insecure determines whether to use insecure connection (HTTP instead of HTTPS).
	Insecure bool

	// Environment is the deployment environment (e.g., "production", "development").
	Environment string
}

// otelTracer implements the Tracer interface using OpenTelemetry.
type otelTracer struct {
	tracer trace.Tracer
	tp     *sdktrace.TracerProvider
}

// otelSpan wraps OpenTelemetry's span.
type otelSpan struct {
	span trace.Span
}

// otelSpanContext wraps OpenTelemetry's span context.
type otelSpanContext struct {
	sc trace.SpanContext
}

// NewTracer creates a new OpenTelemetry-based tracer.
// If endpoint is empty, it creates a no-op tracer for development.
func NewTracer(cfg Config) (Tracer, error) {
	if cfg.ServiceName == "" {
		cfg.ServiceName = "unknown-service"
	}

	// Create resource with service information
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", cfg.ServiceName),
			attribute.String("service.version", cfg.ServiceVersion),
			attribute.String("deployment.environment", cfg.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	var tp *sdktrace.TracerProvider
	if cfg.Endpoint == "" {
		// Use no-op tracer provider for development/testing
		tp = sdktrace.NewTracerProvider(
			sdktrace.WithResource(res),
		)
	} else {
		// Parse endpoint URL to extract host:port
		// WithEndpoint expects host:port without scheme
		endpoint := cfg.Endpoint
		if strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://") {
			parsedURL, err := url.Parse(endpoint)
			if err != nil {
				return nil, fmt.Errorf("failed to parse endpoint URL: %w", err)
			}
			endpoint = parsedURL.Host
		}

		// Create OTLP HTTP exporter for Grafana Tempo
		opts := []otlptracehttp.Option{
			otlptracehttp.WithEndpoint(endpoint),
		}
		if cfg.Insecure {
			opts = append(opts, otlptracehttp.WithInsecure())
		}

		exporter, err := otlptracehttp.New(context.Background(), opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
		}

		// Create tracer provider with exporter
		tp = sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(res),
		)
	}

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	// Set global propagator
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	tracer := tp.Tracer(cfg.ServiceName)

	return &otelTracer{
		tracer: tracer,
		tp:     tp,
	}, nil
}

// Start creates a new span.
func (t *otelTracer) Start(ctx context.Context, name string, opts ...SpanOption) (context.Context, Span) {
	cfg := &spanConfig{
		kind: trace.SpanKindInternal,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	otelOpts := []trace.SpanStartOption{
		trace.WithSpanKind(cfg.kind),
	}

	// Convert attributes
	for _, attr := range cfg.attributes {
		otelOpts = append(otelOpts, trace.WithAttributes(convertAttribute(attr)))
	}

	ctx, span := t.tracer.Start(ctx, name, otelOpts...)

	// Set attributes that were provided
	for _, attr := range cfg.attributes {
		span.SetAttributes(convertAttribute(attr))
	}

	return ctx, &otelSpan{span: span}
}

// Shutdown gracefully shuts down the tracer.
func (t *otelTracer) Shutdown(ctx context.Context) error {
	return t.tp.Shutdown(ctx)
}

// End completes the span.
func (s *otelSpan) End(opts ...SpanEndOption) {
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
func (s *otelSpan) SetAttributes(attrs ...Attribute) {
	for _, attr := range attrs {
		s.span.SetAttributes(convertAttribute(attr))
	}
}

// AddEvent adds an event to the span.
func (s *otelSpan) AddEvent(name string, opts ...EventOption) {
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
func (s *otelSpan) IsRecording() bool {
	return s.span.IsRecording()
}

// SpanContext returns the span context.
func (s *otelSpan) SpanContext() SpanContext {
	return &otelSpanContext{sc: s.span.SpanContext()}
}

// SetStatus sets the status of the span.
func (s *otelSpan) SetStatus(code StatusCode, description string) {
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
func (s *otelSpan) RecordError(err error, opts ...ErrorOption) {
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

// TraceID returns the trace ID.
func (sc *otelSpanContext) TraceID() string {
	return sc.sc.TraceID().String()
}

// SpanID returns the span ID.
func (sc *otelSpanContext) SpanID() string {
	return sc.sc.SpanID().String()
}

// IsValid returns whether the span context is valid.
func (sc *otelSpanContext) IsValid() bool {
	return sc.sc.IsValid()
}

// convertAttribute converts our Attribute to OpenTelemetry's attribute.
func convertAttribute(attr Attribute) attribute.KeyValue {
	key := attribute.Key(attr.Key)

	switch v := attr.Value.(type) {
	case string:
		return key.String(v)
	case int:
		return key.Int(v)
	case int64:
		return key.Int64(v)
	case float64:
		return key.Float64(v)
	case bool:
		return key.Bool(v)
	case []string:
		return key.StringSlice(v)
	case []int:
		return key.IntSlice(v)
	case []int64:
		return key.Int64Slice(v)
	case []float64:
		return key.Float64Slice(v)
	case []bool:
		return key.BoolSlice(v)
	default:
		return key.String(fmt.Sprintf("%v", v))
	}
}
