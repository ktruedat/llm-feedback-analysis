package analysis

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	apprepo "github.com/ktruedat/llm-feedback-analysis/internal/app/repository"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/repository/postgres/analysis/sqlc"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository/sql/utils"
)

func (r *repo) CreateAnalyzedFeedbacks(
	ctx context.Context,
	analysisID uuid.UUID,
	feedbackIDs []uuid.UUID,
	opts ...repository.RepoOption[apprepo.Options],
) error {
	q := utils.GetQuerier(opts, r.defaultQuerier)
	queries := newSQLCQueries(q)

	now := time.Now().UTC()
	for _, feedbackID := range feedbackIDs {
		err := queries.CreateAnalyzedFeedback(
			ctx, sqlc.CreateAnalyzedFeedbackParams{
				AnalysisID: analysisID,
				FeedbackID: feedbackID,
				CreatedAt:  now,
			},
		)
		if err != nil {
			return fmt.Errorf("failed to create analyzed feedback record for feedback %s: %w", feedbackID.String(), err)
		}
	}

	return nil
}

func (r *repo) GetFeedbackIDsByAnalysisID(
	ctx context.Context,
	analysisID uuid.UUID,
	opts ...repository.RepoOption[apprepo.Options],
) ([]uuid.UUID, error) {
	q := utils.GetQuerier(opts, r.defaultQuerier)
	queries := newSQLCQueries(q)

	feedbackIDs, err := queries.GetFeedbackIDsByAnalysisID(ctx, analysisID)
	if err != nil {
		return nil, fmt.Errorf("failed to get feedback IDs by analysis ID: %w", err)
	}

	return feedbackIDs, nil
}
