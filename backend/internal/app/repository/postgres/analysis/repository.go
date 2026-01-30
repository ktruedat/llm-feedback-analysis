package analysis

import (
	"github.com/ktruedat/llm-feedback-analysis/internal/app/repository"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/repository/postgres/analysis/sqlc"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository/sql/querier"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository/sql/utils"
)

type repo struct {
	defaultQuerier querier.PgxQuerier
}

// NewAnalysisRepository creates a new analysis repository.
func NewAnalysisRepository(q querier.PgxQuerier) repository.AnalysisRepository {
	return &repo{
		defaultQuerier: q,
	}
}

var _ sqlc.DBTX = (*utils.QuerierAdapter)(nil)

func newSQLCQueries(q querier.PgxQuerier) *sqlc.Queries {
	return sqlc.New(utils.NewQuerierAdapter(q))
}
