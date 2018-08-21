package icons

import (
	"bitbucket.org/rj/goey/base"
)

type iconElement struct {
	child base.Element
	icon  rune
}

func (w *iconElement) Close() {
	w.child.Close()
	w.child = nil
}

func (*iconElement) Kind() *base.Kind {
	return &kind
}

func (w *iconElement) Layout(bc base.Constraints) base.Size {
	return w.child.Layout(bc)
}

func (w *iconElement) MinIntrinsicHeight(width base.Length) base.Length {
	return w.child.MinIntrinsicHeight(width)
}

func (w *iconElement) MinIntrinsicWidth(height base.Length) base.Length {
	return w.child.MinIntrinsicWidth(height)
}

func (w *iconElement) SetBounds(bounds base.Rectangle) {
	w.child.SetBounds(bounds)
}

func (w *iconElement) UpdateProps(data base.Widget) error {
	// TODO
	return nil
}
