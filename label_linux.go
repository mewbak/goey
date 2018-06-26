package goey

import (
	"unsafe"

	"bitbucket.org/rj/goey/syscall"
	"github.com/gotk3/gotk3/gtk"
)

type mountedLabel struct {
	handle *gtk.Label
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

func (w *mountedLabel) Layout(bc Constraint) Size {
	_, width := w.handle.GetPreferredWidth()
	_, height := w.handle.GetPreferredHeight()
	return bc.Constrain(Size{FromPixelsX(width), FromPixelsY(height)})
}

func (w *mountedLabel) MinimumSize() Size {
	width, _ := w.handle.GetPreferredWidth()
	height, _ := w.handle.GetPreferredHeight()
	return Size{FromPixelsX(width), FromPixelsY(height)}
}

func (w *mountedLabel) Props() Widget {
	text, err := w.handle.GetText()
	if err != nil {
		panic("Could not get text, " + err.Error())
	}

	return &Label{
		Text: text,
	}
}

func (w *mountedLabel) SetBounds(bounds Rectangle) {
	pixels := bounds.Pixels()
	syscall.SetBounds(&w.handle.Widget, pixels.Min.X, pixels.Min.Y, pixels.Dx(), pixels.Dy())
}

func (w *mountedLabel) updateProps(data *Label) error {
	w.handle.SetText(data.Text)
	return nil
}
