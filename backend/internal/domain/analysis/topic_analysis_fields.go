package analysis

import (
	"time"

	"github.com/google/uuid"
)

// ID returns the topic analysis ID.
func (t *TopicAnalysis) ID() uuid.UUID {
	return t.id
}

// AnalysisID returns the analysis ID this topic analysis belongs to.
func (t *TopicAnalysis) AnalysisID() uuid.UUID {
	return t.analysisID
}

// Topic returns the topic enum.
func (t *TopicAnalysis) Topic() Topic {
	return t.topic
}

// TopicName returns the display name of the topic.
func (t *TopicAnalysis) TopicName() string {
	return t.topic.DisplayName()
}

// Summary returns the topic analysis summary.
func (t *TopicAnalysis) Summary() string {
	return t.summary
}

// FeedbackCount returns the number of feedbacks belonging to this topic analysis.
func (t *TopicAnalysis) FeedbackCount() int {
	return t.feedbackCount
}

// Sentiment returns the sentiment for this topic analysis.
func (t *TopicAnalysis) Sentiment() Sentiment {
	return t.sentiment
}

// CreatedAt returns the creation timestamp.
func (t *TopicAnalysis) CreatedAt() time.Time {
	return t.createdAt
}

// UpdatedAt returns the last update timestamp.
func (t *TopicAnalysis) UpdatedAt() time.Time {
	return t.updatedAt
}
