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

func (r *repo) CreateTopicAssignments(
	ctx context.Context,
	analysisID uuid.UUID,
	topicID uuid.UUID,
	feedbackIDs []uuid.UUID,
	opts ...repository.RepoOption[apprepo.Options],
) error {
	q := utils.GetQuerier(opts, r.defaultQuerier)
	queries := newSQLCQueries(q)

	now := time.Now().UTC()
	for _, feedbackID := range feedbackIDs {
		err := queries.CreateTopicAssignment(
			ctx, sqlc.CreateTopicAssignmentParams{
				ID:         uuid.New(),
				AnalysisID: analysisID,
				FeedbackID: feedbackID,
				TopicID:    topicID,
				CreatedAt:  now,
			},
		)
		if err != nil {
			return fmt.Errorf("failed to create topic assignment for feedback %s: %w", feedbackID.String(), err)
		}
	}

	return nil
}
