package goey

import (
	"unsafe"

	"github.com/gotk3/gotk3/gtk"
)

type mountedHR struct {
	NativeWidget
}

func (w *HR) mount(parent NativeWidget) (Element, error) {
	control, err := gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, err
	}
	(*gtk.Container)(unsafe.Pointer(parent.handle)).Add(control)

	retval := &mountedHR{
		NativeWidget: NativeWidget{&control.Widget},
	}

	control.Connect("destroy", hr_onDestroy, retval)
	control.Show()

	return retval, nil
}

func hr_onDestroy(widget *gtk.Separator, mounted *mountedHR) {
	mounted.handle = nil
}
