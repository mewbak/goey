package cocoa

/*
#include "cocoa.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// TabView is a wrapper for a NSTabView.
type TabView struct {
	View
	private int
}

type tabviewCallback struct {
	onChange func(int)
}

var (
	tabviewCallbacks = make(map[unsafe.Pointer]tabviewCallback)
)

func NewTabView(window *View) *TabView {
	handle := C.tabviewNew(unsafe.Pointer(window))
	return (*TabView)(handle)
}

func (w *TabView) Close() {
	C.viewClose(unsafe.Pointer(w))
	delete(tabviewCallbacks, unsafe.Pointer(w))
}

func (w *TabView) AddItem(text string) {
	ctext := C.CString(text)
	defer func() {
		C.free(unsafe.Pointer(ctext))
	}()

	C.tabviewAddItem(unsafe.Pointer(w), ctext)
}

func (w *TabView) SelectItem(index int) {
	C.tabviewSelectItem(unsafe.Pointer(w), C.int(index))
}

func (w *TabView) ContentOrigin() (int, int) {
	origin := C.tabviewContentOrigin(unsafe.Pointer(w))
	return int(origin.width), int(origin.height)
}

func (w *TabView) ContentInsets() (int, int) {
	size := C.tabviewContentInsets(unsafe.Pointer(w))
	return int(size.width), int(size.height)
}

func (w *Window) SetOnChange(cb func(int)) {
	tabviewCallbacks[unsafe.Pointer(w)] = tabviewCallback{
		onChange: cb,
	}
}

//export tabviewDidSelectItem
func tabviewDidSelectItem(handle unsafe.Pointer, index int) {
	if cb := tabviewCallbacks[handle]; cb.onChange != nil {
		cb.onChange(index)
	}
}
