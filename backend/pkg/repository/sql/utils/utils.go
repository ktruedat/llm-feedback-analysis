package utils

import (
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository/sql/querier"
)

func GetQuerier[T any](opts []repository.RepoOption[T], defaultQuerier querier.PgxQuerier) querier.PgxQuerier {
	wrapper := BuildOpts(opts)

	return ExecutorAsQuerierOrDefault(wrapper.Ex, defaultQuerier)
}

func BuildOpts[T any](opts []repository.RepoOption[T]) repository.OptionsWrapper[T] {
	var wrapper repository.OptionsWrapper[T]
	for _, opt := range opts {
		opt(&wrapper)
	}

	return wrapper
}

func ExecutorAsQuerierOrDefault(ex repository.Executor, defaultQuerier querier.PgxQuerier) querier.PgxQuerier {
	if ex != nil {
		if pgxQuerier, ok := ex.(querier.PgxQuerier); ok {
			return pgxQuerier
		}
	}

	return defaultQuerier
}
