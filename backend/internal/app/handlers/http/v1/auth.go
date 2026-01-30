package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/transport/requests"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/transport/responses"
	ce "github.com/ktruedat/llm-feedback-analysis/pkg/errors"
	"github.com/ktruedat/llm-feedback-analysis/pkg/http/responder"
	"github.com/ktruedat/llm-feedback-analysis/pkg/trace"
)

func (h *Handlers) registerAuthRoutes(router chi.Router) {
	router.Route(
		"/auth", func(r chi.Router) {
			r.Post("/register", trace.InstrumentHandlerFunc(h.RegisterUser, "POST /auth/register", h))
			r.Post("/login", trace.InstrumentHandlerFunc(h.LoginUser, "POST /auth/login", h))
		},
	)
}

// RegisterUser handles user registration
//
//	@Summary		Register a new user
//	@Description	Create a new user account with email and password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		requests.RegisterUserRequest	true	"User registration request"
//	@Success		201		{object}	responses.RegisterUserResponse	"User registered successfully"
//	@Failure		400		{object}	map[string]interface{}			"Bad request - invalid request body or validation error"
//	@Failure		409		{object}	map[string]interface{}			"Conflict - email already exists"
//	@Failure		500		{object}	map[string]interface{}			"Internal server error"
//	@Router			/auth/register [post]
func (h *Handlers) RegisterUser(resp http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.WithSpan(ctx)

	var req requests.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.RecordSpanError(ctx, err)
		h.responder.RespondContent(resp, ce.ErrBadRequest("invalid request body", ce.WithCauseError(err)))
		return
	}

	logger.Info("registering user", "email", req.Email)
	user, err := h.userService.RegisterUser(ctx, &req)
	if err != nil {
		logger.RecordSpanError(ctx, err)
		logger.Error("error registering user", err, "email", req.Email)
		h.handleSvcError(resp, err)
		return
	}

	response := responses.RegisterUserResponseFromDomain(user)
	h.responder.RespondContent(resp, responder.NewGenericResponse(http.StatusCreated, response))
}

// LoginUser handles user authentication
//
//	@Summary		Login user
//	@Description	Authenticate user with email and password, returns JWT token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		requests.LoginUserRequest	true	"User login request"
//	@Success		200		{object}	responses.LoginUserResponse	"User authenticated successfully"
//	@Failure		400		{object}	map[string]interface{}		"Bad request - invalid request body"
//	@Failure		401		{object}	map[string]interface{}		"Unauthorized - invalid credentials"
//	@Failure		500		{object}	map[string]interface{}		"Internal server error"
//	@Router			/auth/login [post]
func (h *Handlers) LoginUser(resp http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.WithSpan(ctx)

	var req requests.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.RecordSpanError(ctx, err)
		h.responder.RespondContent(resp, ce.ErrBadRequest("invalid request body", ce.WithCauseError(err)))
		return
	}

	logger.Info("authenticating user", "email", req.Email)
	token, authenticatedUser, err := h.userService.AuthenticateUser(ctx, &req)
	if err != nil {
		logger.RecordSpanError(ctx, err)
		logger.Error("error authenticating user", err, "email", req.Email)
		h.handleSvcError(resp, err)
		return
	}

	// Calculate expiration in seconds from config
	expirationHours := h.jwtCfg.ExpirationHours
	if expirationHours <= 0 {
		expirationHours = 24 // Default to 24 hours
	}
	expiresInSeconds := expirationHours * 3600

	response := responses.LoginUserResponse{
		Token:     token,
		ExpiresIn: expiresInSeconds,
		User:      responses.UserInfoFromDomain(authenticatedUser),
	}

	h.responder.RespondContent(resp, responder.NewGenericResponse(http.StatusOK, response))
}
