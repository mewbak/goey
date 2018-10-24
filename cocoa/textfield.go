package cocoa

/*
#include "cocoa.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// TextField is a wrapper for a NSTextField.
type TextField struct {
	Control
	private int
}

type textfieldCallback struct {
	onChange func(string)
	onFocus  func()
	onBlur   func()
}

var (
	textfieldCallbacks = make(map[unsafe.Pointer]textfieldCallback)
)

func NewTextField(window *View, title string) *TextField {
	ctitle := C.CString(title)
	defer func() {
		C.free(unsafe.Pointer(ctitle))
	}()

	handle := C.textfieldNew(unsafe.Pointer(window), ctitle)
	return (*TextField)(handle)
}

func (w *TextField) Close() {
	C.viewClose(unsafe.Pointer(w))
	delete(textfieldCallbacks, unsafe.Pointer(w))
}

func (w *TextField) SetCallbacks(onchange func(string), onfocus func(), onblur func()) {
	textfieldCallbacks[unsafe.Pointer(w)] = textfieldCallback{
		onChange: onchange,
		onFocus:  onfocus,
		onBlur:   onblur,
	}
}

func (w *TextField) SetValue(text string) {
	ctitle := C.CString(text)
	defer func() {
		C.free(unsafe.Pointer(ctitle))
	}()

	C.textfieldSetValue(unsafe.Pointer(w), ctitle)
}

func (w *TextField) SetPlaceholder(text string) {
	ctitle := C.CString(text)
	defer func() {
		C.free(unsafe.Pointer(ctitle))
	}()

	C.textfieldSetPlaceholder(unsafe.Pointer(w), ctitle)
}

//export textfieldOnChange
func textfieldOnChange(handle unsafe.Pointer, text *C.char) {
	println("textfieldOnChange", handle, C.GoString(text))
	if cb := textfieldCallbacks[handle]; cb.onChange != nil {
		println("found")
		cb.onChange(C.GoString(text))
	}
}

//export textfieldOnFocus
func textfieldOnFocus(handle unsafe.Pointer) {
	if cb := textfieldCallbacks[handle]; cb.onFocus != nil {
		cb.onFocus()
	}
}

//export textfieldOnBlur
func textfieldOnBlur(handle unsafe.Pointer) {
	if cb := textfieldCallbacks[handle]; cb.onBlur != nil {
		cb.onBlur()
	}
}
