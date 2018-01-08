package goey

import (
	"testing"
	"time"
)

func testingSetFocus(t *testing.T, w *Window, i int) {
	// Check the size
	handle := w.vbox.children[i].Handle()

	handle.GrabFocus()
	if !handle.IsFocus() {
		t.Errorf("Wedgit did not grab focus")
	}

	// Note sure why the call to sleep is required, but there may be a debounce
	// provided by the system.  Without this call to sleep, the controls never
	// get the focus events.
	time.Sleep(250 * time.Millisecond)
}
