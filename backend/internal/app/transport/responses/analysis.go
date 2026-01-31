package responses

import (
	"time"

	"github.com/ktruedat/llm-feedback-analysis/internal/domain/analysis"
	"github.com/ktruedat/llm-feedback-analysis/pkg/optional"
)

// AnalysisResponse represents the response payload for an analysis
//
//	@Description	Response payload containing analysis details.
type AnalysisResponse struct {
	ID                 string                       `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	PreviousAnalysisID optional.Optional[string]    `json:"previous_analysis_id,omitempty" swaggertype:"primitive,string"`
	PeriodStart        time.Time                    `json:"period_start" example:"2024-01-01T00:00:00Z"`
	PeriodEnd          time.Time                    `json:"period_end" example:"2024-01-31T23:59:59Z"`
	FeedbackCount      int                          `json:"feedback_count" example:"100"`
	NewFeedbackCount   optional.Optional[int]       `json:"new_feedback_count,omitempty" swaggertype:"primitive,integer"`
	OverallSummary     string                       `json:"overall_summary"`
	Sentiment          string                       `json:"sentiment" example:"positive"`
	KeyInsights        []string                     `json:"key_insights"`
	Model              string                       `json:"model" example:"gpt-5-mini"`
	Tokens             int                          `json:"tokens" example:"5000"`
	AnalysisDurationMs int                          `json:"analysis_duration_ms" example:"5000"`
	Status             string                       `json:"status" example:"success"`
	FailureReason      optional.Optional[string]    `json:"failure_reason,omitempty" swaggertype:"primitive,string"`
	CreatedAt          time.Time                    `json:"created_at" example:"2024-01-01T00:00:00Z"`
	CompletedAt        optional.Optional[time.Time] `json:"completed_at,omitempty" swaggertype:"primitive,string"`
}

// AnalysisResponseFromDomain converts a domain Analysis entity to an AnalysisResponse.
func AnalysisResponseFromDomain(a *analysis.Analysis) *AnalysisResponse {
	resp := &AnalysisResponse{
		ID:                 a.ID().String(),
		PeriodStart:        a.PeriodStart(),
		PeriodEnd:          a.PeriodEnd(),
		FeedbackCount:      a.FeedbackCount(),
		OverallSummary:     a.OverallSummary(),
		Sentiment:          string(a.Sentiment()),
		KeyInsights:        a.KeyInsights(),
		Model:              a.Model(),
		Tokens:             a.Tokens(),
		AnalysisDurationMs: a.AnalysisDurationMs(),
		Status:             string(a.Status()),
		CreatedAt:          a.CreatedAt(),
	}

	if a.PreviousAnalysisID().IsSome() {
		resp.PreviousAnalysisID = optional.Some(a.PreviousAnalysisID().Unwrap().String())
	}

	if a.NewFeedbackCount().IsSome() {
		resp.NewFeedbackCount = optional.Some(a.NewFeedbackCount().Unwrap())
	}

	if a.FailureReason().IsSome() {
		resp.FailureReason = optional.Some(a.FailureReason().Unwrap())
	}

	if a.CompletedAt().IsSome() {
		resp.CompletedAt = optional.Some(a.CompletedAt().Unwrap())
	}

	return resp
}

// TopicAnalysisResponse represents the response payload for a topic analysis
//
//	@Description	Response payload containing topic analysis details.
type TopicAnalysisResponse struct {
	ID            string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Topic         string    `json:"topic" example:"product_functionality_features"`
	TopicName     string    `json:"topic_name" example:"Product Functionality & Features"`
	Summary       string    `json:"summary"`
	FeedbackCount int       `json:"feedback_count" example:"10"`
	Sentiment     string    `json:"sentiment" example:"positive"`
	CreatedAt     time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt     time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// TopicAnalysisResponseFromDomain converts a domain TopicAnalysis entity to a TopicAnalysisResponse.
func TopicAnalysisResponseFromDomain(ta *analysis.TopicAnalysis) *TopicAnalysisResponse {
	return &TopicAnalysisResponse{
		ID:            ta.ID().String(),
		Topic:         string(ta.Topic()),
		TopicName:     ta.TopicName(),
		Summary:       ta.Summary(),
		FeedbackCount: ta.FeedbackCount(),
		Sentiment:     string(ta.Sentiment()),
		CreatedAt:     ta.CreatedAt(),
		UpdatedAt:     ta.UpdatedAt(),
	}
}

// AnalysisDetailResponse represents the detailed response for an analysis with topics and feedbacks
//
//	@Description	Response payload containing detailed analysis with topics and feedback IDs.
type AnalysisDetailResponse struct {
	Analysis  *AnalysisResponse            `json:"analysis"`
	Topics    []TopicAnalysisResponse      `json:"topics"`
	Feedbacks []FeedbackWithTopicsResponse `json:"feedbacks"`
}

// FeedbackWithTopicsResponse represents a feedback with its associated topics
//
//	@Description	Response payload containing feedback details with associated topics.
type FeedbackWithTopicsResponse struct {
	ID        string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Rating    int       `json:"rating" example:"5"`
	Comment   string    `json:"comment" example:"Great product!"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	Topics    []string  `json:"topics"` // Topic enum values
}

// AnalysisListResponse represents a list of analysis responses
//
//	@Description	Response payload containing a list of analyses.
type AnalysisListResponse struct {
	Analyses []AnalysisResponse `json:"analyses"`
	Total    int                `json:"total" example:"10"`
}

// TopicStatsResponse represents statistics for a topic
//
//	@Description	Response payload containing topic statistics from the latest analysis.
type TopicStatsResponse struct {
	Topic         string  `json:"topic" example:"product_functionality_features"`
	TopicName     string  `json:"topic_name" example:"Product Functionality & Features"`
	FeedbackCount int     `json:"feedback_count" example:"10"`
	AverageRating float64 `json:"average_rating" example:"4.5"`
}

// TopicStatsListResponse represents a list of topic statistics
//
//	@Description	Response payload containing a list of topic statistics.
type TopicStatsListResponse struct {
	Topics []TopicStatsResponse `json:"topics"`
	Total  int                  `json:"total" example:"13"`
}

// TopicDetailsResponse represents detailed information about a topic with all associated feedbacks
//
//	@Description	Response payload containing detailed topic information with feedbacks.
type TopicDetailsResponse struct {
	Topic            string             `json:"topic" example:"product_functionality_features"`
	TopicName        string             `json:"topic_name" example:"Product Functionality & Features"`
	TopicDescription string             `json:"topic_description"`
	Summary          string             `json:"summary"`
	FeedbackCount    int                `json:"feedback_count" example:"10"`
	AverageRating    float64            `json:"average_rating" example:"4.5"`
	Sentiment        string             `json:"sentiment" example:"positive"`
	Feedbacks        []FeedbackResponse `json:"feedbacks"`
}
