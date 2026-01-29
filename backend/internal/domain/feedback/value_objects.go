package feedback

import "fmt"

// Rating represents a feedback rating on a scale of 1-5.
type Rating int

const (
	Rating1 Rating = 1
	Rating2 Rating = 2
	Rating3 Rating = 3
	Rating4 Rating = 4
	Rating5 Rating = 5
)

// String returns the string representation of the rating.
func (r Rating) String() string {
	return fmt.Sprintf("%d", r)
}

// IsValid validates that the rating is between 1 and 5.
func (r Rating) IsValid() bool {
	return r >= Rating1 && r <= Rating5
}

// Value returns the integer value of the rating.
func (r Rating) Value() int {
	return int(r)
}

// NewRating creates a new Rating value object with validation.
func NewRating(value int) (Rating, error) {
	rating := Rating(value)
	if !rating.IsValid() {
		return 0, fmt.Errorf("rating must be between 1 and 5, got: %d", value)
	}
	return rating, nil
}

// Comment represents validated feedback comment text.
type Comment struct {
	value string
}

const (
	// MinCommentLength is the minimum length for a comment.
	MinCommentLength = 1
	// MaxCommentLength is the maximum length for a comment.
	MaxCommentLength = 1000
)

// NewComment creates a new Comment value object with validation.
func NewComment(text string) (Comment, error) {
	if text == "" {
		return Comment{}, fmt.Errorf("comment cannot be empty")
	}
	if len(text) < MinCommentLength {
		return Comment{}, fmt.Errorf("comment must be at least %d character(s)", MinCommentLength)
	}
	if len(text) > MaxCommentLength {
		return Comment{}, fmt.Errorf("comment cannot exceed %d characters, got: %d", MaxCommentLength, len(text))
	}
	return Comment{value: text}, nil
}

// Value returns the comment text value.
func (c Comment) Value() string {
	return c.value
}

// String returns the string representation.
func (c Comment) String() string {
	return c.value
}

// Length returns the length of the comment.
func (c Comment) Length() int {
	return len(c.value)
}
