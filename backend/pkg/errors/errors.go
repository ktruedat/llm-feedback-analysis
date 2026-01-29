package errors

// ApplicationError represents an application-level error that can be handled
// by the infrastructure layer (HTTP, gRPC, etc.)
type ApplicationError interface {
	error
	httpResponse

	ErrCode() *ErrorCode // Code represents the unique error code.
	ErrMessage() string  // Message provides a human-readable description of the error.
	ErrCause() error     // Cause holds the underlying error, if any.
	IsUserFacing() bool  // IsUserFacing indicates if the error message is safe to show to end-users.
}

type httpResponse interface {
	IsResponse()
	HTTPCode() int
}

var _ ApplicationError = (*GenericError)(nil)

// GenericError is a generic implementation of ApplicationError.
type GenericError struct {
	Code       *ErrorCode `json:"code"`
	Message    string     `json:"message"`
	Cause      error      `json:"-"`
	UserFacing bool       `json:"-"`
}

func (e *GenericError) Error() string {
	return e.Message
}

func (e *GenericError) Unwrap() error {
	return e.Cause
}

func (e *GenericError) ErrCode() *ErrorCode {
	return e.Code
}

func (e *GenericError) ErrMessage() string {
	return e.Message
}

func (e *GenericError) ErrCause() error {
	return e.Cause
}

func (e *GenericError) IsUserFacing() bool {
	return e.UserFacing
}

func (*GenericError) IsResponse() {}

func (e *GenericError) HTTPCode() int {
	return e.Code.Category.HTTPCode()
}
