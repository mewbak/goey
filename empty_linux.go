package goey

import (
	"unsafe"

	"github.com/gotk3/gotk3/gtk"
)

type mountedEmpty struct {
	Control
}

func (w *Empty) mount(parent Control) (Element, error) {
	control, err := gtk.LabelNew("  ")
	if err != nil {
		return nil, err
	}
	(*gtk.Container)(unsafe.Pointer(parent.handle)).Add(control)

	retval := &mountedEmpty{
		Control: Control{&control.Widget},
	}

	control.Connect("destroy", empty_onDestroy, retval)
	control.Show()

	return retval, nil
}

func empty_onDestroy(widget *gtk.Separator, mounted *mountedHR) {
	mounted.handle = nil
}
