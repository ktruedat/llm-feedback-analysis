package user

import (
	"github.com/ktruedat/llm-feedback-analysis/internal/app/repository"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/repository/postgres/user/sqlc"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository/sql/querier"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository/sql/utils"
)

type repo struct {
	defaultQuerier querier.PgxQuerier
}

// NewUserRepository creates a new user repository.
func NewUserRepository(q querier.PgxQuerier) repository.UserRepository {
	return &repo{
		defaultQuerier: q,
	}
}

var _ sqlc.DBTX = (*utils.QuerierAdapter)(nil)

func newSQLCQueries(q querier.PgxQuerier) *sqlc.Queries {
	return sqlc.New(utils.NewQuerierAdapter(q))
}
