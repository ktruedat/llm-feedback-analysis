package analysis

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	apprepo "github.com/ktruedat/llm-feedback-analysis/internal/app/repository"
	"github.com/ktruedat/llm-feedback-analysis/internal/domain/analysis"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository/sql/utils"
)

func (r *repo) GetTopicsByAnalysisID(
	ctx context.Context,
	analysisID uuid.UUID,
	opts ...repository.RepoOption[apprepo.Options],
) ([]*analysis.TopicAnalysis, error) {
	q := utils.GetQuerier(opts, r.defaultQuerier)
	queries := newSQLCQueries(q)

	sqlcTopics, err := queries.GetTopicsByAnalysisID(ctx, analysisID)
	if err != nil {
		return nil, fmt.Errorf("failed to get topic analyses by analysis ID: %w", err)
	}

	topicAnalyses := make([]*analysis.TopicAnalysis, len(sqlcTopics))
	for i, sqlcTopic := range sqlcTopics {
		topicAnalyses[i] = mapSQLCTopicToDomain(sqlcTopic)
	}

	return topicAnalyses, nil
}

func (r *repo) GetFeedbackIDsByTopicID(
	ctx context.Context,
	topicID uuid.UUID,
	opts ...repository.RepoOption[apprepo.Options],
) ([]uuid.UUID, error) {
	q := utils.GetQuerier(opts, r.defaultQuerier)
	queries := newSQLCQueries(q)

	feedbackIDs, err := queries.GetFeedbackIDsByTopicID(ctx, topicID)
	if err != nil {
		return nil, fmt.Errorf("failed to get feedback IDs by topic ID: %w", err)
	}

	return feedbackIDs, nil
}
