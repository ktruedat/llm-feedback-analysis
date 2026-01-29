package feedback

import (
	"context"
	"fmt"

	apprepo "github.com/ktruedat/llm-feedback-analysis/internal/app/repository"
	"github.com/ktruedat/llm-feedback-analysis/internal/domain/feedback"
	"github.com/ktruedat/llm-feedback-analysis/pkg/errors"
	"github.com/ktruedat/llm-feedback-analysis/pkg/trace"
	"github.com/ktruedat/llm-feedback-analysis/pkg/tracelog"
)

func (s *svc) ListFeedbacks(ctx context.Context, limit, offset int) ([]*feedback.Feedback, error) {
	logger := s.logger.WithSpan(ctx)
	ctx, spanLogger, span := logger.StartSpan(ctx, "feedback_service.list_feedbacks")
	defer span.End()

	span.SetAttributes(
		trace.Attribute{Key: "limit", Value: limit},
		trace.Attribute{Key: "offset", Value: offset},
	)
	spanLogger.Info("listing feedbacks", "limit", limit, "offset", offset)

	feedbacks, err := s.listFeedbacks(ctx, limit, offset, spanLogger)
	if err != nil {
		span.SetStatus(trace.StatusError, err.Error())
		spanLogger.RecordSpanError(ctx, err)
		return nil, s.errChecker.Check(err)
	}

	span.SetStatus(trace.StatusOK, "Successfully listed feedbacks")
	span.SetAttributes(trace.Attribute{Key: "count", Value: len(feedbacks)})
	return feedbacks, nil
}

func (s *svc) listFeedbacks(
	ctx context.Context,
	limit, offset int,
	logger tracelog.TraceLogger,
) ([]*feedback.Feedback, error) {
	if limit <= 0 {
		limit = s.paginationCfg.Limit
	}
	if limit > 1000 {
		return nil, errors.ErrBadRequest("limit cannot exceed 1000")
	}
	if offset < 0 {
		offset = s.paginationCfg.Offset
	}

	feedbacks, err := s.feedRepo.List(
		ctx,
		apprepo.WithOptions(
			&apprepo.Options{
				Limit:  limit,
				Offset: offset,
			},
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list feedbacks: %w", err)
	}

	logger.Info("feedbacks listed successfully", "count", len(feedbacks))
	return feedbacks, nil
}
