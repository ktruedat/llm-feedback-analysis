package feedback

import (
	"time"

	"github.com/google/uuid"
	"github.com/ktruedat/llm-feedback-analysis/pkg/optional"
)

// ID returns the feedback ID.
func (f *Feedback) ID() uuid.UUID {
	return f.id
}

// UserID returns the user ID who submitted the feedback.
func (f *Feedback) UserID() uuid.UUID {
	return f.userID
}

// Rating returns the feedback rating.
func (f *Feedback) Rating() Rating {
	return f.rating
}

// Comment returns the feedback comment.
func (f *Feedback) Comment() Comment {
	return f.comment
}

// CreatedAt returns the creation timestamp.
func (f *Feedback) CreatedAt() time.Time {
	return f.createdAt
}

// UpdatedAt returns the last update timestamp.
func (f *Feedback) UpdatedAt() time.Time {
	return f.updatedAt
}

// DeletedAt returns the deletion timestamp if deleted.
func (f *Feedback) DeletedAt() optional.Optional[time.Time] {
	return f.deletedAt
}
