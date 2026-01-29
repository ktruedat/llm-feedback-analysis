package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/transport/requests"
	"github.com/ktruedat/llm-feedback-analysis/internal/domain/feedback"
)

// FeedbackService defines the interface for feedback business logic operations.
type FeedbackService interface {
	// CreateFeedback creates a new feedback submission.
	CreateFeedback(ctx context.Context, req *requests.CreateFeedbackRequest) (*feedback.Feedback, error)

	// GetFeedbackByID retrieves a feedback entry by its ID.
	GetFeedbackByID(ctx context.Context, feedbackID uuid.UUID) (*feedback.Feedback, error)

	// ListFeedbacks retrieves a list of feedback entries with optional pagination.
	ListFeedbacks(ctx context.Context, limit, offset int) ([]*feedback.Feedback, error)

	// DeleteFeedback performs a soft delete on a feedback entry by its ID.
	DeleteFeedback(ctx context.Context, feedbackID uuid.UUID) error
}
