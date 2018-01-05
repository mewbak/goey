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
	time.Sleep(250 * time.Millisecond)
}
