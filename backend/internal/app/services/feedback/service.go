package feedback

import (
	"github.com/ktruedat/llm-feedback-analysis/internal/app/config"
	apprepo "github.com/ktruedat/llm-feedback-analysis/internal/app/repository"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/services"
	"github.com/ktruedat/llm-feedback-analysis/pkg/errors"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository"
	"github.com/ktruedat/llm-feedback-analysis/pkg/tracelog"
)

type svc struct {
	logger        tracelog.TraceLogger
	paginationCfg *config.Pagination
	errChecker    errors.ErrorChecker
	feedRepo      apprepo.FeedbackRepository
	transactor    repository.Transactor
	analyzer      services.AnalyzerService
}

func NewFeedbackService(
	traceLogger tracelog.TraceLogger,
	paginationCfg *config.Pagination,
	errChecker errors.ErrorChecker,
	feedRepo apprepo.FeedbackRepository,
	transactor repository.Transactor,
	analyzer services.AnalyzerService,
) services.FeedbackService {
	return &svc{
		logger:        traceLogger.NewGroup("feedback_service"),
		paginationCfg: paginationCfg,
		errChecker:    errChecker,
		feedRepo:      feedRepo,
		transactor:    transactor,
		analyzer:      analyzer,
	}
}
