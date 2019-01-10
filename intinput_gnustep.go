// +build gnustep

package goey

import (
	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/cocoa"
)

type intinputElement struct {
	control *cocoa.IntField
}

func (w *IntInput) mount(parent base.Control) (base.Element, error) {
	control := cocoa.NewIntField(parent.Handle, w.Value, -100, 100)
	control.SetPlaceholder(w.Placeholder)
	control.SetEnabled(!w.Disabled)
	control.SetCallbacks(w.OnChange, w.OnFocus, w.OnBlur)

	retval := &intinputElement{
		control: control,
	}
	return retval, nil
}

func (w *intinputElement) Close() {
	if w.control != nil {
		w.control.Close()
		w.control = nil
	}
}

func (w *intinputElement) Layout(bc base.Constraints) base.Size {
	px := w.MinIntrinsicWidth(base.Inf)
	h := w.MinIntrinsicHeight(base.Inf)
	return bc.Constrain(base.Size{px, h})
}

func (w *intinputElement) MinIntrinsicHeight(width base.Length) base.Length {
	_, h := w.control.IntrinsicContentSize()
	return base.FromPixelsY(h)
}

func (w *intinputElement) MinIntrinsicWidth(base.Length) base.Length {
	px, _ := w.control.IntrinsicContentSize()
	return base.FromPixelsX(px)
}

func (w *intinputElement) SetBounds(bounds base.Rectangle) {
	px := bounds.Pixels()
	w.control.SetFrame(px.Min.X, px.Min.Y, px.Dx(), px.Dy())
}

func (w *intinputElement) updateProps(data *IntInput) error {
	w.control.SetValue(data.Value, -100, 100)
	w.control.SetPlaceholder(data.Placeholder)
	w.control.SetEnabled(!data.Disabled)
	w.control.SetCallbacks(data.OnChange, data.OnFocus, data.OnBlur)
	return nil
}
