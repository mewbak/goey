package loop

import (
	"sync/atomic"

	"bitbucket.org/rj/goey/internal/nopanic"
	"bitbucket.org/rj/goey/internal/syscall"
	"github.com/gotk3/gotk3/gtk"
)

var (
	runLevel uint32
)

func init() {
	gtk.Init(nil)
}

func initRun() error {
	// Do nothing
	return nil
}

func terminateRun() {
	// Do nothing
}

func run() {
	// Handle run level.
	if !atomic.CompareAndSwapUint32(&runLevel, 0, 1) {
		panic("internal error")
	}
	defer atomic.StoreUint32(&runLevel, 0)

	// Start the GTK loop.
	gtk.Main()
}

func do(action func() error) error {
	// Make channel for the return value of the action.
	err := make(chan error, 1)

	// Depending on the run level for the main loop, either use an idle
	// callback or a higher priority callback.  The goal with using an
	// idle callback is to ensure that the system is up and running
	// before any new changes.
	if atomic.LoadUint32(&runLevel) < 2 {
		syscall.IdleAdd(func() {
			atomic.StoreUint32(&runLevel, 2)
			err <- nopanic.Wrap(action)
		})
	} else {
		syscall.MainContextInvoke(func() {
			err <- nopanic.Wrap(action)
		})
	}

	// Block on completion of action.
	return nopanic.Unwrap(<-err)
}

func stop() {
	gtk.MainQuit()
}
