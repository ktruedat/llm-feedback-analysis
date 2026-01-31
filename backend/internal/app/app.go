package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/config"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/external/llm"
	handlersv1 "github.com/ktruedat/llm-feedback-analysis/internal/app/handlers/http/v1"
	analysisRepository "github.com/ktruedat/llm-feedback-analysis/internal/app/repository/postgres/analysis"
	feedbackRepository "github.com/ktruedat/llm-feedback-analysis/internal/app/repository/postgres/feedback"
	userRepository "github.com/ktruedat/llm-feedback-analysis/internal/app/repository/postgres/user"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/services"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/services/analysis"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/services/feedback"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/services/user"
	"github.com/ktruedat/llm-feedback-analysis/migrations"
	ce "github.com/ktruedat/llm-feedback-analysis/pkg/errors"
	"github.com/ktruedat/llm-feedback-analysis/pkg/http/responder"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository/sql"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository/sql/querier"
	"github.com/ktruedat/llm-feedback-analysis/pkg/trace"
)

type App struct {
	cfg           *config.Config
	router        *chi.Mux
	tracing       *tracing
	pgxPool       *pgxpool.Pool
	srv           *server
	restResponder responder.RestResponder
	analyzer      services.AnalyzerService
}

func (app *App) Start() error {
	logger := app.tracing.traceLogger

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pgxPool, err := trace.InstrumentPgxPool(ctx, app.cfg.DB.DSN, app.tracing)
	if err != nil {
		return fmt.Errorf("failed to create pgx pool: %w", err)
	}
	app.pgxPool = pgxPool

	if err := migrations.ApplyMigrations(pgxPool, logger); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	q := querier.NewPgxPool(pgxPool)
	feedbackRepo := feedbackRepository.NewFeedbackRepository(q)
	userRepo := userRepository.NewUserRepository(q)
	analysisRepo := analysisRepository.NewAnalysisRepository(q)

	errChecker := ce.NewErrorChecker()
	transactor := sql.NewTransactionManager(pgxPool)

	// Create OpenAI LLM client
	llmClient := llm.NewOpenAIClient(
		app.cfg.LLMAnalysis.OpenAIAPIKey,
		app.cfg.LLMAnalysis.OpenAIModel,
		logger,
	)

	// Create analyzer service (performs analysis)
	analyzerSvc := analysis.NewAnalyzerService(
		logger,
		&app.cfg.LLMAnalysis,
		analysisRepo,
		feedbackRepo,
		llmClient,
	)
	app.analyzer = analyzerSvc

	// Start analyzer in background with signal context
	// This ensures the analyzer stops gracefully when the app receives shutdown signals
	if err := analyzerSvc.Start(ctx); err != nil {
		return fmt.Errorf("failed to start analyzer: %w", err)
	}

	feedbackSvc := feedback.NewFeedbackService(
		logger,
		&app.cfg.Pagination,
		errChecker,
		feedbackRepo,
		transactor,
		analyzerSvc,
	)
	userSvc := user.NewUserService(logger, errChecker, userRepo, &app.cfg.JWT, transactor)

	feedbackSummarySvc := analysis.NewFeedbackSummaryService(logger, analysisRepo, feedbackRepo)
	feedbackV1Handlers := handlersv1.NewHandlers(
		app.router,
		logger,
		app.restResponder,
		feedbackSvc,
		userSvc,
		feedbackSummarySvc,
		&app.cfg.JWT,
		trace.WithTracingEnabled(app.cfg.Tracing.Enabled),
	)

	srv := newServer(&app.cfg.Server, app.router, feedbackV1Handlers)
	app.srv = srv
	go func() {
		log.Printf("starting server on port: %d", app.cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to start server: %s", err)
		}
	}()
	<-ctx.Done()

	// Shutdown application gracefully
	ctx, cancel := context.WithTimeout(ctx, time.Duration(app.cfg.Server.GracefulShutdownSeconds)*time.Second)
	defer cancel()

	if err := app.Close(ctx); err != nil {
		return fmt.Errorf("failed to close application: %w", err)
	}

	return nil
}

func (app *App) Close(ctx context.Context) error {
	if err := app.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown the server: %w", err)
	}

	// Stop analyzer gracefully
	if app.analyzer != nil {
		if err := app.analyzer.Stop(ctx); err != nil {
			return fmt.Errorf("failed to stop analyzer: %w", err)
		}
	}

	if app.pgxPool != nil {
		app.pgxPool.Close()
	}

	if app.tracing.tracer != nil {
		if err := app.tracing.tracer.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown tracer: %w", err)
		}
	}

	return nil
}
