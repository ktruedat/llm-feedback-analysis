package responses

import (
	"time"

	"github.com/ktruedat/llm-feedback-analysis/internal/domain/user"
)

// RegisterUserResponse represents the response payload for user registration
//
//	@Description	Response payload containing user registration details.
type RegisterUserResponse struct {
	ID        string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"` // User unique identifier
	Email     string    `json:"email" example:"user@example.com"`                  // User email address
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`         // Account creation timestamp
}

// RegisterUserResponseFromDomain converts a domain User entity to a RegisterUserResponse.
func RegisterUserResponseFromDomain(u *user.User) *RegisterUserResponse {
	return &RegisterUserResponse{
		ID:        u.ID().String(),
		Email:     u.Email().Value(),
		CreatedAt: u.CreatedAt(),
	}
}

// LoginUserResponse represents the response payload for user authentication
//
//	@Description	Response payload containing authentication token and user details.
type LoginUserResponse struct {
	Token     string   `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // JWT authentication token
	ExpiresIn int      `json:"expires_in" example:"3600"`                               // Token expiration time in seconds
	User      UserInfo `json:"user"`                                                    // User information including roles
}

// UserInfo represents basic user information in authentication responses
//
//	@Description	Basic user information.
type UserInfo struct {
	ID    string   `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"` // User unique identifier
	Email string   `json:"email" example:"user@example.com"`                  // User email address
	Roles []string `json:"roles" example:"[\"user\"]"`                        // User roles
}

// UserInfoFromDomain converts a domain User entity to a UserInfo.
func UserInfoFromDomain(u *user.User) UserInfo {
	roles := u.Roles()
	roleStrings := make([]string, len(roles))
	for i, role := range roles {
		roleStrings[i] = role.String()
	}

	return UserInfo{
		ID:    u.ID().String(),
		Email: u.Email().Value(),
		Roles: roleStrings,
	}
}
