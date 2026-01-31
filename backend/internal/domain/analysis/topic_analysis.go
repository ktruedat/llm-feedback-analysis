package analysis

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// TopicAnalysis represents the analysis of a specific topic within an analysis.
//
// Business Rules:
// - Must belong to an analysis
// - Must have a valid topic (predefined business topic enum)
// - Must have a summary (analysis summary for this topic)
// - Must have at least one feedback assigned
// - Sentiment must be valid
//
// Relationships:
// - Belongs to Analysis (many-to-one)
// - Has many Feedbacks via FeedbackTopicAssignment (many-to-many)
type TopicAnalysis struct {
	id            uuid.UUID
	analysisID    uuid.UUID
	topic         Topic
	summary       string
	feedbackCount int
	sentiment     Sentiment
	createdAt     time.Time
	updatedAt     time.Time
}

// IsValid validates the entire topic analysis entity state.
func (t *TopicAnalysis) IsValid() error {
	if t.id == uuid.Nil {
		return fmt.Errorf("topic analysis ID is required")
	}

	if t.analysisID == uuid.Nil {
		return fmt.Errorf("analysis ID is required")
	}

	if !t.topic.IsValid() {
		return fmt.Errorf("topic is required and must be valid")
	}

	if t.summary == "" {
		return fmt.Errorf("summary is required")
	}

	if t.feedbackCount < 0 {
		return fmt.Errorf("feedback count cannot be negative")
	}

	if t.feedbackCount == 0 {
		return fmt.Errorf("topic analysis must have at least one feedback")
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
