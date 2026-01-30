package user

import (
	"context"
	"fmt"

	appjwt "github.com/ktruedat/llm-feedback-analysis/internal/app/infrastructure/jwt"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/transport/requests"
	"github.com/ktruedat/llm-feedback-analysis/internal/domain/user"
	"github.com/ktruedat/llm-feedback-analysis/pkg/errors"
	"github.com/ktruedat/llm-feedback-analysis/pkg/trace"
	"github.com/ktruedat/llm-feedback-analysis/pkg/tracelog"
	"golang.org/x/crypto/bcrypt"
)

func (s *svc) AuthenticateUser(ctx context.Context, req *requests.LoginUserRequest) (string, error) {
	logger := s.logger.WithSpan(ctx)
	ctx, spanLogger, span := logger.StartSpan(ctx, "user_service.authenticate_user")
	defer span.End()

	span.SetAttributes(
		trace.Attribute{Key: "email", Value: req.Email},
	)

	token, err := s.authenticateUser(ctx, req, spanLogger)
	if err != nil {
		span.SetStatus(trace.StatusError, err.Error())
		spanLogger.RecordSpanError(ctx, err)
		return "", s.errChecker.Check(err)
	}

	span.SetStatus(trace.StatusOK, "Successfully authenticated user")
	return token, nil
}

func (s *svc) authenticateUser(
	ctx context.Context,
	req *requests.LoginUserRequest,
	logger tracelog.TraceLogger,
) (string, error) {
	email, err := user.NewEmail(req.Email)
	if err != nil {
		return "", errors.ErrBadRequest("invalid email format", errors.WithCauseError(err))
	}

	u, err := s.userRepo.GetByEmail(ctx, email.Value())
	if err != nil {
		logger.Warning("user not found", "email", email.Value())
		return "", errors.ErrUnauthorized("invalid email or password")
	}

	if u.IsDeleted() {
		return "", errors.ErrUnauthorized("account has been deleted")
	}

	if !u.IsActive() {
		return "", errors.ErrUnauthorized("account is not active")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash().Value()), []byte(req.Password)); err != nil {
		logger.Warning("invalid password", "email", email.Value())
		return "", errors.ErrUnauthorized("invalid email or password")
	}
	logger.Info("password verified successfully", "user_id", u.ID().String())

	// Convert domain roles to []string
	roles := u.Roles()
	roleStrings := make([]string, len(roles))
	for i, role := range roles {
		roleStrings[i] = role.String()
	}

	// Create claims with user ID and email
	claims := appjwt.NewClaims(u.ID(), u.Email().Value(), roleStrings, s.jwtCfg)

	// Generate JWT token
	token, err := appjwt.GenerateToken(claims, s.jwtCfg)
	if err != nil {
		logger.Error("failed to generate JWT token", err)
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	logger.Info("JWT token generated successfully", "user_id", u.ID().String())
	return token, nil
}
