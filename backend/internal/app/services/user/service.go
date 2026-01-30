package user

import (
	"github.com/ktruedat/llm-feedback-analysis/internal/app/config"
	apprepo "github.com/ktruedat/llm-feedback-analysis/internal/app/repository"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/services"
	"github.com/ktruedat/llm-feedback-analysis/pkg/errors"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository"
	"github.com/ktruedat/llm-feedback-analysis/pkg/tracelog"
)

type svc struct {
	logger     tracelog.TraceLogger
	errChecker errors.ErrorChecker
	userRepo   apprepo.UserRepository
	jwtCfg     *config.JWT
	transactor repository.Transactor
}

func NewUserService(
	traceLogger tracelog.TraceLogger,
	errChecker errors.ErrorChecker,
	userRepo apprepo.UserRepository,
	jwtCfg *config.JWT,
	transactor repository.Transactor,
) services.UserService {
	return &svc{
		logger:     traceLogger.NewGroup("user_service"),
		errChecker: errChecker,
		userRepo:   userRepo,
		jwtCfg:     jwtCfg,
		transactor: transactor,
	}
}
