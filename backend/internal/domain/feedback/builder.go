package feedback

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ktruedat/llm-feedback-analysis/pkg/optional"
)

// Builder provides type-safe, fluent construction of Feedback entities.
// It accumulates validation errors and applies business rules at Build() time.
type Builder struct {
	entity           *Feedback
	validationErrors []error
}

// NewBuilder creates a builder for creating new feedback entities.
// Initialize with sensible defaults.
func NewBuilder() *Builder {
	now := time.Now().UTC()
	return &Builder{
		entity: &Feedback{
			id:        uuid.New(),
			createdAt: now,
			updatedAt: now,
		},
		validationErrors: make([]error, 0),
	}
}

// BuilderFromExisting creates a builder from an existing feedback entity.
// Useful for update operations (though feedback is immutable, this can be used for reconstruction).
func BuilderFromExisting(f *Feedback) *Builder {
	copied := *f
	copied.updatedAt = time.Now()

	return &Builder{
		entity:           &copied,
		validationErrors: make([]error, 0),
	}
}

// WithID sets the feedback ID.
func (b *Builder) WithID(id uuid.UUID) *Builder {
	if id == uuid.Nil {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("feedback ID cannot be nil"))
		return b
	}
	b.entity.id = id
	return b
}

// WithRating sets the feedback rating with validation.
func (b *Builder) WithRating(rating Rating) *Builder {
	if !rating.IsValid() {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("invalid rating: %d (must be between 1 and 5)", rating))
		return b
	}
	b.entity.rating = rating
	return b
}

// WithRatingValue sets the feedback rating from an integer value with validation.
func (b *Builder) WithRatingValue(ratingValue int) *Builder {
	rating, err := NewRating(ratingValue)
	if err != nil {
		b.validationErrors = append(b.validationErrors, err)
		return b
	}
	b.entity.rating = rating
	return b
}

// WithComment sets the feedback comment with validation.
func (b *Builder) WithComment(comment Comment) *Builder {
	if comment.Value() == "" {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("comment cannot be empty"))
		return b
	}
	b.entity.comment = comment
	return b
}

// WithCommentText sets the feedback comment from text with validation.
func (b *Builder) WithCommentText(commentText string) *Builder {
	comment, err := NewComment(commentText)
	if err != nil {
		b.validationErrors = append(b.validationErrors, err)
		return b
	}
	b.entity.comment = comment
	return b
}

// WithCreatedAt sets the creation timestamp (for database reconstruction).
func (b *Builder) WithCreatedAt(t time.Time) *Builder {
	if t.IsZero() {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("createdAt cannot be zero"))
		return b
	}
	b.entity.createdAt = t
	return b
}

// WithUpdatedAt sets the update timestamp.
func (b *Builder) WithUpdatedAt(t time.Time) *Builder {
	if t.IsZero() {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("updatedAt cannot be zero"))
		return b
	}
	b.entity.updatedAt = t
	return b
}

// WithDeletedAt sets the deletion timestamp (for soft delete reconstruction).
func (b *Builder) WithDeletedAt(deletedAt time.Time) *Builder {
	b.entity.deletedAt = optional.Some(deletedAt)
	return b
}

// Build validates all accumulated data and returns the feedback entity.
// Returns an error if any validation failed or required fields are missing.
func (b *Builder) Build() (*Feedback, error) {
	// Return accumulated validation errors first
	if len(b.validationErrors) > 0 {
		return nil, b.validationErrors[0]
	}

	// Validate required fields
	if b.entity.id == uuid.Nil {
		return nil, fmt.Errorf("feedback ID is required")
	}

	if !b.entity.rating.IsValid() {
		return nil, fmt.Errorf("rating is required and must be between 1 and 5")
	}

	if b.entity.comment.Value() == "" {
		return nil, fmt.Errorf("comment is required")
	}

	if b.entity.createdAt.IsZero() {
		return nil, fmt.Errorf("createdAt timestamp is required")
	}

	if b.entity.updatedAt.IsZero() {
		return nil, fmt.Errorf("updatedAt timestamp is required")
	}

	// Validate the entity itself
	if err := b.entity.IsValid(); err != nil {
		return nil, fmt.Errorf("feedback validation failed: %w", err)
	}

	return b.entity, nil
}

// BuildUnchecked returns the feedback entity without validation.
// ONLY use this for database reconstruction where data integrity is already guaranteed.
func (b *Builder) BuildUnchecked() *Feedback {
	return b.entity
}

// BuildNew creates a new feedback entity with required fields only.
// This is a convenience method for common creation scenarios.
func (b *Builder) BuildNew(ratingValue int, commentText string) (*Feedback, error) {
	return b.
		WithRatingValue(ratingValue).
		WithCommentText(commentText).
		Build()
}
