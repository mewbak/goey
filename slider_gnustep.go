// +build gnustep

package goey

import (
	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/cocoa"
)

type sliderElement struct {
	control *cocoa.Text
}

func (w *Slider) mount(parent base.Control) (base.Element, error) {
	control := cocoa.NewText(parent.Handle, "date input")

	retval := &sliderElement{
		control: control,
	}
	return retval, nil
}

func (w *sliderElement) Close() {
	if w.control != nil {
		w.control.Close()
		w.control = nil
	}
}

func (w *sliderElement) Layout(bc base.Constraints) base.Size {
	px := w.MinIntrinsicWidth(base.Inf)
	h := w.MinIntrinsicHeight(base.Inf)
	return bc.Constrain(base.Size{px, h})
}

func (w *sliderElement) MinIntrinsicHeight(width base.Length) base.Length {
	return 20 * base.DIP
}

func (w *sliderElement) MinIntrinsicWidth(base.Length) base.Length {
	return 200 * base.DIP
}

func (w *sliderElement) SetBounds(bounds base.Rectangle) {
	px := bounds.Pixels()
	w.control.SetFrame(px.Min.X, px.Min.Y, px.Dx(), px.Dy())
}

func (w *sliderElement) updateProps(data *Slider) error {
	return nil
}
