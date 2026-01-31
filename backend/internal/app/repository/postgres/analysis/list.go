package analysis

import (
	"context"
	"fmt"

	apprepo "github.com/ktruedat/llm-feedback-analysis/internal/app/repository"
	"github.com/ktruedat/llm-feedback-analysis/internal/domain/analysis"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository/sql/utils"
)

func (r *repo) List(
	ctx context.Context,
	opts ...repository.RepoOption[apprepo.Options],
) ([]*analysis.Analysis, error) {
	q := utils.GetQuerier(opts, r.defaultQuerier)
	queries := newSQLCQueries(q)

	sqlcAnalyses, err := queries.ListAnalyses(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list analyses: %w", err)
	}

	analyses := make([]*analysis.Analysis, len(sqlcAnalyses))
	for i, sqlcAnalysis := range sqlcAnalyses {
		analyses[i] = mapSQLCAnalysisToDomain(sqlcAnalysis)
	}

	return analyses, nil
}
