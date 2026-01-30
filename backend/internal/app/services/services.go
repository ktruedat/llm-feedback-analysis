package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/transport/requests"
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
