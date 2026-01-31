package analysis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/config"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/external"
	apprepo "github.com/ktruedat/llm-feedback-analysis/internal/app/repository"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/services"
	"github.com/ktruedat/llm-feedback-analysis/internal/domain/analysis"
	"github.com/ktruedat/llm-feedback-analysis/internal/domain/feedback"
	"github.com/ktruedat/llm-feedback-analysis/pkg/optional"
	"github.com/ktruedat/llm-feedback-analysis/pkg/tracelog"
)

const (
	defaultBufferSize    = 100
	defaultCheckInterval = 2 * time.Second
)

type analyzer struct {
	logger       tracelog.TraceLogger
	cfg          *config.LLMAnalysis
	analysisRepo apprepo.AnalysisRepository
	feedbackRepo apprepo.FeedbackRepository
	llmClient    external.LLMClient

	// Channel for receiving feedbacks (buffered to avoid blocking)
	feedbackChan chan *feedback.Feedback

	// Internal queue of feedbacks pending analysis
	pendingFeedbacks []*feedback.Feedback
	pendingMutex     sync.Mutex

	// Last analysis time for debounce and rate limiting
	lastAnalysisTime  time.Time
	lastAnalysisMutex sync.Mutex

	// Context and cancellation
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewAnalyzerService creates a new analyzer service. Rate limited by the configuration, for example
// achieving minimum number of new feedbacks before analysis is triggered and debounce time between analyses.
func NewAnalyzerService(
	logger tracelog.TraceLogger,
	cfg *config.LLMAnalysis,
	analysisRepo apprepo.AnalysisRepository,
	feedbackRepo apprepo.FeedbackRepository,
	llmClient external.LLMClient,
) services.AnalyzerService {
	// Buffered channel to avoid blocking feedback creation
	// Buffer size should be large enough to handle bursts
	bufferSize := cfg.MinimumNewFeedbacksForAnalysis * 2
	if bufferSize < defaultBufferSize {
		bufferSize = defaultBufferSize
	}

	return &analyzer{
		logger:           logger.NewGroup("llm_analyzer"),
		cfg:              cfg,
		analysisRepo:     analysisRepo,
		feedbackRepo:     feedbackRepo,
		llmClient:        llmClient,
		feedbackChan:     make(chan *feedback.Feedback, bufferSize),
		pendingFeedbacks: make([]*feedback.Feedback, 0, bufferSize),
	}
}

// EnqueueFeedback adds a feedback to the analysis queue.
// This method is non-blocking and runs in a separate goroutine.
func (a *analyzer) EnqueueFeedback(ctx context.Context, fb *feedback.Feedback) {
	// Send in a separate goroutine to avoid blocking even if buffer is full
	go func() {
		select {
		case a.feedbackChan <- fb:
			a.logger.Info("feedback enqueued for analysis", "feedback_id", fb.ID().String())
		case <-a.ctx.Done():
			a.logger.Info("analyzer stopped, dropping feedback", "feedback_id", fb.ID().String())
		}
	}()
}

// Start starts the analyzer service in a background goroutine.
func (a *analyzer) Start(ctx context.Context) error {
	a.logger.Info("starting LLM analyzer service")

	// Create a cancellable context from the provided context
	a.ctx, a.cancel = context.WithCancel(ctx)

	a.wg.Add(1)
	go a.run(a.ctx)

	return nil
}

// run is the main loop that processes feedbacks and triggers analysis.
func (a *analyzer) run(ctx context.Context) {
	defer a.wg.Done()

	ticker := time.NewTicker(defaultCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			a.logger.Info("analyzer context cancelled, stopping")
			return
		case fb := <-a.feedbackChan:
			a.addFeedbackToQueue(fb)
		case <-ticker.C:
			a.checkAndAnalyze(ctx)
		}
	}
}

