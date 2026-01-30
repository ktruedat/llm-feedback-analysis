package requests

// RegisterUserRequest represents the request payload for user registration
//
//	@Description	Request payload for registering a new user account.
type RegisterUserRequest struct {
	Email    string `json:"email" example:"user@example.com" binding:"required,email"`      // User email address (required, must be valid email)
	Password string `json:"password" example:"SecurePassword123!" binding:"required,min=8"` // User password (required, minimum 8 characters)
}

// LoginUserRequest represents the request payload for user authentication
//
//	@Description	Request payload for user login/authentication.
type LoginUserRequest struct {
	Email    string `json:"email" example:"user@example.com" binding:"required,email"` // User email address (required)
	Password string `json:"password" example:"SecurePassword123!" binding:"required"`  // User password (required)
}
