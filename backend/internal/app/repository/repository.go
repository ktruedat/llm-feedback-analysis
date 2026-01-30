package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ktruedat/llm-feedback-analysis/internal/domain/feedback"
	"github.com/ktruedat/llm-feedback-analysis/internal/domain/user"
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

type UserRepository interface {
	// Create stores a new user in the repository.
	Create(ctx context.Context, u *user.User, opts ...repository.RepoOption[Options]) error
	// GetByID retrieves a user by its ID.
	GetByID(ctx context.Context, userID uuid.UUID, opts ...repository.RepoOption[Options]) (*user.User, error)
	// GetByEmail retrieves a user by email address.
	GetByEmail(ctx context.Context, email string, opts ...repository.RepoOption[Options]) (*user.User, error)
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
