package analysis

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ktruedat/llm-feedback-analysis/pkg/optional"
)

// Analysis represents an AI analysis snapshot of feedback data.
//
// Business Rules:
// - Must have a period (start and end timestamps)
// - Must have at least one feedback analyzed (tracked via analyzed_feedbacks junction table)
// - Status transitions: processing -> success/failed
// - Can reference a previous analysis for incremental updates
//
// Relationships:
// - Can reference a previous Analysis (one-to-one, optional)
// - Has many Feedbacks via AnalyzedFeedback junction table (many-to-many)
type Analysis struct {
	id                 uuid.UUID
	previousAnalysisID optional.Optional[uuid.UUID]
	periodStart        time.Time
	periodEnd          time.Time
	feedbackCount      int
	newFeedbackCount   optional.Optional[int]
	overallSummary     string
	sentiment          Sentiment
	keyInsights        []string
	model              string
	tokens             int
	analysisDurationMs int
	status             Status
	failureReason      optional.Optional[string]
	createdAt          time.Time
	completedAt        optional.Optional[time.Time]
}

// IsValid validates the entire analysis entity state.
func (a *Analysis) IsValid() error {
	if a.id == uuid.Nil {
		return fmt.Errorf("analysis ID is required")
	}

	if a.periodStart.IsZero() {
		return fmt.Errorf("period start timestamp is required")
	}

	if a.periodEnd.IsZero() {
		return fmt.Errorf("period end timestamp is required")
	}

	if a.periodEnd.Before(a.periodStart) {
		return fmt.Errorf("period end must be after period start")
	}

	if a.feedbackCount < 0 {
		return fmt.Errorf("feedback count cannot be negative")
	}

	if a.feedbackCount == 0 {
		return fmt.Errorf("at least one feedback must be analyzed")
	}

	if a.overallSummary == "" {
		return fmt.Errorf("overall summary is required")
	}

	if !a.sentiment.IsValid() {
		return fmt.Errorf("sentiment is required and must be valid")
	}

	if a.keyInsights == nil {
		return fmt.Errorf("key insights is required")
	}

	if a.model == "" {
		return fmt.Errorf("model is required")
	}

	if a.tokens < 0 {
		return fmt.Errorf("tokens cannot be negative")
	}

	if a.analysisDurationMs < 0 {
		return fmt.Errorf("analysis duration cannot be negative")
	}

	if a.createdAt.IsZero() {
		return fmt.Errorf("created_at timestamp is required")
	}

	return nil
}

// CanTransitionTo checks if the analysis can transition to the given status.
func (a *Analysis) CanTransitionTo(newStatus Status) bool {
	switch a.status {
	case StatusProcessing:
		return newStatus == StatusSuccess || newStatus == StatusFailed
	case StatusSuccess, StatusFailed:
		return false // Terminal states
	default:
		return false
	}
}

// MarkSuccess marks the analysis as successful with optional completion time.
func (a *Analysis) MarkSuccess() error {
	if !a.CanTransitionTo(StatusSuccess) {
		return fmt.Errorf("cannot transition from %s to success", a.status)
	}

	a.status = StatusSuccess
	a.completedAt = optional.Some(time.Now().UTC())
	return nil
}

// MarkFailed marks the analysis as failed with a failure reason.
func (a *Analysis) MarkFailed(reason string) error {
	if !a.CanTransitionTo(StatusFailed) {
		return fmt.Errorf("cannot transition from %s to failed", a.status)
	}

	a.status = StatusFailed
	a.failureReason = optional.Some(reason)
	a.completedAt = optional.Some(time.Now().UTC())
	return nil
}
