package analysis

import (
	"time"

	"github.com/google/uuid"
)

// ID returns the topic ID.
func (t *Topic) ID() uuid.UUID {
	return t.id
}

// AnalysisID returns the analysis ID this topic belongs to.
func (t *Topic) AnalysisID() uuid.UUID {
	return t.analysisID
}

// TopicName returns the topic name.
func (t *Topic) TopicName() string {
	return t.topicName
}

// Description returns the topic description.
func (t *Topic) Description() string {
	return t.description
}

// FeedbackCount returns the number of feedbacks belonging to this topic.
func (t *Topic) FeedbackCount() int {
	return t.feedbackCount
}

// Sentiment returns the sentiment for this topic.
func (t *Topic) Sentiment() Sentiment {
	return t.sentiment
}

// CreatedAt returns the creation timestamp.
func (t *Topic) CreatedAt() time.Time {
	return t.createdAt
}

// UpdatedAt returns the last update timestamp.
func (t *Topic) UpdatedAt() time.Time {
	return t.updatedAt
}
