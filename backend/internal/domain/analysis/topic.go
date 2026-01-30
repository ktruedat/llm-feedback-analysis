package analysis

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Topic represents a topic/theme identified by AI analysis.
//
// Business Rules:
// - Must belong to an analysis
// - Must have a name (short, descriptive)
// - Must have a description
// - Must have at least one feedback assigned
// - Sentiment must be valid
// - Importance score must be between 0.0 and 1.0
//
// Relationships:
// - Belongs to Analysis (many-to-one)
// - Has many Feedbacks via FeedbackTopicAssignment (many-to-many)
type Topic struct {
	id            uuid.UUID
	analysisID    uuid.UUID
	topicName     string
	description   string
	feedbackCount int
	sentiment     Sentiment
	createdAt     time.Time
	updatedAt     time.Time
}

// IsValid validates the entire topic entity state.
func (t *Topic) IsValid() error {
	if t.id == uuid.Nil {
		return fmt.Errorf("topic ID is required")
	}

	if t.analysisID == uuid.Nil {
		return fmt.Errorf("analysis ID is required")
	}

	if t.topicName == "" {
		return fmt.Errorf("topic name is required")
	}

	if len(t.topicName) > 100 {
		return fmt.Errorf("topic name cannot exceed 100 characters")
	}

	if t.description == "" {
		return fmt.Errorf("description is required")
	}

	if t.feedbackCount < 0 {
		return fmt.Errorf("feedback count cannot be negative")
	}

	if t.feedbackCount == 0 {
		return fmt.Errorf("topic must have at least one feedback")
	}

	if !t.sentiment.IsValid() {
		return fmt.Errorf("sentiment is required and must be valid")
	}

	if t.createdAt.IsZero() {
		return fmt.Errorf("created_at timestamp is required")
	}

	if t.updatedAt.IsZero() {
		return fmt.Errorf("updated_at timestamp is required")
	}

	return nil
}
