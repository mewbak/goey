package goey

import (
	"image/color"

	"bitbucket.org/rj/goey/base"
)

var (
	decorationKind = base.NewKind("bitbucket.org/rj/goey.Decoration")
)

// Decoration describes a widget that provides a border and background, and
// possibly containing a single child widget.
//
// The size of the control will match the size of the child element, although
// padding will be added between the border of the decoration and the child
// element as specified by the field Insets.
type Decoration struct {
	Fill   color.RGBA  // Background colour used to fill interior.
	Stroke color.RGBA  // Stroke colour used to draw outline.
	Insets Insets      // Space between border of the decoration and the child element.
	Radius base.Length // Radius of the widgets corners.
	Child  base.Widget // Child widget.
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*Decoration) Kind() *base.Kind {
	return &decorationKind
}

// Mount creates a button in the GUI.  The newly created widget
// will be a child of the widget specified by parent.
func (w *Decoration) Mount(parent base.Control) (base.Element, error) {
	// Forward to the platform-dependant code
	return w.mount(parent)
}

func (*decorationElement) Kind() *base.Kind {
	return &decorationKind
}

func (w *decorationElement) Layout(bc base.Constraints) base.Size {
	hinset := w.insets.Left + w.insets.Right
	vinset := w.insets.Top + w.insets.Bottom

	innerConstraints := bc.Inset(hinset, vinset)
	w.childSize = base.Layout(w.child, innerConstraints)
	return base.Size{
		w.childSize.Width + hinset,
		w.childSize.Height + vinset,
	}
}

func (w *decorationElement) MinIntrinsicHeight(width base.Length) base.Length {
	vinset := w.insets.Top + w.insets.Bottom

	if w.child == nil {
		return vinset
	}

	return w.child.MinIntrinsicHeight(width) + vinset
}

func (w *decorationElement) MinIntrinsicWidth(height base.Length) base.Length {
	hinset := w.insets.Left + w.insets.Right

	if w.child == nil {
		return hinset
	}

	return w.child.MinIntrinsicWidth(height) + hinset
}

func (w *decorationElement) UpdateProps(data base.Widget) error {
	// Forward to the platform-dependant code
	return w.updateProps(data.(*Decoration))
}
