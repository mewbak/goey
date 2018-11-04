// +build gnustep

package goey

import (
	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/cocoa"
)

type selectinputElement struct {
	control *cocoa.PopUpButton
}

func (w *SelectInput) mount(parent base.Control) (base.Element, error) {
	control := cocoa.NewPopUpButton(parent.Handle)
	for _, v := range w.Items {
		control.AddItem(v)
	}
	control.SetValue(w.Value)
	control.SetEnabled(!w.Disabled)
	control.SetCallbacks(w.OnChange, w.OnFocus, w.OnBlur)

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
	px, h := w.control.IntrinsicContentSize()
	return bc.Constrain(base.Size{
		base.FromPixelsX(px),
		base.FromPixelsY(h),
	})
}

func (w *selectinputElement) MinIntrinsicHeight(width base.Length) base.Length {
	_, h := w.control.IntrinsicContentSize()
	return base.FromPixelsY(h)
}

func (w *selectinputElement) MinIntrinsicWidth(base.Length) base.Length {
	px, _ := w.control.IntrinsicContentSize()
	return base.FromPixelsX(px)
}

func (w *selectinputElement) TakeFocus() bool {
	return w.control.MakeFirstResponder()
}

func (w *selectinputElement) SetBounds(bounds base.Rectangle) {
	px := bounds.Pixels()
	w.control.SetFrame(px.Min.X, px.Min.Y, px.Dx(), px.Dy())
}

func (w *selectinputElement) updateProps(data *SelectInput) error {
	return nil
}
