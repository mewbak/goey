package goey

import (
	"unsafe"

	"github.com/gotk3/gotk3/gtk"
)

type mountedP struct {
	handle *gtk.Label
}

func (a TextAlignment) native() gtk.Justification {
	switch a {
	case Left:
		return gtk.JUSTIFY_LEFT
	case Center:
		return gtk.JUSTIFY_CENTER
	case Right:
		return gtk.JUSTIFY_RIGHT
	case Justify:
		return gtk.JUSTIFY_FILL
	}

	panic("not reachable")
}

func (w *P) mount(parent NativeWidget) (MountedWidget, error) {
	handle, err := gtk.LabelNew(w.Text)
	if err != nil {
		return nil, err
	}
	handle.SetSingleLineMode(false)
	(*gtk.Container)(unsafe.Pointer(parent.handle)).Add(handle)
	handle.SetJustify(w.Align.native())
	handle.SetLineWrap(true)
	handle.Show()

	retval := &mountedP{handle}
	handle.Connect("destroy", paragraph_onDestroy, retval)

	return retval, nil
}

func paragraph_onDestroy(widget *gtk.Label, mounted *mountedP) {
	mounted.handle = nil
}

func (w *mountedP) Close() {
	if w.handle != nil {
		w.handle.Destroy()
		w.handle = nil
	}
}

func (w *mountedP) Handle() *gtk.Widget {
	return &w.handle.Widget
}

func (w *mountedP) updateProps(data *P) error {
	w.handle.SetText(data.Text)
	w.handle.SetJustify(data.Align.native())
	return nil
}
