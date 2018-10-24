// +build gnustep

package goey

import (
	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/cocoa"
)

type textareaElement struct {
	control *cocoa.TextField
}

func (w *TextArea) mount(parent base.Control) (base.Element, error) {
	control := cocoa.NewTextField(parent.Handle, w.Value)
	control.SetValue(w.Value)
	control.SetPlaceholder(w.Placeholder)
	control.SetCallbacks(w.OnChange, w.OnFocus, w.OnBlur)

	retval := &textareaElement{
		control: control,
	}
	return retval, nil
}

func (w *textareaElement) Close() {
	if w.control != nil {
		w.control.Close()
		w.control = nil
	}
}

func (w *textareaElement) Layout(bc base.Constraints) base.Size {
	px, h := w.control.IntrinsicContentSize()
	return bc.Constrain(base.Size{
		base.FromPixelsX(px),
		base.FromPixelsY(h),
	})
}

func (w *textareaElement) MinIntrinsicHeight(width base.Length) base.Length {
	_, h := w.control.IntrinsicContentSize()
	return base.FromPixelsY(h)
}

func (w *textareaElement) MinIntrinsicWidth(base.Length) base.Length {
	px, _ := w.control.IntrinsicContentSize()
	return base.FromPixelsX(px)
}

func (w *textareaElement) SetBounds(bounds base.Rectangle) {
	px := bounds.Pixels()
	w.control.SetFrame(px.Min.X, px.Min.Y, px.Dx(), px.Dy())
}

func (w *textareaElement) updateProps(data *TextArea) error {
	return nil
}
