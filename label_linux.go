package goey

import (
	"unsafe"

	"github.com/gotk3/gotk3/gtk"
)

type mountedLabel struct {
	NativeWidget
}

func (w *Label) Mount(parent NativeWidget) (MountedWidget, error) {
	control, err := gtk.LabelNew(w.Text)
	if err != nil {
		return nil, err
	}
	control.SetSingleLineMode(false)
	(*gtk.Container)(unsafe.Pointer(parent.handle)).Add(control)
	control.SetJustify(gtk.JUSTIFY_LEFT)
	control.SetLineWrap(false)
	control.Show()

	retval := &mountedLabel{
		NativeWidget: NativeWidget{&control.Widget},
	}

	control.Connect("destroy", label_onDestroy, retval)

	return retval, nil
}

func label_onDestroy(widget *gtk.Label, mounted *mountedLabel) {
	mounted.handle = nil
}

func (w *mountedLabel) UpdateProps(data Widget) error {
	panic("not implemented")
}
