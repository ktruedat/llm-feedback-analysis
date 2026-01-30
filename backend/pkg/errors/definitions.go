package errors

type ErrorOpt func(*GenericError)

func WithCauseError(cause error) ErrorOpt {
	return func(e *GenericError) {
		e.Cause = cause
	}
}

func ErrInternal(err error) ApplicationError {
	return &GenericError{
		Code:    ErrorCodeInternal,
		Message: "An internal error occurred",
		Cause:   err,
	}
}

func ErrBadRequest(msg string, opts ...ErrorOpt) ApplicationError {
	ge := &GenericError{
		Code:       ErrorCodeBadRequest,
		Message:    "Bad request: " + msg,
		UserFacing: true,
	}
	for _, opt := range opts {
		opt(ge)
	}

	return ge
}

var ErrNotFound ApplicationError = &GenericError{
	Code:    ErrorCodeNotFound,
	Message: "The requested resource was not found",
}

func ErrNotImplemented(msg string) ApplicationError {
	return &GenericError{
		Code:       ErrorCodeNotImplemented,
		Message:    "Not implemented: " + msg,
		UserFacing: true,
	}
}

func ErrUnauthorized(msg string, opts ...ErrorOpt) ApplicationError {
	ge := &GenericError{
		Code:       ErrorCodeUnauthorized,
		Message:    "Unauthorized: " + msg,
		UserFacing: true,
	}
	for _, opt := range opts {
		opt(ge)
	}

	return ge
}

func ErrForbidden(msg string, opts ...ErrorOpt) ApplicationError {
	ge := &GenericError{
		Code:       ErrorCodeForbidden,
		Message:    "Forbidden: " + msg,
		UserFacing: true,
	}
	for _, opt := range opts {
		opt(ge)
	}

	return ge
}
