// +build gnustep

package goey

import (
	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/cocoa"
)

type paragraphElement struct {
	control *cocoa.Text
}

func (w *P) mount(parent base.Control) (base.Element, error) {
	control := cocoa.NewText(parent.Handle, w.Text)
	control.SetAlignment(int(w.Align))

	retval := &labelElement{
		control: control,
	}
	return retval, nil
}

func (w *paragraphElement) Close() {
	if w.control != nil {
		w.control.Close()
		w.control = nil
	}
}

func (w *paragraphElement) measureReflowLimits() {
	paragraphMaxWidth = 40 * base.DIP
}

func (w *paragraphElement) MinIntrinsicHeight(width base.Length) base.Length {
	return 20 * base.DIP
}

func (w *paragraphElement) MinIntrinsicWidth(base.Length) base.Length {
	return 200 * base.DIP
}

func (w *paragraphElement) SetBounds(bounds base.Rectangle) {
	px := bounds.Pixels()
	w.control.SetFrame(px.Min.X, px.Min.Y, px.Dx(), px.Dy())
}

func (w *paragraphElement) updateProps(data *P) error {
	w.control.SetText(data.Text)
	w.control.SetAlignment(int(w.Align))
	return nil
}
