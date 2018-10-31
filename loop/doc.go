// Package loop provides a GUI event loop.  The event loop will be locked to
// an OS thread and not return until all GUI elements have been destroyed.
// However, callbacks can be dispatched to the event loop, and will be executed
// on the same thread.
package loop
