package errors

import (
	"encoding/json"
)

// ErrorCategory represents semantic error categories that are transport-agnostic.
// These categories can be mapped to appropriate HTTP status codes, gRPC codes, etc.
// in the infrastructure/transport layer.
type ErrorCategory string

const (
	CategoryValidation ErrorCategory = "ValidationError"
	CategoryNotFound   ErrorCategory = "NotFoundError"
	CategoryConflict   ErrorCategory = "ConflictError"
	CategoryInternal   ErrorCategory = "InternalError"
)

func (c ErrorCategory) HTTPCode() int {
	switch c {
	case CategoryValidation:
		return 400 // Bad Request
	case CategoryNotFound:
		return 404 // Not Found
	case CategoryConflict:
		return 409 // Conflict
	case CategoryInternal:
		return 500 // Internal Server Error
	default:
		return 500 // Default to Internal Server Error
	}
}

// ErrorCode represents an error code with a semantic category.
// Domain packages should use this struct to create their own error codes.
// The error codes should be written in snake_case format.
type ErrorCode struct {
	Code     string        `json:"-"`
	Category ErrorCategory `json:"-"`
}

// MarshalJSON implements json.Marshaler to serialize only the Code string value.
func (e *ErrorCode) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Code)
}

// NewDomainErrorCode creates a new domain-specific error code with a semantic category.
// Domain packages should use this to create their own error code types.
// The category should be one of: CategoryValidation, CategoryNotFound, CategoryConflict, or CategoryInternal.
func NewDomainErrorCode(code string, category ErrorCategory) *ErrorCode {
	return &ErrorCode{
		Code:     code,
		Category: category,
	}
}

var (
	ErrorCodeInternal = &ErrorCode{
		Code:     "internal_error",
		Category: CategoryInternal,
	}
	ErrorCodeBadRequest = &ErrorCode{
		Code:     "bad_request",
		Category: CategoryValidation,
	}
	ErrorCodeNotFound = &ErrorCode{
		Code:     "not_found",
		Category: CategoryNotFound,
	}
	ErrorCodeNotImplemented = &ErrorCode{
		Code:     "not_implemented",
		Category: CategoryInternal,
	}
)
