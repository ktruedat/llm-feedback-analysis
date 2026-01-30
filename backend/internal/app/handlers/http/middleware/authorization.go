package middleware

import (
	"net/http"
	"strings"

	ce "github.com/ktruedat/llm-feedback-analysis/pkg/errors"
	"github.com/ktruedat/llm-feedback-analysis/pkg/http/responder"
	"github.com/ktruedat/llm-feedback-analysis/pkg/tracelog"
)

// RequireRole creates a middleware that ensures the authenticated user has at least one of the required roles.
// This follows RBAC (Role-Based Access Control) pattern.
//
// RBAC Standard Approach:
// - 401 Unauthorized: User is not authenticated (handled by JWTMiddleware)
// - 403 Forbidden: User is authenticated but lacks required permissions/roles
// - Roles are stored in JWT claims as an array: ["admin", "user"]
//
// Usage:
//
//	router.Delete("/{id}", RequireRole("admin", logger, responder)(handler))
func RequireRole(
	requiredRole string,
	logger tracelog.TraceLogger,
	responder responder.RestResponder,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				logger := logger.WithSpan(ctx)

				// Get user claims from context (set by JWTMiddleware)
				claims := GetUserClaims(r)
				if claims == nil {
					// This should not happen if JWTMiddleware is applied before RequireRole
					// But we handle it gracefully
					logger.Warning("user claims not found in context - JWT middleware may be missing")
					responder.RespondContent(w, ce.ErrUnauthorized("authentication required"))
					return
				}

				// Check if user has the required role
				hasRole := false
				if claims.Roles != nil {
					for _, role := range claims.Roles {
						// Case-insensitive comparison (industry standard)
						if strings.EqualFold(strings.TrimSpace(role), strings.TrimSpace(requiredRole)) {
							hasRole = true
							break
						}
					}
				}

				if !hasRole {
					userID := claims.UserID
					logger.Warning(
						"access denied - insufficient permissions",
						"required_role", requiredRole,
						"user_roles", claims.Roles,
						"user_id", userID,
					)
					responder.RespondContent(
						w, ce.ErrForbidden(
							"insufficient permissions - admin role required",
						),
					)
					return
				}

				// User has required role, proceed to handler
				next.ServeHTTP(w, r)
			},
		)
	}
}

// RequireAnyRole creates a middleware that ensures the authenticated user has at least one of the provided roles.
// This is useful when multiple roles can access a resource.
//
// Usage:
//
//	router.Get("/reports", RequireAnyRole([]string{"admin", "manager"}, logger, responder)(handler))
func RequireAnyRole(
	requiredRoles []string,
	logger tracelog.TraceLogger,
	responder responder.RestResponder,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				logger := logger.WithSpan(ctx)

				claims := GetUserClaims(r)
				if claims == nil {
					logger.Warning("user claims not found in context - JWT middleware may be missing")
					responder.RespondContent(w, ce.ErrUnauthorized("authentication required"))
					return
				}

				// Check if user has at least one of the required roles
				hasRequiredRole := false
				if claims.Roles != nil {
					for _, userRole := range claims.Roles {
						for _, requiredRole := range requiredRoles {
							if strings.EqualFold(strings.TrimSpace(userRole), strings.TrimSpace(requiredRole)) {
								hasRequiredRole = true
								break
							}
						}
						if hasRequiredRole {
							break
						}
					}
				}

				if !hasRequiredRole {
					userID := claims.UserID
					if userID == "" {
						userID = claims.Subject // Fallback to Subject if UserID not set
					}
					logger.Warning(
						"access denied - insufficient permissions",
						"required_roles", requiredRoles,
						"user_roles", claims.Roles,
						"user_id", userID,
					)
					responder.RespondContent(
						w, ce.ErrForbidden(
							"insufficient permissions - one of the following roles required: "+strings.Join(
								requiredRoles,
								", ",
							),
						),
					)
					return
				}

				next.ServeHTTP(w, r)
			},
		)
	}
}
