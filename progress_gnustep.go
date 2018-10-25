// +build gnustep

package goey

import (
	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/cocoa"
)

type progressElement struct {
	control *cocoa.Text
}

func (w *Progress) mount(parent base.Control) (base.Element, error) {
	control := cocoa.NewText(parent.Handle, "date input")

	retval := &progressElement{
		control: control,
	}
	return retval, nil
}

func (w *progressElement) Close() {
	if w.control != nil {
		w.control.Close()
		w.control = nil
	}
}

func (w *progressElement) Layout(bc base.Constraints) base.Size {
	px := w.MinIntrinsicWidth(base.Inf)
	h := w.MinIntrinsicHeight(base.Inf)
	return bc.Constrain(base.Size{px, h})
}

func (w *progressElement) MinIntrinsicHeight(width base.Length) base.Length {
	return 20 * base.DIP
}

func (w *progressElement) MinIntrinsicWidth(base.Length) base.Length {
	return 200 * base.DIP
}

func (w *progressElement) SetBounds(bounds base.Rectangle) {
	px := bounds.Pixels()
	w.control.SetFrame(px.Min.X, px.Min.Y, px.Dx(), px.Dy())
}

func (w *progressElement) updateProps(data *Progress) error {
	return nil
}
