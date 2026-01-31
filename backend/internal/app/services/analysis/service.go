package analysis

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	apprepo "github.com/ktruedat/llm-feedback-analysis/internal/app/repository"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/services"
	"github.com/ktruedat/llm-feedback-analysis/internal/domain/analysis"
	"github.com/ktruedat/llm-feedback-analysis/internal/domain/feedback"
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

// GetTopicsWithStats retrieves all predefined topics with their statistics from the latest analysis.
func (s *service) GetTopicsWithStats(ctx context.Context) ([]services.TopicStats, error) {
	logger := s.logger.WithSpan(ctx)
	logger.Info("getting topics with stats")

	// Get latest analysis
	latestAnalysis, err := s.analysisRepo.GetLatest(ctx)
	if err != nil {
		logger.RecordSpanError(ctx, err)
		logger.Error("error getting latest analysis", err)
		return nil, fmt.Errorf("failed to get latest analysis: %w", err)
	}

	if latestAnalysis == nil {
		logger.Info("no analysis found, returning empty topics")
		// Return all predefined topics with zero stats
		allTopics := analysis.AllTopics()
		stats := make([]services.TopicStats, len(allTopics))
		for i, topic := range allTopics {
			stats[i] = services.TopicStats{
				Topic:         topic,
				FeedbackCount: 0,
				AverageRating: 0,
			}
		}
		return stats, nil
	}

	// Get topics from latest analysis
	topics, err := s.analysisRepo.GetTopicsByAnalysisID(ctx, latestAnalysis.ID())
	if err != nil {
		logger.RecordSpanError(ctx, err)
		logger.Error("error getting topics", err, "analysis_id", latestAnalysis.ID())
		return nil, fmt.Errorf("failed to get topics: %w", err)
	}

	// Build map of topic enum -> topic analysis
	topicMap := make(map[analysis.Topic]*analysis.TopicAnalysis)
	for _, topic := range topics {
		topicMap[topic.Topic()] = topic
	}

	// Build stats for all predefined topics
	allTopics := analysis.AllTopics()
	stats := make([]services.TopicStats, len(allTopics))
	for i, topicEnum := range allTopics {
		stats[i] = services.TopicStats{
			Topic:         topicEnum,
			FeedbackCount: 0,
			AverageRating: 0,
		}

		// If this topic exists in the latest analysis, get its stats
		if topicAnalysis, exists := topicMap[topicEnum]; exists {
			stats[i].FeedbackCount = topicAnalysis.FeedbackCount()

			// Get feedback IDs for this topic and calculate average rating
			feedbackIDs, err := s.analysisRepo.GetFeedbackIDsByTopicID(ctx, topicAnalysis.ID())
			if err != nil {
				logger.Warning("error getting feedback IDs for topic", "topic_id", topicAnalysis.ID().String(), "error", err.Error())
				continue
			}

			if len(feedbackIDs) > 0 {
				// Get feedbacks and calculate average rating
				totalRating := 0
				validFeedbacks := 0
				for _, fbID := range feedbackIDs {
					fb, err := s.feedbackRepo.Get(ctx, fbID)
					if err != nil {
						logger.Warning("error getting feedback", "feedback_id", fbID.String(), "error", err.Error())
						continue
					}
					totalRating += fb.Rating().Value()
					validFeedbacks++
				}
				if validFeedbacks > 0 {
					stats[i].AverageRating = float64(totalRating) / float64(validFeedbacks)
				}
			}
		}
	}

	logger.Info("topics with stats retrieved", "topics_count", len(stats))
	return stats, nil
}

// GetTopicDetails retrieves details for a specific topic enum with all associated feedbacks.
func (s *service) GetTopicDetails(ctx context.Context, topicEnum analysis.Topic) (*services.TopicDetails, error) {
	logger := s.logger.WithSpan(ctx)
	logger.Info("getting topic details", "topic_enum", string(topicEnum))

	// Get latest analysis
	latestAnalysis, err := s.analysisRepo.GetLatest(ctx)
	if err != nil {
		logger.RecordSpanError(ctx, err)
		logger.Error("error getting latest analysis", err)
		return nil, fmt.Errorf("failed to get latest analysis: %w", err)
	}

	if latestAnalysis == nil {
		return nil, fmt.Errorf("no analysis found")
	}

	// Get topics from latest analysis
	topics, err := s.analysisRepo.GetTopicsByAnalysisID(ctx, latestAnalysis.ID())
	if err != nil {
		logger.RecordSpanError(ctx, err)
		logger.Error("error getting topics", err, "analysis_id", latestAnalysis.ID())
		return nil, fmt.Errorf("failed to get topics: %w", err)
	}

	// Find the topic analysis for this topic enum
	var topicAnalysis *analysis.TopicAnalysis
	for _, topic := range topics {
		if topic.Topic() == topicEnum {
			topicAnalysis = topic
			break
		}
	}

	if topicAnalysis == nil {
		return nil, fmt.Errorf("topic %s not found in latest analysis", string(topicEnum))
	}

	// Get feedback IDs for this topic
	feedbackIDs, err := s.analysisRepo.GetFeedbackIDsByTopicID(ctx, topicAnalysis.ID())
	if err != nil {
		logger.RecordSpanError(ctx, err)
		logger.Error("error getting feedback IDs", err, "topic_id", topicAnalysis.ID())
		return nil, fmt.Errorf("failed to get feedback IDs: %w", err)
	}

	// Get all feedbacks
	feedbacks := make([]*feedback.Feedback, 0, len(feedbackIDs))
	totalRating := 0
	for _, fbID := range feedbackIDs {
		fb, err := s.feedbackRepo.Get(ctx, fbID)
		if err != nil {
			logger.Warning("error getting feedback", "feedback_id", fbID.String(), "error", err.Error())
			continue
		}
		feedbacks = append(feedbacks, fb)
		totalRating += fb.Rating().Value()
	}

	averageRating := 0.0
	if len(feedbacks) > 0 {
		averageRating = float64(totalRating) / float64(len(feedbacks))
	}

	details := &services.TopicDetails{
		Topic:         topicEnum,
		Summary:       topicAnalysis.Summary(),
		FeedbackCount: len(feedbacks),
		AverageRating: averageRating,
		Sentiment:     topicAnalysis.Sentiment(),
		Feedbacks:     feedbacks,
	}

	logger.Info("topic details retrieved", "topic_enum", string(topicEnum), "feedbacks_count", len(feedbacks))
	return details, nil
}
