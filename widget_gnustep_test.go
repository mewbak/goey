// +build gnustep

package goey

import (
	"testing"
)

func testingClick(t *testing.T, w *Window, i int) {
	// Check the size
	child := w.child.(*vboxElement).children[i]
	if elem, ok := child.(*buttonElement); ok {
		elem.control.PerformClick()
	} else if elem, ok := child.(*checkboxElement); ok {
		elem.control.PerformClick()
	} else {
		panic("Unsupported widget in testingClick")
	}
}

func testingSetFocus(t *testing.T, w *Window, i int) {
	child := w.child.(*vboxElement).children[i]
	if elem, ok := child.(*buttonElement); ok {
		w.handle.MakeFirstResponder(&elem.control.Control)
	} else if elem, ok := child.(*checkboxElement); ok {
		w.handle.MakeFirstResponder(&elem.control.Control)
	} else {
		panic("Unsupported widget in testingClick")
	}
}
