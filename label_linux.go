package goey

import (
	"unsafe"

	"github.com/gotk3/gotk3/gtk"
)

type mountedLabel struct {
	Control
}

func (w *Label) mount(parent Control) (Element, error) {
	handle, err := gtk.LabelNew(w.Text)
	if err != nil {
		return nil, err
	}
	handle.SetSingleLineMode(false)
	(*gtk.Container)(unsafe.Pointer(parent.handle)).Add(handle)
	handle.SetJustify(gtk.JUSTIFY_LEFT)
	handle.SetXAlign(0)
	handle.SetLineWrap(false)
	handle.Show()

	retval := &mountedLabel{Control: Control{&handle.Widget}}
	handle.Connect("destroy", labelOnDestroy, retval)

	return retval, nil
}

func labelOnDestroy(widget *gtk.Label, mounted *mountedLabel) {
	mounted.handle = nil
}

func (w *mountedLabel) label() *gtk.Label {
	return (*gtk.Label)(unsafe.Pointer(w.handle))
}

func (w *mountedLabel) Props() Widget {
	label := w.label()
	text, err := label.GetText()
	if err != nil {
		panic("Could not get text, " + err.Error())
	}

	return &Label{
		Text: text,
	}
}

func (w *mountedLabel) updateProps(data *Label) error {
	label := w.label()
	label.SetText(data.Text)
	return nil
}
