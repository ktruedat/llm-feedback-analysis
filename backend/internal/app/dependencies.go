package app

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/config"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/handlers"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/handlers/http/middleware"
	"github.com/ktruedat/llm-feedback-analysis/pkg/http/responder"
	"github.com/ktruedat/llm-feedback-analysis/pkg/log"
	"github.com/ktruedat/llm-feedback-analysis/pkg/trace"
	"github.com/ktruedat/llm-feedback-analysis/pkg/tracelog"
)

func New() (*App, error) {
	cfg, err := initConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize config: %w", err)
	}
	tracing, err := initTracing(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize tracing: %w", err)
	}

	restResponder := responder.NewRestResponder(tracing.baseLogger)
	app := &App{
		cfg:           cfg,
		router:        initRouter(cfg, tracing.traceLogger, restResponder),
		tracing:       tracing,
		restResponder: restResponder,
	}

	return app, nil
}

func initConfig() (*config.Config, error) {
	cfgPath := flag.String("cfg", "config.yaml", "path to the config file")
	flag.Parse()
	if cfgPath == nil {
		return nil, fmt.Errorf("config path is required")
	}

	return config.New(*cfgPath)
}

type tracing struct {
	tracer      trace.Tracer
	traceLogger tracelog.TraceLogger
	baseLogger  log.Logger
	enabled     bool
}

func (t *tracing) IsTracingEnabled() bool {
	return t.enabled
}

func (t *tracing) SetTracing(enabled bool) {
	t.enabled = enabled
}

func initTracing(cfg *config.Config) (*tracing, error) {
	logger := log.NewLogger(cfg.Profile.String(), log.WithCallersToSkip(4))

	var tracing tracing
	if cfg.Tracing.Enabled {
		tracer, err := trace.NewTracer(
			trace.Config{
				ServiceName: cfg.Tracing.ServiceName,
				Endpoint:    cfg.Tracing.OTelEndpoint,
				Insecure:    cfg.Tracing.Insecure,
				Environment: cfg.Profile.String(),
			},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize tracer: %w", err)
		}
		tracing.tracer = tracer
		tracing.baseLogger = logger
		tracing.traceLogger = tracelog.NewTraceLogger(tracing.baseLogger, tracer)
		tracing.enabled = true

		return &tracing, nil
	}

	// Use regular logger when tracing is disabled
	tracing.baseLogger = logger
	// Create a no-op trace logger that wraps the regular logger
	// This allows the consumers to use tracelog interface even when tracing is off
	noOpTracer, err := trace.NewTracer(
		trace.Config{
			ServiceName: cfg.Tracing.ServiceName,
			Environment: cfg.Profile.String(),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize no-op tracer: %w", err)
	}
	tracing.traceLogger = tracelog.NewTraceLogger(tracing.baseLogger, noOpTracer)
	tracing.enabled = false

	return &tracing, nil
}

func initRouter(cfg *config.Config, logger tracelog.TraceLogger, responder responder.RestResponder) *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		cors.Handler(middleware.CorsOptions()),
		middleware.JWTMiddleware(&cfg.JWT, logger, responder),
	)

	return router
}

type server struct {
	*http.Server
	registrableHandlers []handlers.Handlers
}

func newServer(
	cfg *config.Server,
	router *chi.Mux,
	registrableHandlers ...handlers.Handlers,
) *server {
	srv := &server{
		Server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Handler: router,
		},
		registrableHandlers: registrableHandlers,
	}

	for _, h := range srv.registrableHandlers {
		h.RegisterRoutes()
	}

	return srv
}
