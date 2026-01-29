//nolint:lll // cannot split tags
package requests

// CreateFeedbackRequest represents the request payload for creating a feedback
//
//	@Description	Request payload for creating a new feedback submission.
type CreateFeedbackRequest struct {
	Rating  int    `json:"rating" example:"5" binding:"required,min=1,max=5"` // Rating value from 1 to 5 (required)
	Comment string `json:"comment" example:"Really nice!"`                    // Feedback comment text, 1-1000 characters (required)
}
