// +build !gnustep

package goey

import (
	"unsafe"

	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/internal/syscall"
	"github.com/gotk3/gotk3/gtk"
)

type paragraphElement struct {
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

func (a TextAlignment) halign() gtk.Align {
	switch a {
	case JustifyLeft:
		return gtk.ALIGN_START
	case JustifyCenter:
		return gtk.ALIGN_CENTER
	case JustifyRight:
		return gtk.ALIGN_END
	case JustifyFull:
		return gtk.ALIGN_START
	}

	panic("not reachable")
}

func (w *P) mount(parent base.Control) (base.Element, error) {
	handle, err := gtk.LabelNew(w.Text)
	if err != nil {
		return nil, err
	}
	handle.SetSingleLineMode(false)
	parent.Handle.Add(handle)
	handle.SetJustify(w.Align.native())
	handle.SetHAlign(w.Align.halign())
	handle.SetLineWrap(true)

	retval := &paragraphElement{Control{&handle.Widget}}
	handle.Connect("destroy", paragraphOnDestroy, retval)
	handle.Show()

	return retval, nil
}

func paragraphOnDestroy(widget *gtk.Label, mounted *paragraphElement) {
	mounted.handle = nil
}

func (w *paragraphElement) label() *gtk.Label {
	return (*gtk.Label)(unsafe.Pointer(w.handle))
}

func (w *paragraphElement) Props() base.Widget {
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

func (w *paragraphElement) measureReflowLimits() {
	label := w.label()

	text, err := label.GetText()
	if err != nil {
		panic(err)
	}

	label.SetText("mmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmm")
	width, _ := label.GetPreferredWidth()
	label.SetText(text)

	paragraphMaxWidth = base.FromPixelsX(width)
}

func (w *paragraphElement) MinIntrinsicHeight(width base.Length) base.Length {
	if width == base.Inf {
		width = w.maxReflowWidth()
	}

	height, _ := syscall.WidgetGetPreferredHeightForWidth(w.handle, width.PixelsX())
	return base.FromPixelsY(height)
}

func (w *paragraphElement) MinIntrinsicWidth(height base.Length) base.Length {
	if height != base.Inf {
		panic("not implemented")
	}

	width, _ := w.label().GetPreferredWidth()
	return min(base.FromPixelsX(int(width)), w.minReflowWidth())
}

func (w *paragraphElement) updateProps(data *P) error {
	label := w.label()
	label.SetText(data.Text)
	label.SetJustify(data.Align.native())
	label.SetHAlign(data.Align.halign())
	return nil
}
