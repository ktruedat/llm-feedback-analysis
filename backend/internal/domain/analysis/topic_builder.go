package analysis

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// TopicBuilder provides type-safe, fluent construction of Topic entities.
type TopicBuilder struct {
	entity           *Topic
	validationErrors []error
}

// NewTopicBuilder creates a builder for creating new topic entities.
func NewTopicBuilder() *TopicBuilder {
	now := time.Now().UTC()
	return &TopicBuilder{
		entity: &Topic{
			id:        uuid.New(),
			createdAt: now,
			updatedAt: now,
		},
		validationErrors: make([]error, 0),
	}
}

// BuilderFromExistingTopic creates a builder from an existing topic entity.
func BuilderFromExistingTopic(t *Topic) *TopicBuilder {
	copied := *t
	copied.updatedAt = time.Now().UTC()
	return &TopicBuilder{
		entity:           &copied,
		validationErrors: make([]error, 0),
	}
}

// WithID sets the topic ID.
func (b *TopicBuilder) WithID(id uuid.UUID) *TopicBuilder {
	if id == uuid.Nil {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("topic ID cannot be nil"))
		return b
	}
	b.entity.id = id
	return b
}

// WithAnalysisID sets the analysis ID.
func (b *TopicBuilder) WithAnalysisID(analysisID uuid.UUID) *TopicBuilder {
	if analysisID == uuid.Nil {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("analysis ID cannot be nil"))
		return b
	}
	b.entity.analysisID = analysisID
	return b
}

// WithTopicName sets the topic name.
func (b *TopicBuilder) WithTopicName(name string) *TopicBuilder {
	if name == "" {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("topic name cannot be empty"))
		return b
	}
	if len(name) > 100 {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("topic name cannot exceed 100 characters"))
		return b
	}
	b.entity.topicName = name
	return b
}

// WithDescription sets the topic description.
func (b *TopicBuilder) WithDescription(description string) *TopicBuilder {
	if description == "" {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("description cannot be empty"))
		return b
	}
	b.entity.description = description
	return b
}

// WithFeedbackCount sets the feedback count.
func (b *TopicBuilder) WithFeedbackCount(count int) *TopicBuilder {
	if count < 0 {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("feedback count cannot be negative"))
		return b
	}
	b.entity.feedbackCount = count
	return b
}

// WithSentiment sets the sentiment.
func (b *TopicBuilder) WithSentiment(sentiment Sentiment) *TopicBuilder {
	if !sentiment.IsValid() {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("invalid sentiment: %s", sentiment))
		return b
	}
	b.entity.sentiment = sentiment
	return b
}

// WithCreatedAt sets the creation timestamp.
func (b *TopicBuilder) WithCreatedAt(t time.Time) *TopicBuilder {
	if t.IsZero() {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("created_at cannot be zero"))
		return b
	}
	b.entity.createdAt = t
	return b
}

// WithUpdatedAt sets the update timestamp.
func (b *TopicBuilder) WithUpdatedAt(t time.Time) *TopicBuilder {
	if t.IsZero() {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("updated_at cannot be zero"))
		return b
	}
	b.entity.updatedAt = t
	return b
}

// Build validates all accumulated data and returns the topic entity.
func (b *TopicBuilder) Build() (*Topic, error) {
	// Return accumulated validation errors first
	if len(b.validationErrors) > 0 {
		return nil, b.validationErrors[0]
	}

	// Validate the entity itself
	if err := b.entity.IsValid(); err != nil {
		return nil, fmt.Errorf("topic validation failed: %w", err)
	}

	return b.entity, nil
}

// BuildUnchecked returns the topic entity without validation.
// ONLY use this for database reconstruction where data integrity is already guaranteed.
func (b *TopicBuilder) BuildUnchecked() *Topic {
	return b.entity
}
