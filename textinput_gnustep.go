// +build gnustep

package goey

import (
	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/cocoa"
)

type textinputElement struct {
	control *cocoa.TextField
}

func (w *TextInput) mount(parent base.Control) (base.Element, error) {
	control := cocoa.NewTextField(parent.Handle, w.Value)
	control.SetValue(w.Value)
	control.SetPlaceholder(w.Placeholder)
	control.SetCallbacks(w.OnChange, w.OnFocus, w.OnBlur)

	retval := &textinputElement{
		control: control,
	}
	return retval, nil
}

func (w *textinputElement) Close() {
	if w.control != nil {
		w.control.Close()
		w.control = nil
	}
}

func (w *textinputElement) Layout(bc base.Constraints) base.Size {
	px, h := w.control.IntrinsicContentSize()
	return bc.Constrain(base.Size{
		base.FromPixelsX(px),
		base.FromPixelsY(h),
	})
}

func (w *textinputElement) MinIntrinsicHeight(width base.Length) base.Length {
	_, h := w.control.IntrinsicContentSize()
	return base.FromPixelsY(h)
}

func (w *textinputElement) MinIntrinsicWidth(base.Length) base.Length {
	px, _ := w.control.IntrinsicContentSize()
	return base.FromPixelsX(px)
}

func (w *textinputElement) SetBounds(bounds base.Rectangle) {
	px := bounds.Pixels()
	w.control.SetFrame(px.Min.X, px.Min.Y, px.Dx(), px.Dy())
}

func (w *textinputElement) updateProps(data *TextInput) error {
	w.control.SetValue(data.Value)
	w.control.SetPlaceholder(data.Placeholder)
	w.control.SetCallbacks(data.OnChange, data.OnFocus, data.OnBlur)
	return nil
}
