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

func (r *repo) Update(
	ctx context.Context,
	id uuid.UUID,
	updates *analysis.UpdatableFields,
	opts ...repository.RepoOption[apprepo.Options],
) error {
	q := utils.GetQuerier(opts, r.defaultQuerier)
	queries := newSQLCQueries(q)

	// Get current analysis to use current values for fields we're not updating
	currentAnalysis, err := queries.GetAnalysisByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get current analysis: %w", err)
	}

	// If Results are present, update all LLM fields. Otherwise, keep current values.
	var overallSummary string
	var sentiment sqlc.FeedbackSentiment
	var keyInsights []string
	var tokens int32

	if updates.Results.IsSome() {
		// Success case: update all LLM fields from results
		results := updates.Results.Unwrap()
		overallSummary = results.OverallSummary
		sentiment = sqlc.FeedbackSentiment(results.Sentiment)
		keyInsights = results.KeyInsights
		tokens = int32(results.Tokens)
	} else {
		// Failure case: keep current LLM fields (they were set as placeholders during creation)
		overallSummary = currentAnalysis.OverallSummary
		sentiment = currentAnalysis.Sentiment
		keyInsights = currentAnalysis.KeyInsights
		tokens = currentAnalysis.Tokens
	}

	// Handle failure reason (only set if provided)
	var failureReason *string
	if updates.FailureReason.IsSome() {
		reason := updates.FailureReason.Unwrap()
		failureReason = &reason
	}

	// Handle completed_at
	var completedAt *time.Time
	if !updates.CompletedAt.IsZero() {
		completedAt = &updates.CompletedAt
	}

	err = queries.UpdateAnalysis(
		ctx, sqlc.UpdateAnalysisParams{
			ID:             id,
			OverallSummary: overallSummary,
			Sentiment:      sentiment,
			KeyInsights:    keyInsights,
			Tokens:         tokens,
			Status:         sqlc.FeedbackAnalysisStatus(updates.Status),
			FailureReason:  failureReason,
			CompletedAt:    completedAt,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update analysis: %w", err)
	}

	return nil
}
