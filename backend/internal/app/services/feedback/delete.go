package feedback

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	apprepo "github.com/ktruedat/llm-feedback-analysis/internal/app/repository"
	"github.com/ktruedat/llm-feedback-analysis/pkg/operations"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository"
	"github.com/ktruedat/llm-feedback-analysis/pkg/trace"
	"github.com/ktruedat/llm-feedback-analysis/pkg/tracelog"
)

func (s *svc) DeleteFeedback(ctx context.Context, feedbackID uuid.UUID) error {
	logger := s.logger.WithSpan(ctx)
	ctx, spanLogger, span := logger.StartSpan(ctx, "feedback_service.delete_feedback")
	defer span.End()

	span.SetAttributes(trace.Attribute{Key: "feedback_id", Value: feedbackID.String()})

	spanLogger.Info("deleting feedback", "feedback_id", feedbackID.String())

	err := s.deleteFeedback(ctx, feedbackID, spanLogger)
	if err != nil {
		span.SetStatus(trace.StatusError, err.Error())
		spanLogger.RecordSpanError(ctx, err)
		return s.errChecker.Check(err)
	}

	span.SetStatus(trace.StatusOK, "Successfully deleted feedback")
	return nil
}

func (s *svc) deleteFeedback(
	ctx context.Context,
	feedbackID uuid.UUID,
	logger tracelog.TraceLogger,
) error {
	// Check if feedback exists
	if _, err := s.feedRepo.Get(ctx, feedbackID); err != nil {
		return fmt.Errorf("feedback not found: %w", err)
	}

	if err := operations.RunGenericTransaction(
		ctx,
		s.transactor,
		s.deleteFeedbackRecord(feedbackID, logger),
	); err != nil {
		logger.RecordSpanError(ctx, err)
		return fmt.Errorf("failed to delete feedback in transaction: %w", err)
	}

	logger.Info("feedback deleted successfully", "feedback_id", feedbackID.String())
	return nil
}

func (s *svc) deleteFeedbackRecord(feedbackID uuid.UUID, logger tracelog.TraceLogger) operations.TxExecFunc {
	return func(ctx context.Context, tx repository.Transaction) error {
		logger := logger.WithSpan(ctx)
		logger.Info("deleting feedback record in database", "feedback_id", feedbackID.String())

		if err := s.feedRepo.Delete(ctx, feedbackID, repository.WithExecutor[apprepo.Options](tx)); err != nil {
			logger.RecordSpanError(ctx, err)
			return fmt.Errorf("failed to delete feedback: %w", err)
		}

		logger.Info("feedback record deleted successfully")
		return nil
	}
}
