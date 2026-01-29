//nolint:lll // cannot split tags
package responses

import (
	"time"

	"github.com/ktruedat/llm-feedback-analysis/internal/domain/feedback"
	"github.com/ktruedat/llm-feedback-analysis/pkg/optional"
)

// FeedbackResponse represents the response payload for a feedback
//
//	@Description	Response payload containing feedback details.
type FeedbackResponse struct {
	ID        string                       `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`   // Feedback unique identifier
	Rating    int                          `json:"rating" example:"5"`                                  // Rating value from 1 to 5
	Comment   string                       `json:"comment" example:"Great service!"`                    // Feedback comment text
	CreatedAt time.Time                    `json:"created_at" example:"2024-01-01T00:00:00Z"`           // Creation timestamp
	UpdatedAt time.Time                    `json:"updated_at" example:"2024-01-01T00:00:00Z"`           // Last update timestamp
	DeletedAt optional.Optional[time.Time] `json:"deleted_at,omitempty" example:"2024-01-01T00:00:00Z"` // Deletion timestamp (if deleted)
}

// FeedbackResponseFromDomain converts a domain Feedback entity to a FeedbackResponse.
func FeedbackResponseFromDomain(fb *feedback.Feedback) *FeedbackResponse {
	resp := &FeedbackResponse{
		ID:        fb.ID().String(),
		Rating:    fb.Rating().Value(),
		Comment:   fb.Comment().Value(),
		CreatedAt: fb.CreatedAt(),
		UpdatedAt: fb.UpdatedAt(),
		DeletedAt: fb.DeletedAt(),
	}

	return resp
}

// FeedbackListResponse represents a list of feedback responses
//
//	@Description	Response payload containing a list of feedback entries.
type FeedbackListResponse struct {
	Feedbacks []FeedbackResponse `json:"feedbacks"`          // List of feedback entries
	Total     int                `json:"total" example:"10"` // Total number of feedback entries
}
