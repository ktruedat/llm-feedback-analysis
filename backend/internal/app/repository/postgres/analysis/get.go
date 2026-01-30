package analysis

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	apprepo "github.com/ktruedat/llm-feedback-analysis/internal/app/repository"
	"github.com/ktruedat/llm-feedback-analysis/internal/domain/analysis"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository/sql/utils"
)

func (r *repo) GetByID(
	ctx context.Context,
	analysisID uuid.UUID,
	opts ...repository.RepoOption[apprepo.Options],
) (*analysis.Analysis, error) {
	q := utils.GetQuerier(opts, r.defaultQuerier)
	queries := newSQLCQueries(q)

	sqlcAnalysis, err := queries.GetAnalysisByID(ctx, analysisID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("analysis not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get analysis: %w", err)
	}

	return mapSQLCAnalysisToDomain(sqlcAnalysis), nil
}

func (r *repo) GetLatest(
	ctx context.Context,
	opts ...repository.RepoOption[apprepo.Options],
) (*analysis.Analysis, error) {
	q := utils.GetQuerier(opts, r.defaultQuerier)
	queries := newSQLCQueries(q)

	sqlcAnalysis, err := queries.GetLatestAnalysis(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No previous analysis exists
		}
		return nil, fmt.Errorf("failed to get latest analysis: %w", err)
	}

	return mapSQLCAnalysisToDomain(sqlcAnalysis), nil
}
