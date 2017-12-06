package goey

import (
	"unsafe"

	"github.com/gotk3/gotk3/gtk"
)

type mountedP struct {
	NativeWidget
}

func (w *P) mount(parent NativeWidget) (MountedWidget, error) {
	control, err := gtk.LabelNew(w.Text)
	if err != nil {
		return nil, err
	}
	control.SetSingleLineMode(false)
	(*gtk.Container)(unsafe.Pointer(parent.handle)).Add(control)
	control.SetLineWrap(true)
	control.Show()

	retval := &mountedP{
		NativeWidget: NativeWidget{&control.Widget},
	}

	control.Connect("destroy", paragraph_onDestroy, retval)

	return retval, nil
}

func paragraph_onDestroy(widget *gtk.Label, mounted *mountedP) {
	mounted.handle = nil
}

func (w *mountedP) updateProps(data *P) error {
	panic("not implemented")
}
