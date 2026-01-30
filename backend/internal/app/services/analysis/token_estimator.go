package analysis

import (
	"github.com/ktruedat/llm-feedback-analysis/internal/domain/analysis"
	"github.com/ktruedat/llm-feedback-analysis/internal/domain/feedback"
)

// estimateTokens estimates the number of tokens for a given text.
// Uses a simple approximation: ~4 characters per token for English text.
// This is a conservative estimate that works well for most cases.
func estimateTokens(text string) int {
	// Rough approximation: 4 characters per token
	// Add some overhead for JSON structure, whitespace, etc.
	return len(text)/4 + 10
}

// estimateFeedbackTokens estimates tokens for a single feedback.
func estimateFeedbackTokens(fb *feedback.Feedback) int {
	// Estimate tokens for feedback JSON representation
	// Format: {"id": "...", "rating": 5, "comment": "..."}
	idTokens := estimateTokens(fb.ID().String())
	commentTokens := estimateTokens(fb.Comment().Value())
	// Rating is just a number, minimal tokens
	ratingTokens := 1
	// JSON structure overhead
	structureTokens := 20

	return idTokens + commentTokens + ratingTokens + structureTokens
}

// estimateSystemPromptTokens estimates tokens for the system prompt.
func estimateSystemPromptTokens() int {
	// System prompt is fixed, estimate once
	systemPrompt := `You are an expert feedback analyst. Your task is to analyze customer feedback and provide:
1. An overall summary of all feedback
2. The overall sentiment (positive, mixed, or negative)
3. Key insights as bullet points
4. Topics/themes identified in the feedback, with each topic including:
   - A descriptive name
   - A detailed description
   - The feedback IDs that belong to this topic (a feedback can belong to multiple topics)
   - The sentiment for this specific topic

When analyzing topics:
- Group similar feedback together
- A single feedback can belong to multiple topics if it addresses multiple themes
- Provide clear, actionable insights
- Be specific about which feedback IDs map to which topics`
	return estimateTokens(systemPrompt)
}

// estimatePreviousAnalysisTokens estimates tokens for previous analysis context.
func estimatePreviousAnalysisTokens(prevAnalysis *analysis.Analysis) int {
	if prevAnalysis == nil {
		return 0
	}

	// Estimate tokens for previous analysis summary
	summaryTokens := estimateTokens(prevAnalysis.OverallSummary())
	insightsTokens := 0
	for _, insight := range prevAnalysis.KeyInsights() {
		insightsTokens += estimateTokens(insight)
	}
	// JSON structure overhead
	structureTokens := 50

	return summaryTokens + insightsTokens + structureTokens
}

// estimateResponseTokens estimates tokens for the expected response.
// This is a conservative estimate for the structured JSON response.
func estimateResponseTokens(feedbackCount int) int {
	// Estimate based on feedback count
	// Each feedback might generate some content in topics
	// Conservative estimate: ~100 tokens per feedback for response
	baseTokens := 200        // Base response structure
	perFeedbackTokens := 100 // Estimated tokens per feedback in response

	return baseTokens + (feedbackCount * perFeedbackTokens)
}

// estimateTotalTokens estimates total tokens for an analysis request.
func estimateTotalTokens(feedbacks []*feedback.Feedback, previousAnalysis *analysis.Analysis) int {
	systemPromptTokens := estimateSystemPromptTokens()
	previousAnalysisTokens := estimatePreviousAnalysisTokens(previousAnalysis)

	feedbackTokens := 0
	for _, fb := range feedbacks {
		feedbackTokens += estimateFeedbackTokens(fb)
	}

	// User payload JSON structure overhead
	userPayloadOverhead := 50

	// Estimate response tokens
	responseTokens := estimateResponseTokens(len(feedbacks))

	// Total: system prompt + user payload (previous analysis + feedbacks + overhead) + response
	return systemPromptTokens + previousAnalysisTokens + feedbackTokens + userPayloadOverhead + responseTokens
}

// selectFeedbacksForAnalysis selects feedbacks that fit within token and count limits.
// Returns the selected feedbacks and the remaining feedbacks that should stay in the queue.
func (a *analyzer) selectFeedbacksForAnalysis(
	pendingFeedbacks []*feedback.Feedback,
	previousAnalysis *analysis.Analysis,
) (selected []*feedback.Feedback, remaining []*feedback.Feedback) {
	if len(pendingFeedbacks) == 0 {
		return nil, nil
	}

	selected = make([]*feedback.Feedback, 0)
	maxTokens := a.cfg.MaxTokensPerRequest
	maxFeedbacks := a.cfg.MaxFeedbacksInContext

	// First, apply max feedbacks limit (secondary constraint)
	candidates := pendingFeedbacks
	if len(candidates) > maxFeedbacks {
		candidates = candidates[:maxFeedbacks]
	}

	// Then, apply token limit (primary constraint)
	currentTokens := estimateSystemPromptTokens() + estimatePreviousAnalysisTokens(previousAnalysis)
	userPayloadOverhead := 50
	responseTokensEstimate := 200 // Base response tokens

	for i, fb := range candidates {
		feedbackTokens := estimateFeedbackTokens(fb)
		// Estimate response tokens for this feedback
		estimatedResponseTokens := 100 // Per feedback in response

		// Check if adding this feedback would exceed token limit
		estimatedTotalTokens := currentTokens + feedbackTokens + userPayloadOverhead + responseTokensEstimate + (len(selected)+1)*estimatedResponseTokens

		if estimatedTotalTokens > maxTokens {
			// This feedback would exceed token limit, stop here
			remaining = candidates[i:]
			break
		}

		// Add this feedback
		selected = append(selected, fb)
		currentTokens += feedbackTokens
	}

	// If we didn't use all candidates, add the rest to remaining
	if len(selected) < len(candidates) {
		remaining = append(remaining, candidates[len(selected):]...)
	}

	// Add any feedbacks beyond maxFeedbacks to remaining
	if len(pendingFeedbacks) > maxFeedbacks {
		remaining = append(remaining, pendingFeedbacks[maxFeedbacks:]...)
	}

	return selected, remaining
}
