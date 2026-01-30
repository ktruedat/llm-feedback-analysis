package user

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ktruedat/llm-feedback-analysis/pkg/optional"
)

// User represents a user account in the system.
//
// Business Rules:
// - Email must be unique and valid (enforced by Email value object)
// - Password must be hashed (never store plain text)
// - User must have at least one role
// - Cannot be edited when deleted
// - Email cannot be changed after creation (for audit purposes)
//
// Relationships:
// - One user can have multiple feedbacks.
type User struct {
	id        uuid.UUID
	email     Email
	password  PasswordHash // Always hashed, never plain text
	roles     []Role       // User can have multiple roles
	status    UserStatus
	createdAt time.Time
	updatedAt time.Time
	deletedAt optional.Optional[time.Time]
}

// IsValid validates the entire user entity state.
func (u *User) IsValid() error {
	if u.id == uuid.Nil {
		return fmt.Errorf("user ID is required")
	}

	if err := u.email.IsValid(); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}

	if u.password.Value() == "" {
		return fmt.Errorf("password hash is required")
	}

	if len(u.roles) == 0 {
		return fmt.Errorf("user must have at least one role")
	}

	if !u.status.IsValid() {
		return fmt.Errorf("invalid user status: %s", u.status)
	}

	if u.createdAt.IsZero() {
		return fmt.Errorf("createdAt timestamp is required")
	}

	if u.updatedAt.IsZero() {
		return fmt.Errorf("updatedAt timestamp is required")
	}

	return nil
}

// CanBeEdited returns true if the user can be edited.
func (u *User) CanBeEdited() bool {
	return !u.IsDeleted() && u.status != UserStatusSuspended
}

// CanBeDeleted returns true if the user can be deleted.
func (u *User) CanBeDeleted() bool {
	return !u.IsDeleted()
}

// IsDeleted returns true if the user is soft-deleted.
func (u *User) IsDeleted() bool {
	return u.deletedAt.IsSome()
}

// IsActive returns true if the user is active and not deleted.
func (u *User) IsActive() bool {
	return !u.IsDeleted() && u.status == UserStatusActive
}

// HasRole returns true if the user has the specified role.
func (u *User) HasRole(role Role) bool {
	for _, r := range u.roles {
		if r == role {
			return true
		}
	}
	return false
}

// IsAdmin returns true if the user has admin role.
func (u *User) IsAdmin() bool {
	return u.HasRole(RoleAdmin)
}

// ChangePassword updates the user's password hash.
// This should be called with an already-hashed password from the application layer.
func (u *User) ChangePassword(newPasswordHash PasswordHash) error {
	if !u.CanBeEdited() {
		return fmt.Errorf("user cannot be edited in current state")
	}

	if newPasswordHash.Value() == "" {
		return fmt.Errorf("password hash cannot be empty")
	}

	u.password = newPasswordHash
	u.updatedAt = time.Now()
	return nil
}

// AddRole adds a role to the user if not already present.
func (u *User) AddRole(role Role) error {
	if !u.CanBeEdited() {
		return fmt.Errorf("user cannot be edited in current state")
	}

	if !role.IsValid() {
		return fmt.Errorf("invalid role: %s", role)
	}

	// Check if role already exists
	if u.HasRole(role) {
		return nil // Already has the role, no error
	}

	u.roles = append(u.roles, role)
	u.updatedAt = time.Now()
	return nil
}

// RemoveRole removes a role from the user.
func (u *User) RemoveRole(role Role) error {
	if !u.CanBeEdited() {
		return fmt.Errorf("user cannot be edited in current state")
	}

	// Prevent removing all roles
	if len(u.roles) == 1 && u.HasRole(role) {
		return fmt.Errorf("cannot remove the last role from user")
	}

	// Remove role from slice
	newRoles := make([]Role, 0, len(u.roles))
	for _, r := range u.roles {
		if r != role {
			newRoles = append(newRoles, r)
		}
	}

	u.roles = newRoles
	u.updatedAt = time.Now()
	return nil
}

// ChangeStatus changes the user's status with validation.
func (u *User) ChangeStatus(newStatus UserStatus) error {
	if !newStatus.IsValid() {
		return fmt.Errorf("invalid status: %s", newStatus)
	}

	if u.status == newStatus {
		return nil // No change needed
	}

	// Business rule: Cannot change status of deleted user
	if u.IsDeleted() {
		return fmt.Errorf("cannot change status of deleted user")
	}

	u.status = newStatus
	u.updatedAt = time.Now()
	return nil
}

// Activate activates the user account.
func (u *User) Activate() error {
	return u.ChangeStatus(UserStatusActive)
}

// Suspend suspends the user account.
func (u *User) Suspend() error {
	return u.ChangeStatus(UserStatusSuspended)
}

// Delete performs soft delete on the user.
func (u *User) Delete() error {
	if u.IsDeleted() {
		return fmt.Errorf("user is already deleted")
	}

	now := time.Now()
	u.deletedAt = optional.Some(now)
	u.status = UserStatusInactive // Set status to inactive when deleted
	u.updatedAt = now
	return nil
}

// Restore restores a soft-deleted user.
func (u *User) Restore() error {
	if !u.IsDeleted() {
		return fmt.Errorf("user is not deleted")
	}

	u.deletedAt = optional.None[time.Time]()
	u.status = UserStatusActive // Restore to active status
	u.updatedAt = time.Now()
	return nil
}
