package cocoa

/*
#include "cocoa.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// Button is a wrapper for a NSButton.
type Button struct {
	Control
	private int
}

type buttonCallback struct {
	onClick func()
	onFocus func()
	onBlur  func()
}

var (
	buttonCallbacks = make(map[unsafe.Pointer]buttonCallback)
)

func NewButton(window *Window, title string) *Button {
	ctitle := C.CString(title)
	defer func() {
		C.free(unsafe.Pointer(ctitle))
	}()

	handle := C.buttonNew(unsafe.Pointer(window), ctitle)
	return (*Button)(handle)
}

func (w *Button) Close() {
	C.controlClose(unsafe.Pointer(w))
	delete(buttonCallbacks, unsafe.Pointer(w))
}

func (w *Button) PerformClick() {
	C.buttonPerformClick(unsafe.Pointer(w))
}

func (w *Button) Callbacks() (func(), func(), func()) {
	cb := buttonCallbacks[unsafe.Pointer(w)]
	return cb.onClick, cb.onFocus, cb.onBlur
}

func (w *Button) SetCallbacks(onclick func(), onfocus func(), onblur func()) {
	buttonCallbacks[unsafe.Pointer(w)] = buttonCallback{
		onClick: onclick,
		onFocus: onfocus,
		onBlur:  onblur,
	}
}

func (w *Button) Title() string {
	//cstring := C.buttonTitle(unsafe.Pointer(w))
	return "" //C.GoString(cstring)
}

func (w *Button) SetTitle(title string) {
	ctitle := C.CString(title)
	defer func() {
		C.free(unsafe.Pointer(ctitle))
	}()

	C.buttonSetTitle(unsafe.Pointer(w), ctitle)
}

//export buttonOnClick
func buttonOnClick(handle unsafe.Pointer) {
	if cb := buttonCallbacks[handle]; cb.onClick != nil {
		cb.onClick()
	}
}

//export buttonOnFocus
func buttonOnFocus(handle unsafe.Pointer) {
	if cb := buttonCallbacks[handle]; cb.onFocus != nil {
		cb.onFocus()
	}
}

//export buttonOnBlur
func buttonOnBlur(handle unsafe.Pointer) {
	if cb := buttonCallbacks[handle]; cb.onBlur != nil {
		cb.onBlur()
	}
}
