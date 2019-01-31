package loop

import (
	"errors"
	"runtime"
	"sync/atomic"
)

var (
	// ErrNotRunning indicates that the main loop is not running.
	ErrNotRunning = errors.New("main loop is not running")

	// ErrAlreadyRunning indicates that the main loop is already running.
	ErrAlreadyRunning = errors.New("main loop is already running")
)

var (
	isRunning uint32
	lockCount int32
)

// Run locks the OS thread to act as a GUI thread, and then starts the GUI
// event loop until there are no more instances of Window open.
// If the main loop is already running, this function will return an error
// (ErrAlreadyRunning).
//
// Modification of the GUI should happen only on the GUI thread.  This includes
// creating any windows, mounting any widgets, or updating the properties of any
// elements.
//
// The parameter action takes a closure that can be used to initialize the GUI.
// Any further modifications to the GUI also need to be scheduled on the GUI
// thread, which can be done using the function Do.
func Run(action func() error) error {
	// Pin the GUI message loop to a single thread
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// Want to gate entry into the GUI loop so that only one thread may enter
	// at a time.  Since this is supposed to be non-blocking, we can't use
	// a sync.Mutex without a TryLock method.
	if !atomic.CompareAndSwapUint32(&isRunning, 0, 1) {
		return ErrAlreadyRunning
	}
	defer func() {
		atomic.StoreUint32(&isRunning, 0)
	}()

	// Platform specific initialization ahead of running any user actions.
	err := initRun()
	if err != nil {
		return err
	}
	defer terminateRun()

	// Since we have now locked the OS thread, we can call the initial action.
	// We want to hold a reference to a virtual window by increasing the
	// count to prevent a premature exit if any windows are created and then
	// destroyed during the initialization.  To handle any panics, the call
	// to action needs to be wrapped in a function.
	err = func(action func() error) error {
		atomic.AddInt32(&lockCount, 1)
		defer func() {
			atomic.AddInt32(&lockCount, -1)
		}()
		return action()
	}(action)
	if err != nil {
		return err
	}

	// Check that there is at least on top-level window still open.  Otherwise,
	// there is not point in running the GUI event loop.
	if c := atomic.LoadInt32(&lockCount); c <= 0 {
		return nil
	}

	// Defer to platform-specific code.
	run()
	return nil
}

// Do runs the passed function on the GUI thread.  If the GUI event loop is not
// running, this function will return an error (ErrNotRunning).  Any error from
// the callback will also be returned.
//
// Because this function involves asynchronous communication with the GUI thread,
// it can deadlock if called from the GUI thread.  It is therefore not safe to
// use in any event callbacks from widgets.  However, since those callbacks are
// already executing on the GUI thread, the use of Do is also unnecessary in
// that context.
//
// Note, this function contains a race-condition.  An action may be
// scheduled while the event loop is being terminated, in which case the
// scheduled action may never be run.  Presumably, those actions don't need to
// be run on the GUI thread, so they should be scheduled using a different
// mechanism.
//
// If the passed function panics, the panic will be recovered, and wrapped into
// an error.  That error will be used to create a new panic within the
// caller's goroutine.  If the program terminates because of that panic, there
// will be two active goroutines in the stack trace.  One active goroutine will
// be the GUI thread, where the panic originated, and a secton active goroutine
// from caller's goroutine.
func Do(action func() error) error {
	// Check if the event loop is current running.
	if atomic.LoadUint32(&isRunning) == 0 {
		return ErrNotRunning
	}

	// Race-condition here!  Event loop may terminate between previous check
	// and following call, which will block.

	// Defer to platform-specific code.
	return do(action)
}

// AddLockCount is used to track the number of top-level GUI elements that are
// created.  When the count falls back to zero, the event loop will terminate.
//
// Users should not typically need to call this function.  Top-level GUI
// elements, such as windows, will increment and decrement the count as they
// are created and destroyed.
//
// If the GUI event loop is not running, this function will panic.
func AddLockCount(delta int32) {
	// Check if the event loop is current running.
	if atomic.LoadUint32(&isRunning) == 0 {
		panic(ErrNotRunning)
	}

	// Update the lock count.
	if newval := atomic.AddInt32(&lockCount, delta); newval == 0 {
		stop()
	}
}
