package feedback

import (
	"github.com/ktruedat/llm-feedback-analysis/internal/app/repository/postgres/feedback/sqlc"
	"github.com/ktruedat/llm-feedback-analysis/internal/domain/feedback"
)

// mapSQLCFeedbackToDomain maps a SQLC feedback model to a domain feedback entity.
func mapSQLCFeedbackToDomain(sqlcFeedback sqlc.Feedback) *feedback.Feedback {
	// Build rating value object
	rating, _ := feedback.NewRating(int(sqlcFeedback.Rating))

	// Build comment value object
	comment, _ := feedback.NewComment(sqlcFeedback.Comment)

	// Build domain entity using builder
	builder := feedback.NewBuilder().
		WithID(sqlcFeedback.ID).
		WithUserID(sqlcFeedback.UserID).
		WithRating(rating).
		WithComment(comment).
		WithCreatedAt(sqlcFeedback.CreatedAt).
		WithUpdatedAt(sqlcFeedback.UpdatedAt)

	// Handle deleted_at (nullable timestamp)
	if sqlcFeedback.DeletedAt != nil {
		builder.WithDeletedAt(*sqlcFeedback.DeletedAt)
	}

	return builder.BuildUnchecked()
}
