package feedback

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	apprepo "github.com/ktruedat/llm-feedback-analysis/internal/app/repository"
	"github.com/ktruedat/llm-feedback-analysis/internal/domain/feedback"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository/sql/utils"
)

func (r *repo) Get(
	ctx context.Context,
	feedbackID uuid.UUID,
	opts ...repository.RepoOption[apprepo.Options],
) (*feedback.Feedback, error) {
	q := utils.GetQuerier(opts, r.defaultQuerier)
	queries := newSQLCQueries(q)

	sqlcFeedback, err := queries.GetFeedback(ctx, feedbackID)
	if err != nil {
		return nil, fmt.Errorf("failed to get feedback: %w", err)
	}

	return mapSQLCFeedbackToDomain(sqlcFeedback), nil
}
