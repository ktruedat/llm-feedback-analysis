package feedback

import (
	"github.com/ktruedat/llm-feedback-analysis/internal/app/repository"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/repository/postgres/feedback/sqlc"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository/sql/querier"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository/sql/utils"
)

type repo struct {
	defaultQuerier querier.PgxQuerier
}

// NewFeedbackRepository creates a new feedback repository.
func NewFeedbackRepository(q querier.PgxQuerier) repository.FeedbackRepository {
	return &repo{
		defaultQuerier: q,
	}
}

var _ sqlc.DBTX = (*utils.QuerierAdapter)(nil)

func newSQLCQueries(q querier.PgxQuerier) *sqlc.Queries {
	return sqlc.New(utils.NewQuerierAdapter(q))
}
