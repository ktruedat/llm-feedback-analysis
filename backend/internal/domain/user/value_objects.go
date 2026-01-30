package user

import (
	"fmt"
	"regexp"
	"strings"
)

// Email represents a validated email address.
type Email struct {
	value string
}

var (
	// emailRegex is a basic email validation regex.
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	// MaxEmailLength is the maximum length for an email address.
	MaxEmailLength = 254 // RFC 5321
)

// NewEmail creates a new Email value object with validation.
func NewEmail(email string) (Email, error) {
	email = strings.TrimSpace(email)
	email = strings.ToLower(email) // Normalize to lowercase

	if email == "" {
		return Email{}, fmt.Errorf("email cannot be empty")
	}

	if len(email) > MaxEmailLength {
		return Email{}, fmt.Errorf("email cannot exceed %d characters, got: %d", MaxEmailLength, len(email))
	}

	if !emailRegex.MatchString(email) {
		return Email{}, fmt.Errorf("invalid email format: %s", email)
	}

	return Email{value: email}, nil
}

// Value returns the email value.
func (e Email) Value() string {
	return e.value
}

// String returns the string representation.
func (e Email) String() string {
	return e.value
}

// IsValid validates the email.
func (e Email) IsValid() error {
	if e.value == "" {
		return fmt.Errorf("email cannot be empty")
	}
	if !emailRegex.MatchString(e.value) {
		return fmt.Errorf("invalid email format: %s", e.value)
	}
	return nil
}

// PasswordHash represents a hashed password (never plain text).
// This is a value object to ensure we never accidentally store plain text passwords.
type PasswordHash struct {
	value string
}

// NewPasswordHash creates a new PasswordHash value object.
// This should only be called with an already-hashed password from the application layer.
func NewPasswordHash(hash string) (PasswordHash, error) {
	if hash == "" {
		return PasswordHash{}, fmt.Errorf("password hash cannot be empty")
	}

	// Basic validation: ensure it looks like a hash (starts with $2a$, $2b$, $argon2, etc.)
	// or is a reasonable length (bcrypt hashes are 60 chars, argon2 can be longer)
	if len(hash) < 20 {
		return PasswordHash{}, fmt.Errorf("password hash appears to be invalid (too short)")
	}

	return PasswordHash{value: hash}, nil
}

// Value returns the password hash value.
func (p PasswordHash) Value() string {
	return p.value
}

// String returns the string representation (masked for security).
func (p PasswordHash) String() string {
	if len(p.value) < 8 {
		return "***"
	}
	return p.value[:4] + "..." + p.value[len(p.value)-4:]
}

// Role represents a user role in the system.
type Role string

const (
	// RoleUser represents a regular user role.
	RoleUser Role = "user"
	// RoleAdmin represents an administrator role.
	RoleAdmin Role = "admin"
)

// String returns the string representation of the role.
func (r Role) String() string {
	return string(r)
}

// IsValid validates that the role is valid.
func (r Role) IsValid() bool {
	switch r {
	case RoleUser, RoleAdmin:
		return true
	default:
		return false
	}
}

// NewRole creates a new Role value object with validation.
func NewRole(role string) (Role, error) {
	r := Role(strings.ToLower(strings.TrimSpace(role)))
	if !r.IsValid() {
		return "", fmt.Errorf("invalid role: %s (valid roles: %s, %s)", role, RoleUser, RoleAdmin)
	}
	return r, nil
}

// UserStatus represents the status of a user account.
type UserStatus string

const (
	// UserStatusActive represents an active user account.
	UserStatusActive UserStatus = "active"
	// UserStatusInactive represents an inactive user account.
	UserStatusInactive UserStatus = "inactive"
	// UserStatusSuspended represents a suspended user account.
	UserStatusSuspended UserStatus = "suspended"
)

// String returns the string representation of the status.
func (s UserStatus) String() string {
	return string(s)
}

// IsValid validates that the status is valid.
func (s UserStatus) IsValid() bool {
	switch s {
	case UserStatusActive, UserStatusInactive, UserStatusSuspended:
		return true
	default:
		return false
	}
}

// NewUserStatus creates a new UserStatus value object with validation.
func NewUserStatus(status string) (UserStatus, error) {
	s := UserStatus(strings.ToLower(strings.TrimSpace(status)))
	if !s.IsValid() {
		return "", fmt.Errorf("invalid user status: %s (valid statuses: %s, %s, %s)",
			status, UserStatusActive, UserStatusInactive, UserStatusSuspended)
	}
	return s, nil
}

// CanLogin returns true if a user with this status can log in.
func (s UserStatus) CanLogin() bool {
	return s == UserStatusActive
}
