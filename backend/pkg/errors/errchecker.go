package errors

import (
	"errors"
)

type ErrorChecker interface {
	Check(err error) error
}

type errChecker struct{}

func NewErrorChecker() ErrorChecker {
	return &errChecker{}
}

func (*errChecker) Check(err error) error {
	var appErr *GenericError
	if errors.As(err, &appErr) {
		if appErr.UserFacing {
			return appErr
		}

		return ErrInternal(err)
	}

	return ErrInternal(err)
}

type testErrorChecker struct{}

// NewTestErrorChecker creates an ErrorChecker that returns errors as is, for testing purposes.
func NewTestErrorChecker() ErrorChecker {
	return &testErrorChecker{}
}

func (*testErrorChecker) Check(err error) error {
	return err
}
