package analysis

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	apprepo "github.com/ktruedat/llm-feedback-analysis/internal/app/repository"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/repository/postgres/analysis/sqlc"
	"github.com/ktruedat/llm-feedback-analysis/internal/domain/analysis"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository/sql/utils"
)

func (r *repo) Create(
	ctx context.Context,
	a *analysis.Analysis,
	opts ...repository.RepoOption[apprepo.Options],
) error {
	q := utils.GetQuerier(opts, r.defaultQuerier)
	queries := newSQLCQueries(q)

	var previousAnalysisID *uuid.UUID
	if a.PreviousAnalysisID().IsSome() {
		id := a.PreviousAnalysisID().Unwrap()
		previousAnalysisID = &id
	}

	var newFeedbackCount *int32
	if a.NewFeedbackCount().IsSome() {
		count := int32(a.NewFeedbackCount().Unwrap())
		newFeedbackCount = &count
	}

	var failureReason *string
	if a.FailureReason().IsSome() {
		reason := a.FailureReason().Unwrap()
		failureReason = &reason
	}

	var completedAt *time.Time
	if a.CompletedAt().IsSome() {
		completed := a.CompletedAt().Unwrap()
		completedAt = &completed
	}

	_, err := queries.CreateAnalysis(
		ctx, sqlc.CreateAnalysisParams{
			ID:                 a.ID(),
			PreviousAnalysisID: previousAnalysisID,
			PeriodStart:        a.PeriodStart(),
			PeriodEnd:          a.PeriodEnd(),
			FeedbackCount:      int32(a.FeedbackCount()),
			NewFeedbackCount:   newFeedbackCount,
			OverallSummary:     a.OverallSummary(),
			Sentiment:          sqlc.FeedbackSentiment(a.Sentiment()),
			KeyInsights:        a.KeyInsights(),
			Model:              a.Model(),
			Tokens:             int32(a.Tokens()),
			AnalysisDurationMs: int32(a.AnalysisDurationMs()),
			Status:             sqlc.FeedbackAnalysisStatus(a.Status()),
			FailureReason:      failureReason,
			CreatedAt:          a.CreatedAt(),
			CompletedAt:        completedAt,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create analysis: %w", err)
	}

	return nil
}
