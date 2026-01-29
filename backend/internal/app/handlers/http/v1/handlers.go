package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/handlers"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/services"
	"github.com/ktruedat/llm-feedback-analysis/pkg/http/responder"
	"github.com/ktruedat/llm-feedback-analysis/pkg/trace"
	"github.com/ktruedat/llm-feedback-analysis/pkg/tracelog"
)

type Handlers struct {
	r               chi.Router
	logger          tracelog.TraceLogger
	responder       responder.RestResponder
	feedbackService services.FeedbackService
	tracingEnabled  bool
}

func NewHandlers(
	r chi.Router,
	logger tracelog.TraceLogger,
	responder responder.RestResponder,
	feedbackService services.FeedbackService,
	opts ...trace.InstrumentationOption,
) handlers.Handlers {
	h := &Handlers{
		r:               r,
		logger:          logger.NewGroup("handlers_http"),
		responder:       responder,
		feedbackService: feedbackService,
		tracingEnabled:  false,
	}

	// Apply instrumentation options to the handler
	for _, opt := range opts {
		opt(h)
	}

	return h
}

func (h *Handlers) RegisterRoutes() {
	h.registerFeedbackRoutes(h.r)
}

// SetTracing implements trace.Instrumented interface.
func (h *Handlers) SetTracing(enabled bool) {
	h.tracingEnabled = enabled
}

// IsTracingEnabled implements trace.Instrumented interface.
// This allows automatic HTTP instrumentation via trace.InstrumentHandler.
func (h *Handlers) IsTracingEnabled() bool {
	return h.tracingEnabled
}
