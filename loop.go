package goey

import (
	"errors"
)

var (
	// ErrQuit indicates that the event loop should terminate.  This return
	// will only be used on platforms that expose their loop iteration function
	// in addition to Run.
	ErrQuit = errors.New("quit message")

	// ErrNotRunning indicates that the main loop is not running.
	ErrNotRunning = errors.New("main loop is not running")

	// ErrAlreadyRunning indicates that the main loop is not running.
	ErrAlreadyRunning = errors.New("main loop is already running")
)

// Run iterates over the event loop until there are no more instances of
// MainWindow open.
// If the main loop is already running, this function will return an error.
func Run() error {
	// Defer to platform-specific code.
	return run()
}

// Do runs the passed function on the GUI thread.  If the main loop is not
// running, it will return an error.  Any error from the callback will
// also be returned.
func Do(action func() error) error {
	// Defer to platform-specific code.
	return do(action)
}
