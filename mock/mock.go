package mock

import (
	"bitbucket.org/rj/goey/base"
)

var (
	mockKind = base.NewKind("bitbucket.org/rj/goey/goeytest.Mock")
)

func New(size base.Size) *Element {
	return &Element{
		Size: size,
	}
}

func NewList(sizes ...base.Size) []base.Element {
	ret := make([]base.Element, 0, len(sizes))
	for _, v := range sizes {
		ret = append(ret, &Element{Size: v})
	}
	return ret
}

type Widget struct {
	Size base.Size
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*Widget) Kind() *base.Kind {
	return &mockKind
}

// Mount creates an mock control.
func (w *Widget) Mount(parent base.Control) (base.Element, error) {
	return &Element{
		Size: w.Size,
	}, nil
}

type Element struct {
	Size base.Size

	bounds base.Rectangle
	closed bool
}

func (w *Element) Close() {
	if w.closed {
		panic("Element already closed")
	}
	w.closed = true
}

func (*Element) Kind() *base.Kind {
	return &mockKind
}

func (w *Element) Layout(bc base.Constraints) base.Size {
	return bc.Constrain(w.Size)
}

func (w *Element) MinIntrinsicHeight(base.Length) base.Length {
	return w.Size.Height
}

func (w *Element) MinIntrinsicWidth(base.Length) base.Length {
	return w.Size.Width
}

func (w *Element) Props() base.Widget {
	return &Widget{
		Size: w.Size,
	}
}

func (w *Element) Bounds() base.Rectangle {
	return w.bounds
}

func (w *Element) SetBounds(bounds base.Rectangle) {
	w.bounds = bounds
}

func (w *Element) updateProps(data *Widget) error {
	w.Size = data.Size
	return nil
}

func (w *Element) UpdateProps(data base.Widget) error {
	return w.updateProps(data.(*Widget))
}
