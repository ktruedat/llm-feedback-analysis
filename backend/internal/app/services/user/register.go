package user

import (
	"context"
	"fmt"

	apprepo "github.com/ktruedat/llm-feedback-analysis/internal/app/repository"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/transport/requests"
	"github.com/ktruedat/llm-feedback-analysis/internal/domain/user"
	"github.com/ktruedat/llm-feedback-analysis/pkg/errors"
	"github.com/ktruedat/llm-feedback-analysis/pkg/operations"
	"github.com/ktruedat/llm-feedback-analysis/pkg/repository"
	"github.com/ktruedat/llm-feedback-analysis/pkg/trace"
	"github.com/ktruedat/llm-feedback-analysis/pkg/tracelog"
	"golang.org/x/crypto/bcrypt"
)

func (s *svc) RegisterUser(ctx context.Context, req *requests.RegisterUserRequest) (*user.User, error) {
	logger := s.logger.WithSpan(ctx)
	ctx, spanLogger, span := logger.StartSpan(ctx, "user_service.register_user")
	defer span.End()

	span.SetAttributes(
		trace.Attribute{Key: "email", Value: req.Email},
	)

	u, err := s.registerUser(ctx, req, spanLogger)
	if err != nil {
		span.SetStatus(trace.StatusError, err.Error())
		spanLogger.RecordSpanError(ctx, err)
		return nil, s.errChecker.Check(err)
	}

	span.SetStatus(trace.StatusOK, "Successfully registered user")
	span.SetAttributes(trace.Attribute{Key: "user_id", Value: u.ID().String()})
	return u, nil
}

func (s *svc) registerUser(
	ctx context.Context,
	req *requests.RegisterUserRequest,
	logger tracelog.TraceLogger,
) (*user.User, error) {
	email, err := user.NewEmail(req.Email)
	if err != nil {
		return nil, errors.ErrBadRequest("invalid email format", errors.WithCauseError(err))
	}

	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, email.Value())
	if err == nil && existingUser != nil {
		if existingUser.IsDeleted() {
			return nil, errors.ErrBadRequest("account with this email was previously deleted")
		}
		return nil, &errors.GenericError{
			Code:       errors.NewDomainErrorCode("email_already_exists", errors.CategoryConflict),
			Message:    "Email already exists",
			UserFacing: true,
		}
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("failed to hash password", err)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	passwordHashVO, err := user.NewPasswordHash(string(passwordHash))
	if err != nil {
		return nil, fmt.Errorf("failed to create password hash value object: %w", err)
	}

	builder := user.NewBuilder().
		WithEmail(email).
		WithPasswordHash(passwordHashVO).
		WithRole(user.RoleUser) // default role

	u, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build user: %w", err)
	}
	logger.Info("user built and validated", "user_id", u.ID().String(), "email", email.Value())

	if err := operations.RunGenericTransaction(
		ctx,
		s.transactor,
		s.createUserRecord(u, logger),
	); err != nil {
		logger.RecordSpanError(ctx, err)
		return nil, fmt.Errorf("failed to create user in transaction: %w", err)
	}

	logger.Info("user created successfully", "user_id", u.ID().String())
	return u, nil
}

func (s *svc) createUserRecord(u *user.User, logger tracelog.TraceLogger) operations.TxExecFunc {
	return func(ctx context.Context, tx repository.Transaction) error {
		logger := logger.WithSpan(ctx)
		logger.Info("creating user record in database", "user_id", u.ID().String())

		if err := s.userRepo.Create(ctx, u, repository.WithExecutor[apprepo.Options](tx)); err != nil {
			logger.RecordSpanError(ctx, err)
			return fmt.Errorf("failed to create user: %w", err)
		}

		logger.Info("user record created successfully")
		return nil
	}
}
