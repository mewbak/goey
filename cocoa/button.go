package cocoa

/*
#include "cocoa.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

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
	C.buttonClose(unsafe.Pointer(w))
	delete(buttonCallbacks, unsafe.Pointer(w))
}

func (w *Button) SetCallbacks(onclick func(), onfocus func(), onblur func()) {
	buttonCallbacks[unsafe.Pointer(w)] = buttonCallback{
		onClick: onclick,
		onFocus: onfocus,
		onBlur:  onblur,
	}
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
	println("focus!", handle)
}

//export buttonOnBlur
func buttonOnBlur(handle unsafe.Pointer) {
	println("blur!", handle)
}
