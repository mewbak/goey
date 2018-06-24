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

func (w *mountedP) Layout(bc Box) Size {
	_, width := w.handle.GetPreferredWidth()
	_, height := w.handle.GetPreferredHeight()
	return bc.Constrain(Size{FromPixelsX(width), FromPixelsY(height)})
}

func (w *mountedP) MinimumSize() Size {
	width, _ := w.handle.GetPreferredWidth()
	height, _ := w.handle.GetPreferredHeight()
	return Size{FromPixelsX(width), FromPixelsY(height)}
}

func (w *mountedP) Props() Widget {
	text, err := w.handle.GetText()
	if err != nil {
		panic("Could not get text, " + err.Error())
	}

	align := JustifyLeft
	switch w.handle.GetJustify() {
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

func (w *mountedP) SetBounds(bounds Rectangle) {
	pixels := bounds.Pixels()
	syscall.SetBounds(&w.handle.Widget, pixels.Min.X, pixels.Min.Y, pixels.Dx(), pixels.Dy())
}

func (w *mountedP) updateProps(data *P) error {
	w.handle.SetText(data.Text)
	w.handle.SetJustify(data.Align.native())
	return nil
}
