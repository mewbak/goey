package goey

import (
	"testing"

	"github.com/lxn/win"
)

func testingSetFocus(t *testing.T, w *Window, i int) {
	hwnd := win.GetWindow(w.hWnd, win.GW_CHILD)
	if hwnd == 0 {
		t.Errorf("Internal error to testing, failure in GetWindow")
		return
	}
	for i := i; i > 0; i-- {
		hwnd = win.GetWindow(hwnd, win.GW_HWNDNEXT)
	}
	if hwnd == 0 {
		t.Errorf("Internal error to testing, failure in GetWindow")
		return
	}

	// When starting, the first control may have already been given focus
	// by the main window.  We don't want to double up on setting the focus.
	if win.GetFocus() != hwnd {
		win.SetFocus(hwnd)
	}
}
