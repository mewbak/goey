package goey

import (
	"unsafe"

	"github.com/gotk3/gotk3/gtk"
)

type MountedP struct {
	NativeWidget
}

func (w *P) Mount(parent NativeWidget) (MountedWidget, error) {
	control, err := gtk.LabelNew(w.Text)
	if err != nil {
		return nil, err
	}
	control.SetSingleLineMode(false)
	(*gtk.Container)(unsafe.Pointer(parent.handle)).Add(control)
	control.SetLineWrap(true)
	control.Show()

	retval := &MountedP{
		NativeWidget: NativeWidget{&control.Widget},
	}

	control.Connect("destroy", paragraph_onDestroy, retval)

	return retval, nil
}

func paragraph_onDestroy(widget *gtk.Label, mounted *MountedP) {
	mounted.handle = nil
}

func (w *MountedP) UpdateProps(data Widget) error {
	panic("not implemented")
}
