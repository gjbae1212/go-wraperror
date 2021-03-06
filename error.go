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

// Wrap wraps an error.
func (e *wrapError) Wrap(err error) *wrapError {
	return &wrapError{current: err, child: e}
}

// Errors returns an flatten error array.
func (e *wrapError) Flatten() []error {
	var errs []error

	if e.current != nil {
		if _, ok := e.current.(*wrapError); ok {
			errs = append(errs, e.current.(*wrapError).Flatten()...)
		} else {
			errs = append(errs, e.current)
			if unwrap := errors.Unwrap(e.current); unwrap != nil {
				errs = append(errs, Error(unwrap).Flatten()...)
			}
		}
	}

	if e.child != nil {
		if _, ok := e.child.(*wrapError); ok {
			errs = append(errs, e.child.(*wrapError).Flatten()...)
		} else {
			errs = append(errs, e.child)
			if unwrap := errors.Unwrap(e.child); unwrap != nil {
				errs = append(errs, Error(unwrap).Flatten()...)
			}
		}
	}
	return errs
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

// Unwrap unwraps an error.
// It's implemented for errors package.
func (e *wrapError) Unwrap() error {
	return e.child
}

// Is checks to equal e.current to err.
// It's implemented for errors package.
func (e *wrapError) Is(target error) bool {
	return errors.Is(e.current, target)
}

// As Checks to equal e.child to err.
// It's implemented for errors package.
func (e *wrapError) As(target interface{}) bool {
	return errors.As(e.current, target)
}
