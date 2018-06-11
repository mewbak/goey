package goey

import (
	"unsafe"

	"github.com/gotk3/gotk3/gtk"
)

type mountedLabel struct {
	handle *gtk.Label
}

func (w *Label) mount(parent NativeWidget) (Element, error) {
	handle, err := gtk.LabelNew(w.Text)
	if err != nil {
		return nil, err
	}
	handle.SetSingleLineMode(false)
	(*gtk.Container)(unsafe.Pointer(parent.handle)).Add(handle)
	handle.SetJustify(gtk.JUSTIFY_LEFT)
	handle.SetLineWrap(false)
	handle.Show()

	retval := &mountedLabel{handle}
	handle.Connect("destroy", label_onDestroy, retval)

	return retval, nil
}

func label_onDestroy(widget *gtk.Label, mounted *mountedLabel) {
	mounted.handle = nil
}

func (w *mountedLabel) Close() {
	if w.handle != nil {
		w.handle.Destroy()
		w.handle = nil
	}
}

func (w *mountedLabel) Handle() *gtk.Widget {
	return &w.handle.Widget
}

func (w *mountedLabel) updateProps(data *Label) error {
	w.handle.SetText(data.Text)
	return nil
}
