package goey

import (
	"unsafe"

	"github.com/gotk3/gotk3/gtk"
)

type mountedP struct {
	Control
}

func (a TextAlignment) native() gtk.Justification {
	switch a {
	case JustifyLeft:
		return gtk.JUSTIFY_LEFT
	case JustifyCenter:
		return gtk.JUSTIFY_CENTER
	case JustifyRight:
		return gtk.JUSTIFY_RIGHT
	case JustifyFull:
		return gtk.JUSTIFY_FILL
	}

	panic("not reachable")
}

func (w *P) mount(parent Control) (Element, error) {
	handle, err := gtk.LabelNew(w.Text)
	if err != nil {
		return nil, err
	}
	handle.SetSingleLineMode(false)
	(*gtk.Container)(unsafe.Pointer(parent.handle)).Add(handle)
	handle.SetJustify(w.Align.native())
	handle.SetLineWrap(true)

	retval := &mountedP{Control{&handle.Widget}}
	handle.Connect("destroy", paragraph_onDestroy, retval)
	handle.Show()

	return retval, nil
}

func paragraph_onDestroy(widget *gtk.Label, mounted *mountedP) {
	mounted.handle = nil
}

func (w *mountedP) label() *gtk.Label {
	return (*gtk.Label)(unsafe.Pointer(w.handle))
}

func (w *mountedP) Props() Widget {
	label := w.label()
	text, err := label.GetText()
	if err != nil {
		panic("Could not get text, " + err.Error())
	}

	align := JustifyLeft
	switch label.GetJustify() {
	case gtk.JUSTIFY_CENTER:
		align = JustifyCenter
	case gtk.JUSTIFY_RIGHT:
		align = JustifyRight
	case gtk.JUSTIFY_FILL:
		align = JustifyFull
	}

	return &P{
		Text:  text,
		Align: align,
	}
}

func (w *mountedP) updateProps(data *P) error {
	label := w.label()
	label.SetText(data.Text)
	label.SetJustify(data.Align.native())
	return nil
}
