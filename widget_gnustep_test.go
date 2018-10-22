// +build gnustep

package goey

import (
	"testing"
)

func testingSetFocus(t *testing.T, w *Window, i int) {
	child := w.child.(*vboxElement).children[i]
	if elem, ok := child.(*buttonElement); ok {
		w.handle.MakeFirstResponder(&elem.control.Control)
	} else if elem, ok := child.(*checkboxElement); ok {
		w.handle.MakeFirstResponder(&elem.control.Control)
	} else if elem, ok := child.(*textareaElement); ok {
		w.handle.MakeFirstResponder(&elem.control.Control)
	} else if elem, ok := child.(*textinputElement); ok {
		w.handle.MakeFirstResponder(&elem.control.Control)
	} else {
		panic("Unsupported widget in testingSetFocus")
	}
}
