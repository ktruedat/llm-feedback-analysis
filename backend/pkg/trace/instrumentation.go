package trace

import (
	"context"
	"net/http"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// Instrumented interface allows handlers to opt in to automatic HTTP instrumentation.
// Handlers implementing this interface will be automatically wrapped with otelhttp
// when tracing is enabled.
type Instrumented interface {
	IsTracingEnabled() bool
	SetTracing(enabled bool)
}

// InstrumentationOption configures HTTP handler instrumentation.
type InstrumentationOption func(Instrumented)

func WithTracingEnabled(enabled bool) InstrumentationOption {
	return func(i Instrumented) {
		i.SetTracing(enabled)
	}
}

// InstrumentHandlerFunc is a convenience wrapper for http.HandlerFunc.
// It accepts an optional Instrumented interface to check if tracing is enabled.
// If instrumented is nil or IsTracingEnabled() returns false, the handler is returned as-is.
func InstrumentHandlerFunc(
	handler http.HandlerFunc,
	operation string,
	instrumented Instrumented,
) http.HandlerFunc {
	if instrumented == nil || !instrumented.IsTracingEnabled() {
		return handler
	}

	// Use otelhttp to instrument the handler
	// otelhttp uses the global tracer provider set by our trace package
	wrapped := otelhttp.NewHandler(
		handler, operation, otelhttp.WithSpanNameFormatter(
			func(_ string, r *http.Request) string {
				return r.Method + " " + r.URL.Path
			},
		),
	)
	return wrapped.ServeHTTP
}

// InstrumentPgxPoolConfig instruments a pgxpool.Config with OpenTelemetry tracing.
// If instrumented is nil or IsTracingEnabled() returns false, the config is returned as-is.
// This hides otelpgx implementation details from application code.
func InstrumentPgxPoolConfig(
	connString string,
	instrumented Instrumented,
) (*pgxpool.Config, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	if instrumented != nil && instrumented.IsTracingEnabled() {
		// Attach OpenTelemetry tracer to pgx connection config
		// otelpgx uses the global tracer provider set by our trace package
		cfg.ConnConfig.Tracer = otelpgx.NewTracer()
	}

	return cfg, nil
}

// InstrumentPgxPool creates a new pgxpool.Pool with OpenTelemetry instrumentation.
// If instrumented is nil or IsTracingEnabled() returns false, creates a pool without instrumentation.
func InstrumentPgxPool(
	ctx context.Context,
	connString string,
	instrumented Instrumented,
) (*pgxpool.Pool, error) {
	cfg, err := InstrumentPgxPoolConfig(connString, instrumented)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	// Record database statistics if tracing is enabled
	if instrumented != nil && instrumented.IsTracingEnabled() {
		if err := otelpgx.RecordStats(pool); err != nil {
			pool.Close()
			return nil, err
		}
	}

	return pool, nil
}
