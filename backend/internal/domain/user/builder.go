package user

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ktruedat/llm-feedback-analysis/pkg/optional"
)

// Builder provides type-safe, fluent construction of User entities.
// It accumulates validation errors and applies business rules at Build() time.
type Builder struct {
	entity           *User
	validationErrors []error
}

// NewBuilder creates a builder for creating new user entities.
// Initialize with sensible defaults.
func NewBuilder() *Builder {
	now := time.Now().UTC()
	return &Builder{
		entity: &User{
			id:        uuid.New(),
			roles:     []Role{RoleUser}, // Default role
			status:    UserStatusActive, // Default status
			createdAt: now,
			updatedAt: now,
		},
		validationErrors: make([]error, 0),
	}
}

// BuilderFromExisting creates a builder from an existing user entity.
// Useful for update operations.
func BuilderFromExisting(u *User) *Builder {
	copied := *u
	// Deep copy the roles slice
	copied.roles = make([]Role, len(u.roles))
	copy(copied.roles, u.roles)
	copied.updatedAt = time.Now().UTC()

	return &Builder{
		entity:           &copied,
		validationErrors: make([]error, 0),
	}
}

// WithID sets the user ID.
func (b *Builder) WithID(id uuid.UUID) *Builder {
	if id == uuid.Nil {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("user ID cannot be nil"))
		return b
	}
	b.entity.id = id
	return b
}

// WithEmail sets the user email with validation.
func (b *Builder) WithEmail(email Email) *Builder {
	if err := email.IsValid(); err != nil {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("invalid email: %w", err))
		return b
	}
	b.entity.email = email
	return b
}

// WithEmailString sets the user email from a string with validation.
func (b *Builder) WithEmailString(emailStr string) *Builder {
	email, err := NewEmail(emailStr)
	if err != nil {
		b.validationErrors = append(b.validationErrors, err)
		return b
	}
	b.entity.email = email
	return b
}

// WithPasswordHash sets the user password hash.
func (b *Builder) WithPasswordHash(passwordHash PasswordHash) *Builder {
	if passwordHash.Value() == "" {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("password hash cannot be empty"))
		return b
	}
	b.entity.password = passwordHash
	return b
}

// WithPasswordHashString sets the user password hash from a string.
func (b *Builder) WithPasswordHashString(hashStr string) *Builder {
	hash, err := NewPasswordHash(hashStr)
	if err != nil {
		b.validationErrors = append(b.validationErrors, err)
		return b
	}
	b.entity.password = hash
	return b
}

// WithRoles sets the user roles with validation.
func (b *Builder) WithRoles(roles []Role) *Builder {
	if len(roles) == 0 {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("user must have at least one role"))
		return b
	}

	// Validate each role
	for _, role := range roles {
		if !role.IsValid() {
			b.validationErrors = append(b.validationErrors, fmt.Errorf("invalid role: %s", role))
			return b
		}
	}

	// Deep copy the roles slice
	b.entity.roles = make([]Role, len(roles))
	copy(b.entity.roles, roles)
	return b
}

// WithRole adds a single role to the user (replaces existing roles).
func (b *Builder) WithRole(role Role) *Builder {
	if !role.IsValid() {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("invalid role: %s", role))
		return b
	}
	b.entity.roles = []Role{role}
	return b
}

// WithRoleString sets a single role from a string.
func (b *Builder) WithRoleString(roleStr string) *Builder {
	role, err := NewRole(roleStr)
	if err != nil {
		b.validationErrors = append(b.validationErrors, err)
		return b
	}
	b.entity.roles = []Role{role}
	return b
}

// WithStatus sets the user status with validation.
func (b *Builder) WithStatus(status UserStatus) *Builder {
	if !status.IsValid() {
		b.validationErrors = append(b.validationErrors, fmt.Errorf("invalid user status: %s", status))
		return b
	}
	b.entity.status = status
	return b
}

// WithStatusString sets the user status from a string with validation.
func (b *Builder) WithStatusString(statusStr string) *Builder {
	status, err := NewUserStatus(statusStr)
	if err != nil {
		b.validationErrors = append(b.validationErrors, err)
		return b
	}
	b.entity.status = status
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

// Build validates all accumulated data and returns the user entity.
// Returns an error if any validation failed or required fields are missing.
func (b *Builder) Build() (*User, error) {
	// Return accumulated validation errors first
	if len(b.validationErrors) > 0 {
		return nil, b.validationErrors[0]
	}

	// Validate required fields
	if b.entity.id == uuid.Nil {
		return nil, fmt.Errorf("user ID is required")
	}

	if err := b.entity.email.IsValid(); err != nil {
		return nil, fmt.Errorf("email is required and must be valid: %w", err)
	}

	if b.entity.password.Value() == "" {
		return nil, fmt.Errorf("password hash is required")
	}

	if len(b.entity.roles) == 0 {
		return nil, fmt.Errorf("user must have at least one role")
	}

	if !b.entity.status.IsValid() {
		return nil, fmt.Errorf("user status is required and must be valid")
	}

	if b.entity.createdAt.IsZero() {
		return nil, fmt.Errorf("createdAt timestamp is required")
	}

	if b.entity.updatedAt.IsZero() {
		return nil, fmt.Errorf("updatedAt timestamp is required")
	}

	// Validate the entity itself
	if err := b.entity.IsValid(); err != nil {
		return nil, fmt.Errorf("user validation failed: %w", err)
	}

	return b.entity, nil
}

// BuildUnchecked returns the user entity without validation.
// ONLY use this for database reconstruction where data integrity is already guaranteed.
func (b *Builder) BuildUnchecked() *User {
	return b.entity
}

// BuildNew creates a new user entity with required fields only.
// This is a convenience method for common creation scenarios.
func (b *Builder) BuildNew(emailStr string, passwordHashStr string) (*User, error) {
	return b.
		WithEmailString(emailStr).
		WithPasswordHashString(passwordHashStr).
		Build()
}
