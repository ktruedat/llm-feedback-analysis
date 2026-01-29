package responder

import (
	"encoding/json"
)

type ErrorResponse struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (*ErrorResponse) IsResponse() {}

func (r *ErrorResponse) HTTPCode() int {
	return r.Code
}

func NewErrorResponse(statusCode int, message string) Response {
	return &ErrorResponse{
		Status:  "error",
		Code:    statusCode,
		Message: message,
	}
}

type GenericResponse struct {
	Payload    any `json:"-"`
	StatusCode int
}

func (*GenericResponse) IsResponse() {}

func (r *GenericResponse) HTTPCode() int {
	return r.StatusCode
}

func (r *GenericResponse) MarshalJSON() ([]byte, error) {
	if r.Payload == nil {
		return []byte("null"), nil
	}
	return json.Marshal(r.Payload)
}

func NewGenericResponse(statusCode int, payload any) Response {
	return &GenericResponse{
		Payload:    payload,
		StatusCode: statusCode,
	}
}
