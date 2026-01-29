# TraceLog Package

Hybrid logger that combines standard logging with distributed tracing. Automatically adds trace context to log messages and integrates with OpenTelemetry spans.

## Quick Start

### Creating a TraceLogger

**Manual construction:**
```go
logger := log.NewLogger("development")
tracer, _ := trace.NewTracer(trace.Config{ServiceName: "my-service"})
traceLogger := tracelog.NewTraceLogger(logger, tracer)
```

### Basic Logging

TraceLogger implements the standard Logger interface:

```go
traceLogger.Info("application started", "version", "1.0.0")
traceLogger.Warning("deprecated feature used", "feature", "old-api")
traceLogger.Debug("processing data", "count", 42)
traceLogger.Error("operation failed", err, "user", "alice")
```

### Starting Spans

Create a span and get a context-aware logger:

```go
ctx, spanLogger, span := traceLogger.StartSpan(context.Background(), "process-request",
    trace.WithAttributes(
        trace.Attribute{Key: "request.id", Value: "req-123"},
    ),
)
defer span.End()

// Log messages automatically include trace_id and span_id
spanLogger.Info("processing request", "user", "alice")
spanLogger.Error("validation failed", err, "field", "email")
```

**Note:** When logging with a span-aware logger:
- Log messages automatically include `trace_id` and `span_id`
- `Error()` calls automatically record errors on the span
- `Info()`, `Warning()`, `Debug()` calls automatically add events to the span

### Working with Existing Spans

Use a logger that operates within an existing span context (e.g., from HTTP middleware):

```go
// Extract span from context
spanLogger := traceLogger.WithSpan(ctx)

// Log messages will include trace context if span exists
spanLogger.Info("operation completed", "result", "success")
```

### Span Operations

**Set attributes:**
```go
traceLogger.SetSpanAttributes(ctx,
    trace.Attribute{Key: "processing.stage", Value: "validation"},
)
```

**Add events:**
```go
traceLogger.AddSpanEvent(ctx, "checkpoint-reached",
    trace.WithEventAttributes(
        trace.Attribute{Key: "items.processed", Value: 100},
    ),
)
```

**Record errors:**
```go
traceLogger.RecordSpanError(ctx, err,
    trace.WithErrorAttributes(
        trace.Attribute{Key: "error.type", Value: "validation"},
    ),
)
```

### Logger Composition

**Create groups:**
```go
dbLogger := traceLogger.NewGroup("database")
dbLogger.Info("query executed", "table", "users")
```

**Add context:**
```go
userLogger := traceLogger.With("user_id", "123")
userLogger.Info("user action", "action", "login")
```

### Converting to Standard Logger

Use `AsLogger()` to convert TraceLogger to log.Logger when needed:

```go
standardLogger := tracelog.AsLogger(traceLogger)
// standardLogger.NewGroup() and With() return log.Logger, not TraceLogger
```

### Key Features

- **Automatic trace context**: Log messages from span-aware loggers include trace/span IDs
- **Error recording**: `Error()` calls automatically record errors on spans
- **Event integration**: Log calls automatically add events to spans
- **Context-aware**: Works seamlessly with spans from HTTP middleware or other sources
- **Backward compatible**: Implements standard Logger interface

