// +build gnustep

package goey

import (
	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/cocoa"
)

type decorationElement struct {
	control *cocoa.Decoration
	insets  Insets
	radius  base.Length

	child     base.Element
	childSize base.Size
}

func (w *Decoration) mount(parent base.Control) (base.Element, error) {
	control := cocoa.NewDecoration(parent.Handle, w.Fill, w.Stroke)

	retval := &decorationElement{
		control: control,
		insets:  w.Insets,
		radius:  w.Radius,
	}

	child, err := base.DiffChild(base.Control{&control.View}, nil, w.Child)
	if err != nil {
		control.Close()
		return nil, err
	}
	retval.child = child

	return retval, nil
}

func (w *decorationElement) Close() {
	if w.child != nil {
		w.child.Close()
		w.child = nil
	}
	if w.control != nil {
		w.control.Close()
		w.control = nil
	}
}

func (w *decorationElement) props() *Decoration {
	return &Decoration{}
}

func (w *decorationElement) SetBounds(bounds base.Rectangle) {
	px := bounds.Pixels()
	w.control.SetFrame(px.Min.X, px.Min.Y, px.Dx(), px.Dy())

	if w.child != nil {
		width := bounds.Dx()
		height := bounds.Dy()

		bounds.Min.X = w.insets.Left
		bounds.Min.Y = w.insets.Top
		bounds.Max.X = width - w.insets.Right
		bounds.Max.Y = height - w.insets.Bottom
		w.child.SetBounds(bounds)
	}
}

func (w *decorationElement) updateProps(data *Decoration) error {
	child, err := base.DiffChild(base.Control{&w.control.View}, w.child, data.Child)
	if err != nil {
		return err
	}
	w.child = child

	w.control.SetBorderRadius(data.Radius.PixelsX(), data.Radius.PixelsY())
	w.control.SetFillColor(data.Fill)
	w.control.SetStrokeColor(data.Stroke)
	return nil
}