// performAnalysis performs the actual LLM analysis.
func (a *analyzer) performAnalysis(ctx context.Context, feedbacks []*feedback.Feedback) {
	defer a.wg.Done()

	logger := a.logger.WithSpan(ctx)
	logger.Info("starting analysis", "feedback_count", len(feedbacks))

	// Get the latest analysis for incremental updates
	previousAnalysis, err := a.analysisRepo.GetLatest(ctx)
	if err != nil {
		logger.Info("no previous analysis found, starting fresh")
	}

	// Determine period
	periodStart := time.Now().UTC()
	periodEnd := time.Now().UTC()
	if len(feedbacks) > 0 {
		// Use the earliest feedback creation time as period start
		periodStart = feedbacks[0].CreatedAt()
		for _, fb := range feedbacks {
			if fb.CreatedAt().Before(periodStart) {
				periodStart = fb.CreatedAt()
			}
			if fb.CreatedAt().After(periodEnd) {
				periodEnd = fb.CreatedAt()
			}
		}
	}

	// Collect feedback IDs
	feedbackIDs := make([]uuid.UUID, len(feedbacks))
	for i, fb := range feedbacks {
		feedbackIDs[i] = fb.ID()
	}

	// Create analysis record with status 'processing'
	// Provide placeholder values for required fields that will be updated after LLM analysis
	analysisBuilder := analysis.NewBuilder().
		WithPeriod(periodStart, periodEnd).
		WithFeedbackCount(len(feedbacks)).
		WithOverallSummary("Processing...").
		WithSentiment(analysis.SentimentMixed).
		WithKeyInsights([]string{}).
		WithModel(a.cfg.OpenAIModel).
		WithTokens(0).
		WithAnalysisDurationMs(0).
		WithStatus(analysis.StatusProcessing)

	if previousAnalysis != nil {
		analysisBuilder.WithPreviousAnalysisID(previousAnalysis.ID())
		newCount := len(feedbacks)
		analysisBuilder.WithNewFeedbackCount(newCount)
	}

	analysisEntity, err := analysisBuilder.Build()
	if err != nil {
		logger.RecordSpanError(ctx, fmt.Errorf("failed to build analysis: %w", err))
		return
	}

	if err := a.analysisRepo.Create(ctx, analysisEntity); err != nil {
		logger.RecordSpanError(ctx, fmt.Errorf("failed to create analysis record: %w", err))
		return
	}
	logger.Info("analysis record created", "analysis_id", analysisEntity.ID().String())

	// Create analyzed feedback records (junction table)
	if err := a.analysisRepo.CreateAnalyzedFeedbacks(ctx, analysisEntity.ID(), feedbackIDs); err != nil {
		logger.RecordSpanError(ctx, fmt.Errorf("failed to create analyzed feedback records: %w", err))
		return
	}
	logger.Info(
		"analyzed feedback records created",
		"analysis_id",
		analysisEntity.ID().String(),
		"feedback_count",
		len(feedbackIDs),
	)

	// Call LLM client
	startTime := time.Now()
	var llmResult *external.AnalysisResult
	if a.llmClient != nil {
		llmResult, err = a.llmClient.AnalyzeFeedbacks(ctx, feedbacks, previousAnalysis)
	} else {
		// Stub implementation - return error for now
		err = fmt.Errorf("LLM client not implemented yet")
	}
	duration := time.Since(startTime)

	if err != nil {
		if err := analysisEntity.MarkFailed(err.Error()); err != nil {
			logger.RecordSpanError(ctx, fmt.Errorf("failed to mark analysis as failed: %w", err))
			return
		}
		if updateErr := a.analysisRepo.Update(
			ctx, analysisEntity.ID(), &analysis.UpdatableFields{
				FailureReason: analysisEntity.FailureReason(),
				Status:        analysis.StatusFailed,
				CompletedAt:   analysisEntity.CompletedAt().Unwrap(),
			},
		); updateErr != nil {
			logger.RecordSpanError(ctx, fmt.Errorf("failed to update analysis with failure: %w", updateErr))
		}
		logger.RecordSpanError(ctx, fmt.Errorf("LLM analysis failed: %w", err))
		return
	}

	// Log topics received from LLM
	logger.Info(
		"LLM analysis completed",
		"topics_count",
		len(llmResult.Topics),
		"overall_summary_length",
		len(llmResult.OverallSummary),
	)
	if len(llmResult.Topics) > 0 {
		for i, t := range llmResult.Topics {
			logger.Info(
				"LLM topic received",
				"index",
				i,
				"topic_enum",
				string(t.Topic),
				"feedback_ids_count",
				len(t.FeedbackIDs),
			)
		}
	} else {
		logger.Info("No topics returned from LLM")
	}

	// Use topics directly from LLM result (already converted)
	topics := llmResult.Topics
	logger.Info("topics array prepared", "topics_count", len(topics))

	// Update analysis with results
	if err := analysisEntity.MarkSuccess(); err != nil {
		logger.RecordSpanError(ctx, fmt.Errorf("failed to mark analysis as success: %w", err))
		return
	}
	logger.Info("analysis marked as success")

	updateBuilder := analysis.BuilderFromExisting(analysisEntity).
		WithOverallSummary(llmResult.OverallSummary).
		WithSentiment(llmResult.Sentiment).
		WithKeyInsights(llmResult.KeyInsights).
		WithTokens(llmResult.TokensUsed).
		WithAnalysisDurationMs(int(duration.Milliseconds()))

	updatedAnalysis, err := updateBuilder.Build()
	if err != nil {
		logger.RecordSpanError(ctx, fmt.Errorf("failed to build updated analysis: %w", err))
		return
	}
	logger.Info("updated analysis built successfully")

	logger.Info("updating analysis in database", "analysis_id", analysisEntity.ID().String())
	if err := a.analysisRepo.Update(
		ctx, analysisEntity.ID(), &analysis.UpdatableFields{
			Results: optional.Some(
				&analysis.UpdatedResults{
					OverallSummary: updatedAnalysis.OverallSummary(),
					Sentiment:      updatedAnalysis.Sentiment(),
					KeyInsights:    updatedAnalysis.KeyInsights(),
					Tokens:         updatedAnalysis.Tokens(),
				},
			),
			Status:      analysis.StatusSuccess,
			CompletedAt: updatedAnalysis.CompletedAt().Unwrap(),
		},
	); err != nil {
		logger.RecordSpanError(ctx, fmt.Errorf("failed to update analysis with results: %w", err))
		return
	}
	logger.Info("analysis updated in database successfully")

	// Create topics and their assignments
	logger.Info("creating topics", "topics_count", len(topics), "analysis_id", analysisEntity.ID().String())
	if err := a.createTopics(ctx, analysisEntity.ID(), topics, logger); err != nil {
		logger.Error(
			"failed to create topics",
			err,
			"topics_count",
			len(topics),
			"analysis_id",
			analysisEntity.ID().String(),
		)
		logger.RecordSpanError(ctx, err)
		// Don't return - analysis is already marked as success, topics are supplementary
	} else {
		logger.Info("topics creation completed successfully", "topics_count", len(topics))
	}

	logger.Info("analysis completed successfully", "analysis_id", updatedAnalysis.ID().String())

	// Note: Pending feedbacks are already managed in checkAndAnalyze
	// We only clear the ones that were selected for analysis, which is already done there

	a.lastAnalysisMutex.Lock()
	a.lastAnalysisTime = time.Now()
	a.lastAnalysisMutex.Unlock()
}

