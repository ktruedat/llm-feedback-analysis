package middleware

import (
	"context"
	"net/http"

	"github.com/ktruedat/llm-feedback-analysis/internal/app/config"
	appjwt "github.com/ktruedat/llm-feedback-analysis/internal/app/infrastructure/jwt"
	ce "github.com/ktruedat/llm-feedback-analysis/pkg/errors"
	"github.com/ktruedat/llm-feedback-analysis/pkg/http/responder"
	"github.com/ktruedat/llm-feedback-analysis/pkg/tracelog"
)

// ContextKey is a type for context keys to avoid collisions.
type ContextKey string

const (
	// UserClaimsContextKey is the key used to store JWT claims in the request context.
	UserClaimsContextKey ContextKey = "user_claims"
)

// JWTMiddleware creates a middleware that validates JWT bearer tokens
// according to RFC 6750 (OAuth 2.0 Bearer Token Usage)
func JWTMiddleware(
	cfg *config.JWT,
	logger tracelog.TraceLogger,
	responder responder.RestResponder,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				logger := logger.WithSpan(ctx)

				// Extract bearer token from Authorization header
				authHeader := r.Header.Get("Authorization")
				tokenString, err := appjwt.ExtractBearerToken(authHeader)
				if err != nil {
					logger.Warning("failed to extract bearer token", "error", err)
					responder.RespondContent(w, ce.ErrUnauthorized(err.Error()))
					return
				}

				// Parse and validate the token
				claims, err := appjwt.ParseToken(tokenString, cfg)
				if err != nil {
					logger.Warning("jwt validation failed", "error", err)
					responder.RespondContent(w, ce.ErrUnauthorized("invalid or expired token", ce.WithCauseError(err)))
					return
				}

				// Store claims in context for handlers to access
				ctx = context.WithValue(ctx, UserClaimsContextKey, claims)
				next.ServeHTTP(w, r.WithContext(ctx))
			},
		)
	}
}

// GetUserClaims extracts user claims from the request context
// Returns nil if no claims are found in the context
func GetUserClaims(r *http.Request) *appjwt.Claims {
	claims, ok := r.Context().Value(UserClaimsContextKey).(*appjwt.Claims)
	if !ok {
		return nil
	}
	return claims
}
