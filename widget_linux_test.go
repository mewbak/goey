package goey

import (
	"testing"
	"time"

	"github.com/gotk3/gotk3/gtk"
)

type Handler interface {
	Handle() *gtk.Widget
}

func testingSetFocus(t *testing.T, w *Window, i int) {
	// Check the size
	handle := w.child.(*mountedVBox).children[i].(Handler).Handle()
	if !handle.GetCanFocus() {
		t.Errorf("Widget can not grab focus.")
		return
	}

	handle.GrabFocus()
	time.Sleep(500 * time.Millisecond)
	if !handle.IsFocus() {
		t.Errorf("Widget did not grab focus")
	}

	// Note sure why the call to sleep is required, but there may be a debounce
	// provided by the system.  Without this call to sleep, the controls never
	// get the focus events.
	time.Sleep(250 * time.Millisecond)
}
