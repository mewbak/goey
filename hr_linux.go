package goey

import (
	"unsafe"

	"github.com/gotk3/gotk3/gtk"
)

type MountedHR struct {
	NativeWidget
}

func (w *HR) Mount(parent NativeWidget) (MountedWidget, error) {
	control, err := gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL)
	if err != nil {
		return nil, err
	}
	(*gtk.Container)(unsafe.Pointer(parent.handle)).Add(control)

	retval := &MountedHR{
		NativeWidget: NativeWidget{&control.Widget},
	}

	control.Connect("destroy", hr_onDestroy, retval)
	control.Show()

	return retval, nil
}

func hr_onDestroy(widget *gtk.Separator, mounted *MountedHR) {
	mounted.handle = nil
}

func (w *MountedHR) UpdateProps(data Widget) error {
	return nil
}
