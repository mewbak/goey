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
	px, h := w.control.IntrinsicContentSize()
	return bc.Constrain(base.Size{
		base.FromPixelsX(px),
		base.FromPixelsY(h),
	})
}

func (w *buttonElement) MinIntrinsicHeight(width base.Length) base.Length {
	_, h := w.control.IntrinsicContentSize()
	return base.FromPixelsY(h)
}

func (w *buttonElement) MinIntrinsicWidth(base.Length) base.Length {
	px, _ := w.control.IntrinsicContentSize()
	return base.FromPixelsX(px)
}

func (w *buttonElement) SetBounds(bounds base.Rectangle) {
	px := bounds.Pixels()
	w.control.SetFrame(px.Min.X, px.Min.Y, px.Dx(), px.Dy())
}

func (w *buttonElement) updateProps(data *Button) error {
	w.control.SetTitle(data.Text)
	return nil
}
