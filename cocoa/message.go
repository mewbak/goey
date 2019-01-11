package cocoa

/*
#include "cocoa.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

func MessageDialog(handle *Window, text string, title string, icon byte) {
	ctext := C.CString(text)
	defer func() {
		C.free(unsafe.Pointer(ctext))
	}()
	ctitle := C.CString(title)
	defer func() {
		C.free(unsafe.Pointer(ctitle))
	}()

	C.messageDialog(unsafe.Pointer(handle), ctext, ctitle, C.char(icon))
}

func OpenPanel(handle *Window) {
	C.openPanel(unsafe.Pointer(handle))
}

func SavePanel(handle *Window) {
	C.savePanel(unsafe.Pointer(handle))
}
