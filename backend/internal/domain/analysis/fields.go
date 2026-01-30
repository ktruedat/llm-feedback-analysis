package analysis

import (
	"time"

	"github.com/google/uuid"
	"github.com/ktruedat/llm-feedback-analysis/pkg/optional"
)

// ID returns the analysis ID.
func (a *Analysis) ID() uuid.UUID {
	return a.id
}

// PreviousAnalysisID returns the previous analysis ID if set.
func (a *Analysis) PreviousAnalysisID() optional.Optional[uuid.UUID] {
	return a.previousAnalysisID
}

// PeriodStart returns the start timestamp of the analysis period.
func (a *Analysis) PeriodStart() time.Time {
	return a.periodStart
}

// PeriodEnd returns the end timestamp of the analysis period.
func (a *Analysis) PeriodEnd() time.Time {
	return a.periodEnd
}

// FeedbackCount returns the total number of feedbacks in this analysis.
func (a *Analysis) FeedbackCount() int {
	return a.feedbackCount
}

// NewFeedbackCount returns the number of new feedbacks since the previous analysis.
func (a *Analysis) NewFeedbackCount() optional.Optional[int] {
	return a.newFeedbackCount
}

// OverallSummary returns the human-readable summary of all feedback.
func (a *Analysis) OverallSummary() string {
	return a.overallSummary
}

// Sentiment returns the overall sentiment.
func (a *Analysis) Sentiment() Sentiment {
	return a.sentiment
}

// KeyInsights returns the array of key insights/takeaways.
func (a *Analysis) KeyInsights() []string {
	return a.keyInsights
}

// Model returns the LLM model used for this analysis.
func (a *Analysis) Model() string {
	return a.model
}

// Tokens returns the total tokens consumed.
func (a *Analysis) Tokens() int {
	return a.tokens
}

// AnalysisDurationMs returns the analysis duration in milliseconds.
func (a *Analysis) AnalysisDurationMs() int {
	return a.analysisDurationMs
}

// Status returns the analysis status.
func (a *Analysis) Status() Status {
	return a.status
}

// FailureReason returns the failure reason if the analysis failed.
func (a *Analysis) FailureReason() optional.Optional[string] {
	return a.failureReason
}

// CreatedAt returns the creation timestamp.
func (a *Analysis) CreatedAt() time.Time {
	return a.createdAt
}

// CompletedAt returns the completion timestamp if completed.
func (a *Analysis) CompletedAt() optional.Optional[time.Time] {
	return a.completedAt
}
