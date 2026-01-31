package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/transport/requests"
	"github.com/ktruedat/llm-feedback-analysis/internal/domain/analysis"
	"github.com/ktruedat/llm-feedback-analysis/internal/domain/feedback"
	"github.com/ktruedat/llm-feedback-analysis/internal/domain/user"
)

// FeedbackService defines the interface for feedback business logic operations.
type FeedbackService interface {
	// CreateFeedback creates a new feedback submission for the authenticated user.
	CreateFeedback(ctx context.Context, userID uuid.UUID, req *requests.CreateFeedbackRequest) (
		*feedback.Feedback,
		error,
	)

	// GetFeedbackByID retrieves a feedback entry by its ID.
	GetFeedbackByID(ctx context.Context, feedbackID uuid.UUID) (*feedback.Feedback, error)

	// ListFeedbacks retrieves a list of feedback entries with optional pagination.
	ListFeedbacks(ctx context.Context, limit, offset int) ([]*feedback.Feedback, error)

	// DeleteFeedback performs a soft delete on a feedback entry by its ID.
	DeleteFeedback(ctx context.Context, feedbackID uuid.UUID) error
}

// UserService defines the interface for user authentication and management operations.
type UserService interface {
	// RegisterUser creates a new user account and returns the created user.
	// Returns an error if the email already exists or validation fails.
	RegisterUser(ctx context.Context, req *requests.RegisterUserRequest) (*user.User, error)

	// AuthenticateUser authenticates a user with email and password.
	// Returns a JWT token string and user info if authentication succeeds.
	// Returns an error if credentials are invalid or user is inactive.
	AuthenticateUser(ctx context.Context, req *requests.LoginUserRequest) (string, *user.User, error)
}

// AnalyzerService defines the interface for LLM analysis operations.
type AnalyzerService interface {
	// EnqueueFeedback adds a feedback to the analysis queue.
	// This method is non-blocking and runs in a separate goroutine.
	EnqueueFeedback(ctx context.Context, fb *feedback.Feedback)

	// Start starts the analyzer service in a background goroutine.
	// It should be called once during application initialization.
	Start(ctx context.Context) error

	// Stop stops the analyzer service gracefully.
	Stop(ctx context.Context) error
}

// FeedbackSummaryService defines the interface for querying analysis data.
// This is separate from AnalyzerService which only performs the analysis.
type FeedbackSummaryService interface {
	// GetLatestAnalysis retrieves the latest completed analysis.
	GetLatestAnalysis(ctx context.Context) (*analysis.Analysis, error)

	// GetAllAnalyses retrieves all analyses ordered by creation date (newest first).
	GetAllAnalyses(ctx context.Context) ([]*analysis.Analysis, error)

	// GetAnalysisByID retrieves an analysis by ID with its topics and analyzed feedbacks with their topics.
	GetAnalysisByID(ctx context.Context, analysisID uuid.UUID) (
		*analysis.Analysis,
		[]*analysis.TopicAnalysis,
		map[uuid.UUID][]*analysis.TopicAnalysis, // feedback ID -> topics
		error,
	)
	// GetTopicsWithStats retrieves all predefined topics with their statistics from the latest analysis.
	// Returns topics with feedback count and average rating.
	GetTopicsWithStats(ctx context.Context) ([]TopicStats, error)
	// GetTopicDetails retrieves details for a specific topic enum with all associated feedbacks.
	GetTopicDetails(ctx context.Context, topicEnum analysis.Topic) (*TopicDetails, error)
}

// TopicStats represents statistics for a topic from the latest analysis.
type TopicStats struct {
	Topic         analysis.Topic
	FeedbackCount int
	AverageRating float64
}

// TopicDetails represents detailed information about a topic with all associated feedbacks.
type TopicDetails struct {
	Topic         analysis.Topic
	Summary       string
	FeedbackCount int
	AverageRating float64
	Sentiment     analysis.Sentiment
	Feedbacks     []*feedback.Feedback
}
