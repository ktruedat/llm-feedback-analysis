package feedback

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ktruedat/llm-feedback-analysis/pkg/optional"
)

// Feedback represents a user feedback submission.
//
// Business Rules:
// - Rating must be between 1 and 5 (enforced by Rating value object)
// - Comment is required and must be between 1 and 1000 characters
// - Cannot be edited once created (immutable after creation)
// - Can be soft-deleted
//
// Relationships:
// - Standalone entity (no parent/child relationships).
type Feedback struct {
	id        uuid.UUID
	rating    Rating
	comment   Comment
	createdAt time.Time
	updatedAt time.Time
	deletedAt optional.Optional[time.Time]
}

// IsValid validates the entire feedback entity state.
func (f *Feedback) IsValid() error {
	if f.id == uuid.Nil {
		return fmt.Errorf("feedback ID is required")
	}

	if !f.rating.IsValid() {
		return fmt.Errorf("invalid rating: %d", f.rating)
	}

	if f.comment.Value() == "" {
		return fmt.Errorf("comment is required")
	}

	if f.createdAt.IsZero() {
		return fmt.Errorf("createdAt timestamp is required")
	}

	if f.updatedAt.IsZero() {
		return fmt.Errorf("updatedAt timestamp is required")
	}

	return nil
}

// CanBeDeleted returns true if the feedback can be deleted.
// All feedback can be soft-deleted regardless of state.
func (f *Feedback) CanBeDeleted() bool {
	return !f.IsDeleted()
}

// IsDeleted returns true if the feedback is soft-deleted.
func (f *Feedback) IsDeleted() bool {
	return f.deletedAt.IsSome()
}

// Delete performs soft delete on the feedback.
func (f *Feedback) Delete() error {
	if f.IsDeleted() {
		return fmt.Errorf("feedback is already deleted")
	}

	now := time.Now()
	f.deletedAt = optional.Some(now)
	f.updatedAt = now
	return nil
}

// Restore restores a soft-deleted feedback.
func (f *Feedback) Restore() error {
	if !f.IsDeleted() {
		return fmt.Errorf("feedback is not deleted")
	}

	f.deletedAt = optional.None[time.Time]()
	f.updatedAt = time.Now()
	return nil
}
