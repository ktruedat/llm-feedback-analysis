package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/transport/requests"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/transport/responses"
	ce "github.com/ktruedat/llm-feedback-analysis/pkg/errors"
	"github.com/ktruedat/llm-feedback-analysis/pkg/http/responder"
	"github.com/ktruedat/llm-feedback-analysis/pkg/trace"
)

func (h *Handlers) registerFeedbackRoutes(router chi.Router) {
	router.Route(
		"/feedbacks", func(r chi.Router) {
			r.Post("/", trace.InstrumentHandlerFunc(h.CreateFeedback, "POST /feedbacks", h))
			r.Get("/{id}", trace.InstrumentHandlerFunc(h.GetFeedbackByID, "GET /feedbacks/{id}", h))
			r.Get("/", trace.InstrumentHandlerFunc(h.ListFeedbacks, "GET /feedbacks", h))
			// TODO: only admin users should be able to delete feedbacks
			r.Delete("/{id}", trace.InstrumentHandlerFunc(h.DeleteFeedback, "DELETE /feedbacks/{id}", h))
		},
	)
}

// CreateFeedback creates a new feedback submission
//
//	@Summary		Create a new feedback
//	@Description	Create a new feedback submission with rating and comment
//	@Tags			feedbacks
//	@Accept			json
//	@Produce		json
//	@Param			request	body		requests.CreateFeedbackRequest	true	"Feedback creation request"
//	@Success		201		{object}	responses.FeedbackResponse		"Feedback created successfully"
//	@Failure		400		{object}	map[string]interface{}			"Bad request - invalid request body"
//	@Failure		500		{object}	map[string]interface{}			"Internal server error"
//	@Router			/feedbacks [post]
func (h *Handlers) CreateFeedback(resp http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.WithSpan(ctx)

	var req requests.CreateFeedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.RecordSpanError(ctx, err)
		h.responder.RespondContent(resp, ce.ErrBadRequest("invalid request body", ce.WithCauseError(err)))
		return
	}

	logger.Info("creating feedback", "rating", req.Rating)
	feedback, err := h.feedbackService.CreateFeedback(ctx, &req)
	if err != nil {
		logger.RecordSpanError(ctx, err)
		logger.Error("error creating feedback", err)
		h.handleSvcError(resp, err)
		return
	}

	response := responses.FeedbackResponseFromDomain(feedback)
	h.responder.RespondContent(resp, responder.NewGenericResponse(http.StatusCreated, response))
}

// GetFeedbackByID retrieves a feedback entry by its ID
//
//	@Summary		Get feedback by ID
//	@Description	Retrieve a specific feedback entry by its unique identifier
//	@Tags			feedbacks
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Feedback ID"
//	@Success		200	{object}	responses.FeedbackResponse	"Feedback retrieved successfully"
//	@Failure		400	{object}	map[string]interface{}		"Bad request - invalid feedback ID"
//	@Failure		404	{object}	map[string]interface{}		"Feedback not found"
//	@Failure		500	{object}	map[string]interface{}		"Internal server error"
//	@Router			/feedbacks/{id} [get]
func (h *Handlers) GetFeedbackByID(resp http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.WithSpan(ctx)

	feedbackIDStr := chi.URLParam(r, "id")
	feedbackID, err := uuid.Parse(feedbackIDStr)
	if err != nil {
		h.responder.RespondContent(resp, ce.ErrBadRequest("invalid feedback ID format"))
		return
	}

	logger.Info("getting feedback by ID", "feedback_id", feedbackID)
	feedback, err := h.feedbackService.GetFeedbackByID(ctx, feedbackID)
	if err != nil {
		logger.RecordSpanError(ctx, err)
		logger.Error("error getting feedback", err, "feedback_id", feedbackID)
		h.handleSvcError(resp, err)
		return
	}

	response := responses.FeedbackResponseFromDomain(feedback)
	h.responder.RespondContent(resp, responder.NewGenericResponse(http.StatusOK, response))
}

// ListFeedbacks retrieves a list of feedback entries
//
//	@Summary		List feedbacks
//	@Description	Retrieve a list of feedback entries with optional pagination
//	@Tags			feedbacks
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int		false	"Maximum number of feedbacks to return (default: 100)"
//	@Param			offset	query		int		false	"Number of feedbacks to skip (default: 0)"
//	@Success		200		{object}	responses.FeedbackListResponse	"Feedbacks retrieved successfully"
//	@Failure		400		{object}	map[string]interface{}			"Bad request - invalid query parameters"
//	@Failure		500		{object}	map[string]interface{}			"Internal server error"
//	@Router			/feedbacks [get]
func (h *Handlers) ListFeedbacks(resp http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.WithSpan(ctx)

	var limit, offset int
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := parseInt(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if parsedOffset, err := parseInt(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	logger.Info("listing feedbacks", "limit", limit, "offset", offset)
	feedbacks, err := h.feedbackService.ListFeedbacks(ctx, limit, offset)
	if err != nil {
		logger.RecordSpanError(ctx, err)
		logger.Error("error listing feedbacks", err)
		h.handleSvcError(resp, err)
		return
	}

	// Convert to response format
	feedbackResponses := make([]responses.FeedbackResponse, len(feedbacks))
	for i, fb := range feedbacks {
		feedbackResponses[i] = *responses.FeedbackResponseFromDomain(fb)
	}

	response := responses.FeedbackListResponse{
		Feedbacks: feedbackResponses,
		Total:     len(feedbackResponses),
	}

	h.responder.RespondContent(resp, responder.NewGenericResponse(http.StatusOK, response))
}

// DeleteFeedback performs a soft delete on a feedback entry
//
//	@Summary		Delete feedback
//	@Description	Soft delete a feedback entry by its unique identifier
//	@Tags			feedbacks
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Feedback ID"
//	@Success		204	{object}	nil		"Feedback deleted successfully"
//	@Failure		400	{object}	map[string]interface{}	"Bad request - invalid feedback ID"
//	@Failure		404	{object}	map[string]interface{}	"Feedback not found"
//	@Failure		500	{object}	map[string]interface{}	"Internal server error"
//	@Router			/feedbacks/{id} [delete]
func (h *Handlers) DeleteFeedback(resp http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.WithSpan(ctx)

	feedbackIDStr := chi.URLParam(r, "id")
	feedbackID, err := uuid.Parse(feedbackIDStr)
	if err != nil {
		h.responder.RespondContent(resp, ce.ErrBadRequest("invalid feedback ID format"))
		return
	}

	logger.Info("deleting feedback", "feedback_id", feedbackID)
	err = h.feedbackService.DeleteFeedback(ctx, feedbackID)
	if err != nil {
		logger.RecordSpanError(ctx, err)
		logger.Error("error deleting feedback", err, "feedback_id", feedbackID)
		h.handleSvcError(resp, err)
		return
	}

	h.responder.RespondContent(resp, responder.NewGenericResponse(http.StatusNoContent, nil))
}
