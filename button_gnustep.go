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
	return base.Size{}
}

func (w *buttonElement) MinIntrinsicHeight(width base.Length) base.Length {
	return 0
}

func (w *buttonElement) MinIntrinsicWidth(base.Length) base.Length {
	return 0
}

func (w *buttonElement) SetBounds(bounds base.Rectangle) {

}

func (w *buttonElement) updateProps(data *Button) error {
	return nil
}
