// +build gnustep

package goey

import (
	"bitbucket.org/rj/goey/base"
)

type hrElement struct {
}

func (w *HR) mount(parent base.Control) (base.Element, error) {
	//control := cocoa.NewText(parent.Handle, w.Text)

	retval := &hrElement{
		//control: control,
	}
	return retval, nil
}

func (w *hrElement) Close() {
	/*if w.control != nil {
		w.control.Close()
		w.control = nil
	}*/
}

func (w *hrElement) SetBounds(bounds base.Rectangle) {
	//px := bounds.Pixels()
	//w.control.SetFrame(px.Min.X, px.Min.Y, px.Dx(), px.Dy())
}
