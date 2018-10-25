// +build gnustep

package goey

import (
	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/cocoa"
)

type selectinputElement struct {
	control *cocoa.Text
}

func (w *SelectInput) mount(parent base.Control) (base.Element, error) {
	control := cocoa.NewText(parent.Handle, "date input")

	retval := &selectinputElement{
		control: control,
	}
	return retval, nil
}

func (w *selectinputElement) Close() {
	if w.control != nil {
		w.control.Close()
		w.control = nil
	}
}

func (w *selectinputElement) Layout(bc base.Constraints) base.Size {
	px := w.MinIntrinsicWidth(base.Inf)
	h := w.MinIntrinsicHeight(base.Inf)
	return bc.Constrain(base.Size{px, h})
}

func (w *selectinputElement) MinIntrinsicHeight(width base.Length) base.Length {
	return 20 * base.DIP
}

func (w *selectinputElement) MinIntrinsicWidth(base.Length) base.Length {
	return 200 * base.DIP
}

func (w *selectinputElement) SetBounds(bounds base.Rectangle) {
	px := bounds.Pixels()
	w.control.SetFrame(px.Min.X, px.Min.Y, px.Dx(), px.Dy())
}

func (w *selectinputElement) updateProps(data *SelectInput) error {
	return nil
}
