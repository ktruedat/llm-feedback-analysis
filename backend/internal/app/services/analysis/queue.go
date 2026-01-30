package analysis

import (
	"context"
	"time"

	"github.com/ktruedat/llm-feedback-analysis/internal/domain/feedback"
)

// addFeedbackToQueue adds a feedback to the pending queue.
func (a *analyzer) addFeedbackToQueue(fb *feedback.Feedback) {
	a.pendingMutex.Lock()
	defer a.pendingMutex.Unlock()

	a.pendingFeedbacks = append(a.pendingFeedbacks, fb)
	a.logger.Info(
		"feedback added to pending queue",
		"feedback_id",
		fb.ID().String(),
		"pending_count",
		len(a.pendingFeedbacks),
	)
}

// checkAndAnalyze checks if we should trigger an analysis based on configuration.
func (a *analyzer) checkAndAnalyze(ctx context.Context) {
	a.pendingMutex.Lock()
	pendingCount := len(a.pendingFeedbacks)
	pendingFeedbacks := make([]*feedback.Feedback, len(a.pendingFeedbacks))
	copy(pendingFeedbacks, a.pendingFeedbacks)
	a.pendingMutex.Unlock()

	if pendingCount < a.cfg.MinimumNewFeedbacksForAnalysis {
		return
	}

	// Check debounce if enabled
	if a.cfg.EnableDebounce {
		a.lastAnalysisMutex.Lock()
		timeSinceLastAnalysis := time.Since(a.lastAnalysisTime)
		a.lastAnalysisMutex.Unlock()

		if timeSinceLastAnalysis < time.Duration(a.cfg.DebounceMinutes)*time.Minute {
			return
		}
	}

	// Get previous analysis for token estimation
	previousAnalysis, err := a.analysisRepo.GetLatest(ctx)
	if err != nil {
		// No previous analysis, continue with nil
		previousAnalysis = nil
	}

	// Select feedbacks that fit within token and count limits
	selectedFeedbacks, remainingFeedbacks := a.selectFeedbacksForAnalysis(pendingFeedbacks, previousAnalysis)

	if len(selectedFeedbacks) == 0 {
		a.logger.Info("no feedbacks selected for analysis (token limit too restrictive)")
		return
	}

	// Update pending queue: remove selected feedbacks, keep remaining ones
	a.pendingMutex.Lock()
	// Remove selected feedbacks from the queue
	// We need to find and remove them
	selectedIDs := make(map[string]bool)
	for _, fb := range selectedFeedbacks {
		selectedIDs[fb.ID().String()] = true
	}

	newPending := make([]*feedback.Feedback, 0, len(a.pendingFeedbacks))
	for _, fb := range a.pendingFeedbacks {
		if !selectedIDs[fb.ID().String()] {
			newPending = append(newPending, fb)
		}
	}
	// Add remaining feedbacks back to queue (they weren't selected due to limits)
	newPending = append(newPending, remainingFeedbacks...)
	a.pendingFeedbacks = newPending
	a.pendingMutex.Unlock()

	if len(remainingFeedbacks) > 0 {
		a.logger.Info(
			"feedbacks returned to queue due to token/limit constraints",
			"selected_count", len(selectedFeedbacks),
			"remaining_count", len(remainingFeedbacks),
		)
	}

	// Trigger analysis with selected feedbacks only
	a.wg.Add(1)
	go a.performAnalysis(ctx, selectedFeedbacks)
}