// createTopics creates topics and their feedback assignments for an analysis.
func (a *analyzer) createTopics(
	ctx context.Context,
	analysisID uuid.UUID,
	llmTopics []external.Topic,
	logger tracelog.TraceLogger,
) error {
	if len(llmTopics) == 0 {
		logger.Info("no topics to create (empty topics array)")
		return nil // No topics to create
	}

	logger.Info("starting topic creation", "topics_count", len(llmTopics), "analysis_id", analysisID.String())

	for i, llmTopic := range llmTopics {
		logger.Info(
			"processing topic",
			"index", i,
			"topic_enum", string(llmTopic.Topic),
			"summary_length", len(llmTopic.Summary),
			"feedback_ids_count", len(llmTopic.FeedbackIDs),
			"sentiment", string(llmTopic.Sentiment),
		)

		// Build topic analysis domain object
		topicAnalysisBuilder := analysis.NewTopicAnalysisBuilder().
			WithAnalysisID(analysisID).
			WithTopic(llmTopic.Topic).
			WithSummary(llmTopic.Summary).
			WithSentiment(llmTopic.Sentiment).
			WithFeedbackCount(len(llmTopic.FeedbackIDs))

		topicAnalysis, err := topicAnalysisBuilder.Build()
		if err != nil {
			buildErr := fmt.Errorf("failed to build topic analysis %s: %w", string(llmTopic.Topic), err)
			logger.Error("failed to build topic analysis", buildErr, "topic_enum", string(llmTopic.Topic), "index", i)
			logger.RecordSpanError(ctx, buildErr)
			return buildErr
		}

		logger.Info(
			"topic analysis built successfully",
			"topic_analysis_id",
			topicAnalysis.ID().String(),
			"topic_enum",
			string(topicAnalysis.Topic()),
		)

		// Create topic analysis in database
		if err := a.analysisRepo.CreateTopicAnalysis(ctx, topicAnalysis); err != nil {
			createErr := fmt.Errorf("failed to create topic analysis %s in database: %w", string(llmTopic.Topic), err)
			logger.Error(
				"failed to create topic analysis in database",
				createErr,
				"topic_analysis_id",
				topicAnalysis.ID().String(),
				"topic_enum",
				string(llmTopic.Topic),
			)
			logger.RecordSpanError(ctx, createErr)
			return createErr
		}

		logger.Info(
			"topic analysis created in database",
			"topic_analysis_id",
			topicAnalysis.ID().String(),
			"topic_enum",
			string(topicAnalysis.Topic()),
		)

		// Create feedback-topic assignments
		if len(llmTopic.FeedbackIDs) > 0 {
			logger.Info(
				"creating topic assignments",
				"topic_analysis_id", topicAnalysis.ID().String(),
				"feedback_ids_count", len(llmTopic.FeedbackIDs),
			)
			if err := a.analysisRepo.CreateTopicAssignments(
				ctx,
				analysisID,
				topicAnalysis.ID(),
				llmTopic.FeedbackIDs,
			); err != nil {
				assignErr := fmt.Errorf(
					"failed to create topic assignments for topic %s: %w",
					string(llmTopic.Topic),
					err,
				)
				logger.Error(
					"failed to create topic assignments",
					assignErr,
					"topic_analysis_id",
					topicAnalysis.ID().String(),
					"topic_enum",
					string(llmTopic.Topic),
					"feedback_ids_count",
					len(llmTopic.FeedbackIDs),
				)
				logger.RecordSpanError(ctx, assignErr)
				return assignErr
			}

			logger.Info(
				"topic assignments created successfully",
				"topic_analysis_id", topicAnalysis.ID().String(),
				"feedback_count", len(llmTopic.FeedbackIDs),
			)
		} else {
			logger.Info(
				"topic analysis has no feedback IDs, skipping assignments",
				"topic_analysis_id",
				topicAnalysis.ID().String(),
			)
		}
	}

	logger.Info("all topics created successfully", "topics_count", len(llmTopics))
	return nil
}
