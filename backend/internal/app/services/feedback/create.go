package feedback

import (
	"context"
	"fmt"

	apprepo "github.com/ktruedat/llm-feedback-analysis/internal/app/repository"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/transport/requests"
	"github.com/ktruedat/llm-feedback-analysis/internal/domain/feedback"
	"github.com/ktruedat/llm-feedback-analysis/pkg/errors"
	"github.com/ktruedat/llm-feedback-analysis/pkg/operations"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository"
	"github.com/ktruedat/llm-feedback-analysis/pkg/trace"
	"github.com/ktruedat/llm-feedback-analysis/pkg/tracelog"
)

func (s *svc) CreateFeedback(ctx context.Context, req *requests.CreateFeedbackRequest) (*feedback.Feedback, error) {
	logger := s.logger.WithSpan(ctx)
	ctx, spanLogger, span := logger.StartSpan(ctx, "feedback_service.create_feedback")
	defer span.End()

	span.SetAttributes(
		trace.Attribute{Key: "rating", Value: req.Rating},
		trace.Attribute{Key: "comment_length", Value: len(req.Comment)},
	)

	fb, err := s.createFeedback(ctx, req, spanLogger)
	if err != nil {
		span.SetStatus(trace.StatusError, err.Error())
		spanLogger.RecordSpanError(ctx, err)
		return nil, s.errChecker.Check(err)
	}

	span.SetStatus(trace.StatusOK, "Successfully created feedback")
	span.SetAttributes(trace.Attribute{Key: "feedback_id", Value: fb.ID().String()})
	return fb, nil
}

func (s *svc) createFeedback(
	ctx context.Context,
	req *requests.CreateFeedbackRequest,
	logger tracelog.TraceLogger,
) (*feedback.Feedback, error) {
	// Build domain value objects
	rating, err := feedback.NewRating(req.Rating)
	if err != nil {
		return nil, errors.ErrBadRequest("invalid rating", errors.WithCauseError(err))
	}

	comment, err := feedback.NewComment(req.Comment)
	if err != nil {
		return nil, errors.ErrBadRequest("invalid comment", errors.WithCauseError(err))
	}

	// Build domain entity
	builder := feedback.NewBuilder().
		WithRating(rating).
		WithComment(comment)

	fb, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build feedback: %w", err)
	}
	logger.Info("feedback built and validated", "feedback_id", fb.ID().String())

	// Create feedback in transaction
	if err := operations.RunGenericTransaction(
		ctx,
		s.transactor,
		s.createFeedbackRecord(fb, logger),
	); err != nil {
		logger.RecordSpanError(ctx, err)
		return nil, fmt.Errorf("failed to create feedback in transaction: %w", err)
	}

	logger.Info("feedback created successfully", "feedback_id", fb.ID().String())
	return fb, nil
}

func (s *svc) createFeedbackRecord(fb *feedback.Feedback, logger tracelog.TraceLogger) operations.TxExecFunc {
	return func(ctx context.Context, tx repository.Transaction) error {
		logger := logger.WithSpan(ctx)
		logger.Info("creating feedback record in database", "feedback_id", fb.ID().String())

		if err := s.feedRepo.Create(ctx, fb, repository.WithExecutor[apprepo.Options](tx)); err != nil {
			logger.RecordSpanError(ctx, err)
			return fmt.Errorf("failed to create feedback: %w", err)
		}

		logger.Info("feedback record created successfully")
		return nil
	}
}
