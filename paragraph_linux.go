package goey

import (
	"unsafe"

	"bitbucket.org/rj/goey/syscall"
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

func (w *P) mount(parent NativeWidget) (Element, error) {
	handle, err := gtk.LabelNew(w.Text)
	if err != nil {
		return nil, err
	}
	handle.SetSingleLineMode(false)
	(*gtk.Container)(unsafe.Pointer(parent.handle)).Add(handle)
	handle.SetJustify(w.Align.native())
	handle.SetLineWrap(true)

	retval := &mountedP{handle}
	handle.Connect("destroy", paragraph_onDestroy, retval)
	handle.Show()

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

func (w *mountedP) MeasureWidth() (Length, Length) {
	min, max := w.handle.GetPreferredWidth()
	return FromPixelsX(min), FromPixelsX(max)
}

func (w *mountedP) MeasureHeight(width Length) (Length, Length) {
	min, max := syscall.WidgetGetPreferredHeightForWidth(&w.handle.Widget, width.PixelsX())
	return FromPixelsY(min), FromPixelsY(max)
}

func (w *mountedP) SetBounds(bounds Rectangle) {
	pixels := bounds.Pixels()
	syscall.SetBounds(&w.handle.Widget, pixels.Min.X, pixels.Min.Y, pixels.Dx(), pixels.Dy())
}

func (w *mountedP) updateProps(data *P) error {
	w.handle.SetText(data.Text)
	w.handle.SetJustify(data.Align.native())
	return nil
}
