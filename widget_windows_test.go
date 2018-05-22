package goey

import (
	"github.com/lxn/win"
	"testing"
)

func TestCalculateHGap(t *testing.T) {
	cases := []struct {
		w1, w2   MountedWidget
		expected Length
	}{
		{(*mountedTextInput)(nil), (*mountedTextInput)(nil), 11}, // Space between unrelated controls
		{(*mountedTextInput)(nil), (*mountedButton)(nil), 11},    // Space between unrelated controls
		{(*mountedButton)(nil), (*mountedTextInput)(nil), 11},    // Space between unrelated controls
		{(*mountedButton)(nil), (*mountedButton)(nil), 7},        // Space between adjacent buttons
	}

	for _, v := range cases {
		out := calculateHGap(v.w1, v.w2)
		if out != v.expected {
			t.Errorf("Incorrect horizontal gap calculated, %d =/= %d", out, v.expected)
		}
	}
}

func TestCalculateVGap(t *testing.T) {
	cases := []struct {
		w1, w2   MountedWidget
		expected Length
	}{
		{(*mountedTextInput)(nil), (*mountedTextInput)(nil), 11},   // Space between unrelated controls
		{(*mountedTextInput)(nil), (*mountedP)(nil), 11},           // Space between unrelated controls
		{(*mountedTextInput)(nil), (*mountedSelectInput)(nil), 11}, // Space between unrelated controls
		{(*mountedLabel)(nil), (*mountedTextInput)(nil), 5},        // Space between text labels and associated fields
		{(*mountedLabel)(nil), (*mountedSelectInput)(nil), 5},      // Space between text labels and associated fields
		{(*mountedLabel)(nil), (*mountedTextArea)(nil), 5},         // Space between text labels and associated fields
		{(*mountedCheckbox)(nil), (*mountedCheckbox)(nil), 7},      // Space between related controls
		{(*mountedP)(nil), (*mountedP)(nil), 11},                   // Space between paragraphs of text
	}

	for _, v := range cases {
		out := calculateVGap(v.w1, v.w2)
		if out != v.expected {
			t.Errorf("Incorrect vertical gap calculated, %d =/= %d", out, v.expected)
		}
	}
}

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
