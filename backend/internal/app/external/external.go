package external

import (
	"context"

	"github.com/google/uuid"
	"github.com/ktruedat/llm-feedback-analysis/internal/domain/analysis"
	"github.com/ktruedat/llm-feedback-analysis/internal/domain/feedback"
)

// LLMClient defines the interface for LLM operations.
type LLMClient interface {
	// AnalyzeFeedbacks performs LLM analysis on the given feedbacks.
	// Returns the analysis result with summary, sentiment, insights, etc.
	AnalyzeFeedbacks(
		ctx context.Context,
		feedbacks []*feedback.Feedback,
		previousAnalysis *analysis.Analysis,
	) (*AnalysisResult, error)
}

// AnalysisResult contains the result of an LLM analysis.
type AnalysisResult struct {
	OverallSummary string
	Sentiment      analysis.Sentiment
	KeyInsights    []string
	TokensUsed     int
	Topics         []Topic
}

// Topic represents a topic identified by the LLM.
type Topic struct {
	Name        string
	Description string
	FeedbackIDs []uuid.UUID
	Sentiment   analysis.Sentiment
}
