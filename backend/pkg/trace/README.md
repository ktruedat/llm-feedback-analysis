# Trace Package

Abstract tracing interface for distributed tracing with OpenTelemetry implementation. Supports pluggable backends while maintaining a consistent API.

## Quick Start

### Creating a Tracer

```go
tracer, err := trace.NewTracer(trace.Config{
    ServiceName:    "my-service",
    ServiceVersion: "1.0.0",
    Endpoint:       "http://tempo:4318", // OTLP endpoint (optional, uses no-op if empty)
    Insecure:       true,                // Use HTTP instead of HTTPS
    Environment:    "production",
})
if err != nil {
    panic(err)
}
defer tracer.Shutdown(context.Background())
```

**Note:** If `Endpoint` is empty, a no-op tracer is created for development/testing.

### Starting Spans

```go
ctx, span := tracer.Start(context.Background(), "operation-name",
    trace.WithSpanKind(oteltrace.SpanKindServer),
    trace.WithAttributes(
        trace.Attribute{Key: "user.id", Value: "123"},
        trace.Attribute{Key: "operation.type", Value: "read"},
    ),
)
defer span.End()
```

### Working with Spans

**Set attributes:**
```go
span.SetAttributes(
    trace.Attribute{Key: "result.count", Value: 42},
)
```

**Add events:**
```go
span.AddEvent("processing-complete",
    trace.WithEventAttributes(
        trace.Attribute{Key: "duration.ms", Value: 150},
    ),
)
```

**Record errors:**
```go
span.RecordError(err,
    trace.WithErrorAttributes(
        trace.Attribute{Key: "error.type", Value: "validation"},
    ),
)
span.SetStatus(trace.StatusError, err.Error())
span.End(trace.WithError(err))
```

### Nested Spans

```go
ctx, parentSpan := tracer.Start(ctx, "http.request",
    trace.WithSpanKind(oteltrace.SpanKindServer),
)

ctx, childSpan := tracer.Start(ctx, "database.query",
    trace.WithSpanKind(oteltrace.SpanKindClient),
)

childSpan.End()
parentSpan.End()
```

### Extracting Spans from Context

```go
span := trace.SpanFromContext(ctx)
if span != nil && span.IsRecording() {
    span.SetAttributes(trace.Attribute{Key: "key", Value: "value"})
    span.End()
}
```

### Span Status Codes

- `trace.StatusUnset` - Status is unset (default)
- `trace.StatusOK` - Operation completed successfully
- `trace.StatusError` - Operation ended with an error

### Supported Attribute Types

- `string`, `int`, `int64`, `float64`, `bool`
- `[]string`, `[]int`, `[]int64`, `[]float64`, `[]bool`
- Other types are converted to strings

### Span Context

```go
spanCtx := span.SpanContext()
traceID := spanCtx.TraceID()
spanID := spanCtx.SpanID()
isValid := spanCtx.IsValid()
```

