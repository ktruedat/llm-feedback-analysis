package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/ktruedat/llm-feedback-analysis/pkg/tracelog"
)

// LoggerMiddleware creates a middleware that logs HTTP requests using TraceLogger.
// It integrates with chi's logger middleware pattern and adds tracing support.
func LoggerMiddleware(logger tracelog.TraceLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				start := time.Now()
				ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor) // Wrap the resp writer to capture status code

				next.ServeHTTP(ww, r)

				duration := time.Since(start)
				logger.Info(
					"http request completed",
					"method", r.Method,
					"path", r.URL.Path,
					"status", ww.Status(),
					"bytes", ww.BytesWritten(),
					"duration_ms", duration.Milliseconds(),
					"remote_addr", r.RemoteAddr,
				)
			},
		)
	}
}
