package analysis

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	apprepo "github.com/ktruedat/llm-feedback-analysis/internal/app/repository"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/services"
	"github.com/ktruedat/llm-feedback-analysis/internal/domain/analysis"
	"github.com/ktruedat/llm-feedback-analysis/pkg/tracelog"
)

// service provides read-only access to analysis data.
// This is separate from AnalyzerService which performs the actual analysis.
type service struct {
	logger       tracelog.TraceLogger
	analysisRepo apprepo.AnalysisRepository
	feedbackRepo apprepo.FeedbackRepository
}

// NewFeedbackSummaryService creates a new analysis service.
func NewFeedbackSummaryService(
	logger tracelog.TraceLogger,
	analysisRepo apprepo.AnalysisRepository,
	feedbackRepo apprepo.FeedbackRepository,
) services.FeedbackSummaryService {
	return &service{
		logger:       logger.NewGroup("feedback_summary_service"),
		analysisRepo: analysisRepo,
		feedbackRepo: feedbackRepo,
	}
}

// GetLatestAnalysis retrieves the latest completed analysis.
func (s *service) GetLatestAnalysis(ctx context.Context) (*analysis.Analysis, error) {
	logger := s.logger.WithSpan(ctx)
	logger.Info("getting latest analysis")

	latestAnalysis, err := s.analysisRepo.GetLatest(ctx)
	if err != nil {
		logger.RecordSpanError(ctx, err)
		logger.Error("error getting latest analysis", err)
		return nil, fmt.Errorf("failed to get latest analysis: %w", err)
	}

	if latestAnalysis == nil {
		logger.Info("no analysis found")
		return nil, nil
	}

	logger.Info("latest analysis retrieved", "analysis_id", latestAnalysis.ID().String())
	return latestAnalysis, nil
}

// GetAllAnalyses retrieves all analyses ordered by creation date (newest first).
func (s *service) GetAllAnalyses(ctx context.Context) ([]*analysis.Analysis, error) {
	logger := s.logger.WithSpan(ctx)
	logger.Info("getting all analyses")

	analyses, err := s.analysisRepo.List(ctx)
	if err != nil {
		logger.RecordSpanError(ctx, err)
		logger.Error("error getting all analyses", err)
		return nil, fmt.Errorf("failed to get analyses: %w", err)
	}

	logger.Info("all analyses retrieved", "count", len(analyses))
	return analyses, nil
}

// GetAnalysisByID retrieves an analysis by ID with its topics and analyzed feedbacks with their topics.
func (s *service) GetAnalysisByID(ctx context.Context, analysisID uuid.UUID) (
	*analysis.Analysis,
	[]*analysis.TopicAnalysis,
	map[uuid.UUID][]*analysis.TopicAnalysis, // feedback ID -> topics
	error,
) {
	logger := s.logger.WithSpan(ctx)
	logger.Info("getting analysis by ID", "analysis_id", analysisID.String())

	// Get the analysis
	analysisEntity, err := s.analysisRepo.GetByID(ctx, analysisID)
	if err != nil {
		logger.RecordSpanError(ctx, err)
		logger.Error("error getting analysis", err, "analysis_id", analysisID)
		return nil, nil, nil, fmt.Errorf("failed to get analysis: %w", err)
	}

	// Get topics for this analysis
	topics, err := s.analysisRepo.GetTopicsByAnalysisID(ctx, analysisID)
	if err != nil {
		logger.RecordSpanError(ctx, err)
		logger.Error("error getting topics", err, "analysis_id", analysisID)
		return nil, nil, nil, fmt.Errorf("failed to get topics: %w", err)
	}

	// Get feedback IDs analyzed in this analysis
	feedbackIDs, err := s.analysisRepo.GetFeedbackIDsByAnalysisID(ctx, analysisID)
	if err != nil {
		logger.RecordSpanError(ctx, err)
		logger.Error("error getting feedback IDs", err, "analysis_id", analysisID)
		return nil, nil, nil, fmt.Errorf("failed to get feedback IDs: %w", err)
	}

	// Build map of feedback ID -> topics
	feedbackTopics := make(map[uuid.UUID][]*analysis.TopicAnalysis)
	for _, topic := range topics {
		topicFeedbackIDs, err := s.analysisRepo.GetFeedbackIDsByTopicID(ctx, topic.ID())
		if err != nil {
			logger.Warning(
				"error getting feedback IDs for topic",
				"topic_id",
				topic.ID().String(),
				"error",
				err.Error(),
			)
			continue
		}
		for _, fbID := range topicFeedbackIDs {
			feedbackTopics[fbID] = append(feedbackTopics[fbID], topic)
		}
	}

	logger.Info(
		"analysis retrieved",
		"analysis_id",
		analysisID.String(),
		"topics_count",
		len(topics),
		"feedbacks_count",
		len(feedbackIDs),
	)
	return analysisEntity, topics, feedbackTopics, nil
}
