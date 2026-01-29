package feedback

import (
	"context"

	"github.com/google/uuid"
	"github.com/ktruedat/llm-feedback-analysis/internal/domain/feedback"
	"github.com/ktruedat/llm-feedback-analysis/pkg/trace"
)

func (s *svc) GetFeedbackByID(ctx context.Context, feedbackID uuid.UUID) (*feedback.Feedback, error) {
	logger := s.logger.WithSpan(ctx)
	logger.SetSpanAttributes(ctx, trace.Attribute{Key: "feedback_id", Value: feedbackID.String()})
	logger.Info("getting feedback by id", "feedback_id", feedbackID.String())

	fb, err := s.feedRepo.Get(ctx, feedbackID)
	if err != nil {
		logger.RecordSpanError(ctx, err)
		return nil, s.errChecker.Check(err)
	}

	return fb, nil
}
