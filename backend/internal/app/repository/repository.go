package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ktruedat/llm-feedback-analysis/internal/domain/feedback"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository"
)

type FeedbackRepository interface {
	// Create stores a new feedback entry in the repository.
	Create(ctx context.Context, feedback *feedback.Feedback, opts ...repository.RepoOption[Options]) error
	// Get retrieves a feedback entry by its ID.
	Get(ctx context.Context, feedbackID uuid.UUID, opts ...repository.RepoOption[Options]) (*feedback.Feedback, error)
	// List retrieves a list of feedback entries from the repository.
	List(ctx context.Context, opts ...repository.RepoOption[Options]) ([]*feedback.Feedback, error)
	// Delete performs a soft delete on a feedback entry by setting deleted_at timestamp.
	Delete(ctx context.Context, feedbackID uuid.UUID, opts ...repository.RepoOption[Options]) error
}

type Options struct {
	Limit  int
	Offset int
}

func WithOptions(opts *Options) repository.RepoOption[Options] {
	return func(o *repository.OptionsWrapper[Options]) {
		o.Ext = opts
	}
}
