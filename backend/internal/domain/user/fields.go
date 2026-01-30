package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/ktruedat/llm-feedback-analysis/pkg/optional"
)

// ID returns the user ID.
func (u *User) ID() uuid.UUID {
	return u.id
}

// Email returns the user email.
func (u *User) Email() Email {
	return u.email
}

// PasswordHash returns the user password hash.
func (u *User) PasswordHash() PasswordHash {
	return u.password
}

// Roles returns a copy of the user's roles.
func (u *User) Roles() []Role {
	roles := make([]Role, len(u.roles))
	copy(roles, u.roles)
	return roles
}

// Status returns the user status.
func (u *User) Status() UserStatus {
	return u.status
}

// CreatedAt returns the creation timestamp.
func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

// UpdatedAt returns the last update timestamp.
func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

// DeletedAt returns the deletion timestamp if deleted.
func (u *User) DeletedAt() optional.Optional[time.Time] {
	return u.deletedAt
}
