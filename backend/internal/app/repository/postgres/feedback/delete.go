package feedback

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	apprepo "github.com/ktruedat/llm-feedback-analysis/internal/app/repository"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository/sql/utils"
)

func (r *repo) Delete(
	ctx context.Context,
	feedbackID uuid.UUID,
	opts ...repository.RepoOption[apprepo.Options],
) error {
	q := utils.GetQuerier(opts, r.defaultQuerier)
	queries := newSQLCQueries(q)

	rowsAffected, err := queries.DeleteFeedback(ctx, feedbackID)
	if err != nil {
		return fmt.Errorf("failed to delete feedback: %w", err)
	}

	// Check if any rows were affected
	if rowsAffected == 0 {
		return fmt.Errorf("feedback with ID %s not found or already deleted", feedbackID)
	}

	return nil
}
