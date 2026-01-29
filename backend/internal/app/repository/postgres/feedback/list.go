package feedback

import (
	"context"
	"fmt"

	apprepo "github.com/ktruedat/llm-feedback-analysis/internal/app/repository"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/repository/postgres/feedback/sqlc"
	"github.com/ktruedat/llm-feedback-analysis/internal/domain/feedback"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository/sql/utils"
)

func (r *repo) List(
	ctx context.Context,
	opts ...repository.RepoOption[apprepo.Options],
) ([]*feedback.Feedback, error) {
	q := utils.GetQuerier(opts, r.defaultQuerier)
	queries := newSQLCQueries(q)

	wrapper := utils.BuildOpts(opts)
	options := wrapper.Ext

	var limit *int32
	var offset *int32
	if options != nil {
		if options.Limit > 0 {
			l := int32(options.Limit)
			limit = &l
		}
		if options.Offset > 0 {
			o := int32(options.Offset)
			offset = &o
		}
	}

	// Default limit if not specified
	if limit == nil {
		defaultLimit := int32(100)
		limit = &defaultLimit
	}
	if offset == nil {
		defaultOffset := int32(0)
		offset = &defaultOffset
	}

	var sqlcFeedbacks []sqlc.Feedback
	sqlcFeedbacks, err := queries.ListFeedbacks(ctx, *limit, *offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list feedbacks: %w", err)
	}

	// Map to domain
	feedbacks := make([]*feedback.Feedback, len(sqlcFeedbacks))
	for i, sqlcFeedback := range sqlcFeedbacks {
		feedbacks[i] = mapSQLCFeedbackToDomain(sqlcFeedback)
	}

	return feedbacks, nil
}
