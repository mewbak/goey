package nopanic

import (
	"fmt"
	"runtime/debug"
)

// PanicError represents a panic that occurred.
type PanicError struct {
	value interface{}
	stack string
}

// Error returns a description of the error.
func (pe PanicError) Error() string {
	return fmt.Sprintf("%v\n\n%s", pe.value, pe.stack)
}

// Panic panics with itself as the value.  This method exist in case future
// changes to the language allow better tracking of the original stack trace.
func (pe PanicError) Panic() {
	panic(pe)
}

// Value returns the value returned by recover after a panic.
func (pe PanicError) Value() interface{} {
	return pe.value
}

// Stack returns a formatted stack trace of the goroutine that originally
// paniced.
func (pe PanicError) Stack() string {
	return pe.stack
}

// New wraps the value returned by recover into a PanicError.
func New(value interface{}) PanicError {
	stack := string(debug.Stack())
	return PanicError{value, stack}
}

// Wrap ensures that no panics escape.  Action will be called, and if it
// returns normally, its return value will be returned.  However, if action
// panics, the panic will be converted to an error, and will be returned.
func Wrap(action func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = New(r)
		}
	}()

	return action()
}
