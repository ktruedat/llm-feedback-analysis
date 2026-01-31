package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/transport/responses"
	ce "github.com/ktruedat/llm-feedback-analysis/pkg/errors"
	"github.com/ktruedat/llm-feedback-analysis/pkg/http/responder"
	"github.com/ktruedat/llm-feedback-analysis/pkg/trace"
)

func (h *Handlers) registerAnalysisRoutes(router chi.Router) {
	router.Route(
		"/analyses", func(r chi.Router) {
			r.Get("/latest", trace.InstrumentHandlerFunc(h.GetLatestAnalysis, "GET /analyses/latest", h))
			r.Get("/", trace.InstrumentHandlerFunc(h.ListAnalyses, "GET /analyses", h))
			r.Get("/{id}", trace.InstrumentHandlerFunc(h.GetAnalysisByID, "GET /analyses/{id}", h))
		},
	)
}

// GetLatestAnalysis retrieves the latest completed analysis
//
//	@Summary		Get latest analysis
//	@Description	Retrieve the most recent completed analysis for the dashboard
//	@Tags			analyses
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	responses.AnalysisResponse	"Latest analysis retrieved successfully"
//	@Success		204	{object}	nil							"No analysis found"
//	@Failure		401	{object}	map[string]interface{}		"Unauthorized - invalid or missing JWT token"
//	@Failure		500	{object}	map[string]interface{}		"Internal server error"
//	@Router			/analyses/latest [get]
func (h *Handlers) GetLatestAnalysis(resp http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.WithSpan(ctx)

	logger.Info("getting latest analysis")
	analysis, err := h.feedbackSummaryService.GetLatestAnalysis(ctx)
	if err != nil {
		logger.RecordSpanError(ctx, err)
		logger.Error("error getting latest analysis", err)
		h.handleSvcError(resp, err)
		return
	}

	if analysis == nil {
		h.responder.RespondContent(resp, responder.NewGenericResponse(http.StatusNoContent, nil))
		return
	}

	response := responses.AnalysisResponseFromDomain(analysis)
	h.responder.RespondContent(resp, responder.NewGenericResponse(http.StatusOK, response))
}

// ListAnalyses retrieves all analyses ordered by creation date (newest first)
//
//	@Summary		List all analyses
//	@Description	Retrieve all analyses ordered by creation date (newest first) for the history page
//	@Tags			analyses
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	responses.AnalysisListResponse	"Analyses retrieved successfully"
//	@Failure		401	{object}	map[string]interface{}			"Unauthorized - invalid or missing JWT token"
//	@Failure		500	{object}	map[string]interface{}			"Internal server error"
//	@Router			/analyses [get]
func (h *Handlers) ListAnalyses(resp http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.WithSpan(ctx)

	logger.Info("listing all analyses")
	analyses, err := h.feedbackSummaryService.GetAllAnalyses(ctx)
	if err != nil {
		logger.RecordSpanError(ctx, err)
		logger.Error("error listing analyses", err)
		h.handleSvcError(resp, err)
		return
	}

	// Convert to response format
	analysisResponses := make([]responses.AnalysisResponse, len(analyses))
	for i, a := range analyses {
		analysisResponses[i] = *responses.AnalysisResponseFromDomain(a)
	}

	response := responses.AnalysisListResponse{
		Analyses: analysisResponses,
		Total:    len(analysisResponses),
	}

	h.responder.RespondContent(resp, responder.NewGenericResponse(http.StatusOK, response))
}

// GetAnalysisByID retrieves an analysis by ID with its topics and analyzed feedbacks
//
//	@Summary		Get analysis by ID
//	@Description	Retrieve a specific analysis with its topics and analyzed feedbacks with their associated topics
//	@Tags			analyses
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Analysis ID"	example(550e8400-e29b-41d4-a716-446655440000)
//	@Success		200	{object}	responses.AnalysisDetailResponse	"Analysis retrieved successfully"
//	@Failure		400	{object}	map[string]interface{}			"Bad request - invalid analysis ID format"
//	@Failure		401	{object}	map[string]interface{}			"Unauthorized - invalid or missing JWT token"
//	@Failure		404	{object}	map[string]interface{}			"Analysis not found"
//	@Failure		500	{object}	map[string]interface{}			"Internal server error"
//	@Router			/analyses/{id} [get]
func (h *Handlers) GetAnalysisByID(resp http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.WithSpan(ctx)

	analysisIDStr := chi.URLParam(r, "id")
	analysisID, err := uuid.Parse(analysisIDStr)
	if err != nil {
		h.responder.RespondContent(resp, ce.ErrBadRequest("invalid analysis ID format"))
		return
	}

	logger.Info("getting analysis by ID", "analysis_id", analysisID)
	analysisEntity, topics, feedbackTopics, err := h.feedbackSummaryService.GetAnalysisByID(ctx, analysisID)
	if err != nil {
		logger.RecordSpanError(ctx, err)
		logger.Error("error getting analysis", err, "analysis_id", analysisID)
		h.handleSvcError(resp, err)
		return
	}

	// Convert topics to response format
	topicResponses := make([]responses.TopicAnalysisResponse, len(topics))
	for i, topic := range topics {
		topicResponses[i] = *responses.TopicAnalysisResponseFromDomain(topic)
	}

	// Get all feedbacks and build response with topics
	feedbackResponses := make([]responses.FeedbackWithTopicsResponse, 0)
	for feedbackID, associatedTopics := range feedbackTopics {
		feedback, err := h.feedbackService.GetFeedbackByID(ctx, feedbackID)
		if err != nil {
			logger.Warning("error getting feedback", "feedback_id", feedbackID, "error", err.Error())
			continue
		}

		// Extract topic enum values
		topicEnums := make([]string, len(associatedTopics))
		for i, topic := range associatedTopics {
			topicEnums[i] = string(topic.Topic())
		}

		feedbackResponses = append(
			feedbackResponses, responses.FeedbackWithTopicsResponse{
				ID:        feedback.ID().String(),
				Rating:    feedback.Rating().Value(),
				Comment:   feedback.Comment().Value(),
				CreatedAt: feedback.CreatedAt(),
				Topics:    topicEnums,
			},
		)
	}

	response := responses.AnalysisDetailResponse{
		Analysis:  responses.AnalysisResponseFromDomain(analysisEntity),
		Topics:    topicResponses,
		Feedbacks: feedbackResponses,
	}

	h.responder.RespondContent(resp, responder.NewGenericResponse(http.StatusOK, response))
}
