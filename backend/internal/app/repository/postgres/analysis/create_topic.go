package analysis

import (
	"context"
	"fmt"

	apprepo "github.com/ktruedat/llm-feedback-analysis/internal/app/repository"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/repository/postgres/analysis/sqlc"
	"github.com/ktruedat/llm-feedback-analysis/internal/domain/analysis"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository/sql/utils"
)

func (r *repo) CreateTopic(
	ctx context.Context,
	topic *analysis.Topic,
	opts ...repository.RepoOption[apprepo.Options],
) error {
	q := utils.GetQuerier(opts, r.defaultQuerier)
	queries := newSQLCQueries(q)

	_, err := queries.CreateTopic(
		ctx, sqlc.CreateTopicParams{
			ID:            topic.ID(),
			AnalysisID:    topic.AnalysisID(),
			TopicName:     topic.TopicName(),
			Description:   topic.Description(),
			FeedbackCount: int32(topic.FeedbackCount()),
			Sentiment:     sqlc.FeedbackSentiment(topic.Sentiment()),
			CreatedAt:     topic.CreatedAt(),
			UpdatedAt:     topic.UpdatedAt(),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create topic: %w", err)
	}

	return nil
}
