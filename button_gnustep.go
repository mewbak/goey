// +build gnustep

package goey

import (
	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/cocoa"
)

type buttonElement struct {
	control *cocoa.Button
}

func (w *Button) mount(parent base.Control) (base.Element, error) {
	control := cocoa.NewButton(parent.Handle, w.Text)
	control.SetCallbacks(w.OnClick, w.OnFocus, w.OnBlur)

	retval := &buttonElement{
		control: control,
	}
	return retval, nil
}

func (w *buttonElement) Close() {
	if w.control != nil {
		w.control.Close()
		w.control = nil
	}
}

func (w *buttonElement) Layout(bc base.Constraints) base.Size {
	return bc.Constrain(base.Size{100 * base.DIP, 17 * base.DIP})
}

func (w *buttonElement) MinIntrinsicHeight(width base.Length) base.Length {
	return 17 * base.DIP
}

func (w *buttonElement) MinIntrinsicWidth(base.Length) base.Length {
	return 100 * base.DIP
}

func (w *buttonElement) SetBounds(bounds base.Rectangle) {
	px := bounds.Pixels()
	w.control.SetBounds(px.Min.X, px.Min.Y, px.Dx(), px.Dy())
}

func (w *buttonElement) updateProps(data *Button) error {
	w.control.SetTitle(data.Text)
	return nil
}
