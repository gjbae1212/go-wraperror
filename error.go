package wraperror

import (
	"errors"
)

// WrapError is custom error struct implemented error interface, supporting As, Is, Unwrap.
type WrapError struct {
	current error
	child   error
}

// Current returns current error.
func (e *WrapError) Current() error {
	return e.current
}

// Child returns child error.
func (e *WrapError) Child() error {
	return e.child
}

// Wrap wraps an error.
func (e *WrapError) Wrap(err error) *WrapError {
	return &WrapError{current: err, child: e}
}

// Flatten returns an flatten error array.
func (e *WrapError) Flatten() []error {
	var errs []error

	if e.current != nil {
		if _, ok := e.current.(*WrapError); ok {
			errs = append(errs, e.current.(*WrapError).Flatten()...)
		} else {
			errs = append(errs, e.current)
			if unwrap := errors.Unwrap(e.current); unwrap != nil {
				errs = append(errs, Error(unwrap).Flatten()...)
			}
		}
	}

	if e.child != nil {
		if _, ok := e.child.(*WrapError); ok {
			errs = append(errs, e.child.(*WrapError).Flatten()...)
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
func (e *WrapError) Error() string {
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
func (e *WrapError) Unwrap() error {
	return e.child
}

// Is checks to equal e.current to err.
// It's implemented for errors package.
func (e *WrapError) Is(target error) bool {
	return errors.Is(e.current, target)
}

// As Checks to equal e.child to err.
// It's implemented for errors package.
func (e *WrapError) As(target interface{}) bool {
	return errors.As(e.current, target)
}

// Error converts error to WrapError.
func Error(err error) *WrapError {
	if err == nil {
		return &WrapError{}
	}

	switch err.(type) {
	case *WrapError:
		return err.(*WrapError)
	default:
		return &WrapError{current: err}
	}
}

// FromError returns an error to WrapError,
// but if an error doesn't convert to WrapError, returning nil and false.
func FromError(err error) (*WrapError, bool) {
	if err == nil {
		return nil, false
	}
	wrapErr, ok := err.(*WrapError)
	return wrapErr, ok
}
