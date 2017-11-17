package goey

import (
	"errors"
)

var (
	// ErrQuit indicates that the event loop should terminate.  This return
	// will only be used on platforms that expose their loop iteration function
	// in addition to Run.
	ErrQuit = errors.New("quit message")
)

// Run iterates over the event loop until there are no more MainWindows open.
func Run() error {
	return run()
}
