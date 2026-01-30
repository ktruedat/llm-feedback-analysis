package analysis

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ktruedat/llm-feedback-analysis/pkg/optional"
)

// Builder provides type-safe, fluent construction of Analysis entities.
type Builder struct {
	entity           *Analysis
	validationErrors []error
}

// NewBuilder creates a builder for creating new analysis entities.
func NewBuilder() *Builder {
	now := time.Now().UTC()
	return &Builder{
		entity: &Analysis{
			id:                 uuid.New(),
			status:             StatusProcessing,
			createdAt:          now,
			keyInsights:        []string{},
			sentiment:          SentimentMixed, // Default sentiment
			tokens:             0,              // Must be set explicitly
			analysisDurationMs: 0,              // Must be set explicitly
		},
		validationErrors: make([]error, 0),
	}
}

// BuilderFromExisting creates a builder from an existing analysis entity.
func BuilderFromExisting(a *Analysis) *Builder {
	copied := *a
	return &Builder{
		entity:           &copied,
		validationErrors: make([]error, 0),
	}
}

// WithID sets the analysis ID.
func (b *Builder) WithID(id uuid.UUID) *Builder {
	if id == uuid.Nil {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("analysis ID cannot be nil"))
		return b
	}
	b.entity.id = id
	return b
}

// WithPreviousAnalysisID sets the previous analysis ID.
func (b *Builder) WithPreviousAnalysisID(id uuid.UUID) *Builder {
	if id == uuid.Nil {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("previous analysis ID cannot be nil"))
		return b
	}
	b.entity.previousAnalysisID = optional.Some(id)
	return b
}

// WithPeriod sets the analysis period.
func (b *Builder) WithPeriod(start, end time.Time) *Builder {
	if start.IsZero() {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("period start cannot be zero"))
		return b
	}
	if end.IsZero() {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("period end cannot be zero"))
		return b
	}
	if end.Before(start) {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("period end must be after period start"))
		return b
	}
	b.entity.periodStart = start
	b.entity.periodEnd = end
	return b
}

// WithFeedbackCount sets the feedback count.
func (b *Builder) WithFeedbackCount(count int) *Builder {
	if count < 0 {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("feedback count cannot be negative"))
		return b
	}
	b.entity.feedbackCount = count
	return b
}

// WithNewFeedbackCount sets the new feedback count.
func (b *Builder) WithNewFeedbackCount(count int) *Builder {
	if count < 0 {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("new feedback count cannot be negative"))
		return b
	}
	b.entity.newFeedbackCount = optional.Some(count)
	return b
}

// WithOverallSummary sets the overall summary.
func (b *Builder) WithOverallSummary(summary string) *Builder {
	if summary == "" {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("overall summary cannot be empty"))
		return b
	}
	b.entity.overallSummary = summary
	return b
}

// WithSentiment sets the sentiment.
func (b *Builder) WithSentiment(sentiment Sentiment) *Builder {
	if !sentiment.IsValid() {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("invalid sentiment: %s", sentiment))
		return b
	}
	b.entity.sentiment = sentiment
	return b
}

// WithKeyInsights sets the key insights.
func (b *Builder) WithKeyInsights(insights []string) *Builder {
	if insights == nil {
		b.entity.keyInsights = []string{}
	} else {
		b.entity.keyInsights = insights
	}
	return b
}

// WithModel sets the model name.
func (b *Builder) WithModel(model string) *Builder {
	if model == "" {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("model cannot be empty"))
		return b
	}
	b.entity.model = model
	return b
}

// WithTokens sets the tokens consumed.
func (b *Builder) WithTokens(tokens int) *Builder {
	if tokens < 0 {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("tokens cannot be negative"))
		return b
	}
	b.entity.tokens = tokens
	return b
}

// WithAnalysisDurationMs sets the analysis duration in milliseconds.
func (b *Builder) WithAnalysisDurationMs(ms int) *Builder {
	if ms < 0 {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("analysis duration cannot be negative"))
		return b
	}
	b.entity.analysisDurationMs = ms
	return b
}

// WithStatus sets the status.
func (b *Builder) WithStatus(status Status) *Builder {
	if !status.IsValid() {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("invalid status: %s", status))
		return b
	}
	b.entity.status = status
	return b
}

// WithFailureReason sets the failure reason.
func (b *Builder) WithFailureReason(reason string) *Builder {
	if reason != "" {
		b.entity.failureReason = optional.Some(reason)
	}
	return b
}

// WithCreatedAt sets the creation timestamp.
func (b *Builder) WithCreatedAt(t time.Time) *Builder {
	if t.IsZero() {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("created_at cannot be zero"))
		return b
	}
	b.entity.createdAt = t
	return b
}

// WithCompletedAt sets the completion timestamp.
func (b *Builder) WithCompletedAt(t time.Time) *Builder {
	if !t.IsZero() {
		b.entity.completedAt = optional.Some(t)
	}
	return b
}

// Build validates all accumulated data and returns the analysis entity.
func (b *Builder) Build() (*Analysis, error) {
	// Return accumulated validation errors first
	if len(b.validationErrors) > 0 {
		return nil, b.validationErrors[0]
	}

	// Validate the entity itself
	if err := b.entity.IsValid(); err != nil {
		return nil, fmt.Errorf("analysis validation failed: %w", err)
	}

	return b.entity, nil
}

// BuildUnchecked returns the analysis entity without validation.
// ONLY use this for database reconstruction where data integrity is already guaranteed.
func (b *Builder) BuildUnchecked() *Analysis {
	return b.entity
}
