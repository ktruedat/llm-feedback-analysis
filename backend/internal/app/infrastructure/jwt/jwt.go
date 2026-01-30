package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/ktruedat/llm-feedback-analysis/internal/app/config"
)

const (
	defaultExpirationHours = 24
	defaultAlgorithm       = "HS256"
)

// Claims represents the JWT claims structure used in the application.
// It includes standard registered claims plus custom application claims.
type Claims struct {
	jwt.RegisteredClaims
	// UserID is the unique identifier of the user (stored as subject in standard claim)
	// This is a convenience field that maps to Subject.
	UserID string `json:"user_id,omitempty"`
	// Email is the user's email address.
	Email string `json:"email,omitempty"`
	// Roles contains the user's roles/permissions.
	// Example values: ["admin"], ["user", "admin"], ["user"], etc.
	Roles []string `json:"roles,omitempty"`
}

// NewClaims creates a new Claims struct with the provided user information.
// The expiration time is calculated from the JWT config's ExpirationHours.
func NewClaims(userID uuid.UUID, email string, roles []string, cfg *config.JWT) *Claims {
	expirationHours := cfg.ExpirationHours
	if expirationHours <= 0 {
		expirationHours = defaultExpirationHours
	}
	expirationTime := time.Now().Add(time.Duration(expirationHours) * time.Hour)

	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: userID.String(),
		Email:  email,
		Roles:  roles,
	}

	return claims
}

// GenerateToken generates a JWT token with the provided claims.
// Returns the signed token string.
func GenerateToken(claims *Claims, cfg *config.JWT) (string, error) {
	algorithm := cfg.Algorithm
	if algorithm == "" {
		algorithm = defaultAlgorithm
	}

	signingMethod := jwt.GetSigningMethod(algorithm)
	if signingMethod == nil {
		return "", fmt.Errorf("unsupported signing method: %s", algorithm)
	}

	// create and sign the token
	token := jwt.NewWithClaims(signingMethod, claims)
	tokenString, err := token.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ParseToken parses and validates a JWT token string.
// Returns the claims if the token is valid, otherwise returns an error.
func ParseToken(tokenString string, cfg *config.JWT) (*Claims, error) {
	algorithm := cfg.Algorithm
	if algorithm == "" {
		algorithm = defaultAlgorithm
	}

	// Parse token with claims
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
		tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Validate signing algorithm
			expectedAlg := jwt.GetSigningMethod(algorithm)
			if expectedAlg == nil {
				return nil, jwt.ErrSignatureInvalid
			}

			if token.Method != expectedAlg {
				return nil, jwt.ErrSignatureInvalid
			}

			// Return the secret key for validation
			return []byte(cfg.Secret), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Ensure UserID is set from Subject if not explicitly set
	if claims.UserID == "" && claims.Subject != "" {
		claims.UserID = claims.Subject
	}

	return claims, nil
}

// ExtractBearerToken extracts the bearer token from the Authorization header.
// Returns the token string if found, otherwise returns an error.
func ExtractBearerToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", fmt.Errorf("missing authorization header")
	}

	// Validate Bearer token format: "Bearer <token>"
	// Find the first space to split prefix and token
	spaceIndex := -1
	for i := 0; i < len(authHeader); i++ {
		if authHeader[i] == ' ' {
			spaceIndex = i
			break
		}
	}

	if spaceIndex == -1 {
		return "", fmt.Errorf("invalid authorization header format. Expected: Bearer <token>")
	}

	prefix := authHeader[:spaceIndex]
	tokenString := authHeader[spaceIndex+1:]

	if prefix != "Bearer" {
		return "", fmt.Errorf("invalid authorization header format. Expected: Bearer <token>")
	}

	if tokenString == "" {
		return "", fmt.Errorf("bearer token is empty")
	}

	return tokenString, nil
}
