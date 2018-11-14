// +build gnustep

package goey

import (
	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/cocoa"
)

type intinputElement struct {
	control *cocoa.Text
}

func (w *IntInput) mount(parent base.Control) (base.Element, error) {
	control := cocoa.NewText(parent.Handle, "date input")

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
	return 20 * base.DIP
}

func (w *intinputElement) MinIntrinsicWidth(base.Length) base.Length {
	return 200 * base.DIP
}

func (w *intinputElement) SetBounds(bounds base.Rectangle) {
	px := bounds.Pixels()
	w.control.SetFrame(px.Min.X, px.Min.Y, px.Dx(), px.Dy())
}

func (w *intinputElement) updateProps(data *IntInput) error {
	return nil
}
