package feedback

import (
	"context"
	"fmt"
	"time"

	apprepo "github.com/ktruedat/llm-feedback-analysis/internal/app/repository"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/repository/postgres/feedback/sqlc"
	"github.com/ktruedat/llm-feedback-analysis/internal/domain/feedback"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository/sql/utils"
)

func (r *repo) Create(
	ctx context.Context,
	fb *feedback.Feedback,
	opts ...repository.RepoOption[apprepo.Options],
) error {
	q := utils.GetQuerier(opts, r.defaultQuerier)
	queries := newSQLCQueries(q)

	var deletedAt *time.Time
	if fb.DeletedAt().IsSome() {
		dt := fb.DeletedAt().Unwrap()
		deletedAt = &dt
	}

	_, err := queries.CreateFeedback(
		ctx, sqlc.CreateFeedbackParams{
			ID:        fb.ID(),
			Rating:    int32(fb.Rating().Value()),
			Comment:   fb.Comment().Value(),
			CreatedAt: fb.CreatedAt(),
			UpdatedAt: fb.UpdatedAt(),
			DeletedAt: deletedAt,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create feedback: %w", err)
	}

	return nil
}
