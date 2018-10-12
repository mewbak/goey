package cocoa

/*
#include "cocoa.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// Window is a wrapper for a NSWindow.
type Window struct {
	private int
}

type WindowCallbacks interface {
	OnClosing() bool
	OnClose()
}

var (
	windowCallbacks = make(map[unsafe.Pointer]WindowCallbacks)
)

func NewWindow(title string, width, height uint) *Window {
	ctitle := C.CString(title)
	defer func() {
		C.free(unsafe.Pointer(ctitle))
	}()

	handle := C.windowNew(ctitle, C.unsigned(width), C.unsigned(height))
	return (*Window)(handle)
}

func (w *Window) Close() {
	C.windowClose(unsafe.Pointer(w))
}

func (w *Window) ContentSize() (int, int) {
	var h C.int
	px := C.windowContentSize(unsafe.Pointer(w), &h)
	return int(px), int(h)
}

func (w *Window) SetCallbacks(cb WindowCallbacks) {
	windowCallbacks[unsafe.Pointer(w)] = cb
}

//export windowShouldClose
func windowShouldClose(handle unsafe.Pointer) bool {
	if cb := windowCallbacks[handle]; cb != nil {
		return !cb.OnClosing()
	}
	
	return true
}

//export windowWillClose
func windowWillClose(handle unsafe.Pointer) {
	if cb := windowCallbacks[handle]; cb != nil {
		cb.OnClose()
	}
	delete( windowCallbacks, handle )
}
