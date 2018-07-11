package goey

import (
	"unsafe"

	"github.com/gotk3/gotk3/gtk"
)

type hrElement struct {
	Control
}

func (w *HR) mount(parent Control) (Element, error) {
	control, err := gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, err
	}
	(*gtk.Container)(unsafe.Pointer(parent.handle)).Add(control)

	retval := &hrElement{
		Control: Control{&control.Widget},
	}

	control.Connect("destroy", hrOnDestroy, retval)
	control.Show()

	return retval, nil
}

func hrOnDestroy(widget *gtk.Separator, mounted *hrElement) {
	mounted.handle = nil
}
