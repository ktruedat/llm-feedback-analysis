package analysis

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// TopicAnalysisBuilder provides type-safe, fluent construction of TopicAnalysis entities.
type TopicAnalysisBuilder struct {
	entity           *TopicAnalysis
	validationErrors []error
}

// NewTopicAnalysisBuilder creates a builder for creating new topic analysis entities.
func NewTopicAnalysisBuilder() *TopicAnalysisBuilder {
	now := time.Now().UTC()
	return &TopicAnalysisBuilder{
		entity: &TopicAnalysis{
			id:        uuid.New(),
			createdAt: now,
			updatedAt: now,
		},
		validationErrors: make([]error, 0),
	}
}

// BuilderFromExistingTopicAnalysis creates a builder from an existing topic analysis entity.
func BuilderFromExistingTopicAnalysis(t *TopicAnalysis) *TopicAnalysisBuilder {
	copied := *t
	copied.updatedAt = time.Now().UTC()
	return &TopicAnalysisBuilder{
		entity:           &copied,
		validationErrors: make([]error, 0),
	}
}

// WithID sets the topic analysis ID.
func (b *TopicAnalysisBuilder) WithID(id uuid.UUID) *TopicAnalysisBuilder {
	if id == uuid.Nil {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("topic analysis ID cannot be nil"))
		return b
	}
	b.entity.id = id
	return b
}

// WithAnalysisID sets the analysis ID.
func (b *TopicAnalysisBuilder) WithAnalysisID(analysisID uuid.UUID) *TopicAnalysisBuilder {
	if analysisID == uuid.Nil {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("analysis ID cannot be nil"))
		return b
	}
	b.entity.analysisID = analysisID
	return b
}

// WithTopic sets the topic enum.
func (b *TopicAnalysisBuilder) WithTopic(topic Topic) *TopicAnalysisBuilder {
	if !topic.IsValid() {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("topic must be valid"))
		return b
	}
	b.entity.topic = topic
	return b
}

// WithSummary sets the topic analysis summary.
func (b *TopicAnalysisBuilder) WithSummary(summary string) *TopicAnalysisBuilder {
	if summary == "" {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("summary cannot be empty"))
		return b
	}
	b.entity.summary = summary
	return b
}

// WithFeedbackCount sets the feedback count.
func (b *TopicAnalysisBuilder) WithFeedbackCount(count int) *TopicAnalysisBuilder {
	if count < 0 {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("feedback count cannot be negative"))
		return b
	}
	b.entity.feedbackCount = count
	return b
}

// WithSentiment sets the sentiment.
func (b *TopicAnalysisBuilder) WithSentiment(sentiment Sentiment) *TopicAnalysisBuilder {
	if !sentiment.IsValid() {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("invalid sentiment: %s", sentiment))
		return b
	}
	b.entity.sentiment = sentiment
	return b
}

// WithCreatedAt sets the creation timestamp.
func (b *TopicAnalysisBuilder) WithCreatedAt(t time.Time) *TopicAnalysisBuilder {
	if t.IsZero() {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("created_at cannot be zero"))
		return b
	}
	b.entity.createdAt = t
	return b
}

// WithUpdatedAt sets the update timestamp.
func (b *TopicAnalysisBuilder) WithUpdatedAt(t time.Time) *TopicAnalysisBuilder {
	if t.IsZero() {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("updated_at cannot be zero"))
		return b
	}
	b.entity.updatedAt = t
	return b
}

// Build validates all accumulated data and returns the topic analysis entity.
func (b *TopicAnalysisBuilder) Build() (*TopicAnalysis, error) {
	// Return accumulated validation errors first
	if len(b.validationErrors) > 0 {
		return nil, b.validationErrors[0]
	}

	// Validate the entity itself
	if err := b.entity.IsValid(); err != nil {
		return nil, fmt.Errorf("topic analysis validation failed: %w", err)
	}

	return b.entity, nil
}

// BuildUnchecked returns the topic analysis entity without validation.
// ONLY use this for database reconstruction where data integrity is already guaranteed.
func (b *TopicAnalysisBuilder) BuildUnchecked() *TopicAnalysis {
	return b.entity
}
