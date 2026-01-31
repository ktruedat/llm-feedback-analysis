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

func (r *repo) CreateTopicAnalysis(
	ctx context.Context,
	topicAnalysis *analysis.TopicAnalysis,
	opts ...repository.RepoOption[apprepo.Options],
) error {
	q := utils.GetQuerier(opts, r.defaultQuerier)
	queries := newSQLCQueries(q)

	if _, err := queries.CreateTopicAnalysis(
		ctx, sqlc.CreateTopicAnalysisParams{
			ID:            topicAnalysis.ID(),
			AnalysisID:    topicAnalysis.AnalysisID(),
			TopicEnum:     sqlc.FeedbackTopicEnum(topicAnalysis.Topic()),
			Summary:       topicAnalysis.Summary(),
			FeedbackCount: int32(topicAnalysis.FeedbackCount()),
			Sentiment:     sqlc.FeedbackSentiment(topicAnalysis.Sentiment()),
			CreatedAt:     topicAnalysis.CreatedAt(),
			UpdatedAt:     topicAnalysis.UpdatedAt(),
		},
	); err != nil {
		return fmt.Errorf("failed to create topic analysis: %w", err)
	}

	return nil
}
