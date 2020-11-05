package wraperror

import (
	"errors"
)

// wrapError is custom error struct implemented error interface, supporting As, Is, Unwrap.
type wrapError struct {
	current error
	child   error
}

// Error converts error to WrapError.
func Error(err error) *wrapError {
	if err == nil {
		return &wrapError{}
	}

	switch err.(type) {
	case *wrapError:
		return err.(*wrapError)
	default:
		return &wrapError{current: err}
	}
}

// Error returns a chained error string.
// It's implemented error interface.
func (e *wrapError) Error() string {
	if e.current == nil {
		return ""
	}
	msg := e.current.Error()
	if e.child != nil {
		msg += " " + e.child.Error()
	}
	return msg
}

// Wrap wraps an error.
func (e *wrapError) Wrap(err error) *wrapError {
	return &wrapError{current: err, child: e}
}

// Unwrap unwraps an error.
func (e *wrapError) Unwrap() error {
	return e.child
}

// Is checks to equal e.current to err.
func (e *wrapError) Is(target error) bool {
	return errors.Is(e.current, target)
}

// As Checks to equal e.child to err.
func (e *wrapError) As(target interface{}) bool {
	return errors.As(e.current, target)
}
