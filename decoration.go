package goey

import (
	"image/color"
)

var (
	decorationKind = Kind{"bitbucket.org/rj/goey.Decoration"}
)

// Decoration describes a widget that provides a border and background, and
// possibly containing a single child widget.
//
// The size of the control will match the size of the child element, although
// padding will be added between the border of the decoration and the child
// element as specified by the field Insets.
type Decoration struct {
	Fill   color.RGBA // Background colour used to fill interior.
	Stroke color.RGBA // Stroke colour used to draw outline.
	Insets Insets     // Space between border of the decoration and the child element.
	Radius Length     // Radius of the widgets corners.
	Child  Widget     // Child.
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*Decoration) Kind() *Kind {
	return &decorationKind
}

// Mount creates a button in the GUI.  The newly created widget
// will be a child of the widget specified by parent.
func (w *Decoration) Mount(parent Control) (Element, error) {
	// Forward to the platform-dependant code
	return w.mount(parent)
}

func (w *decorationElement) Layout(bc Constraint) Size {
	hinset := w.insets.Left + w.insets.Right
	vinset := w.insets.Top + w.insets.Bottom

	if w.child == nil {
		return bc.Constrain(Size{hinset, vinset})
	}

	innerConstraints := bc.Inset(hinset, vinset)
	w.childSize = w.child.Layout(innerConstraints)
	return Size{
		w.childSize.Width + hinset,
		w.childSize.Height + vinset,
	}
}

func (w *decorationElement) MinimumSize() Size {
	hinset := w.insets.Left + w.insets.Right
	vinset := w.insets.Top + w.insets.Bottom

	if w.child == nil {
		return Size{hinset, vinset}
	}

	size := w.child.MinimumSize()
	size.Width += hinset
	size.Height += vinset
	return size
}

func (*decorationElement) Kind() *Kind {
	return &decorationKind
}

func (w *decorationElement) UpdateProps(data Widget) error {
	// Forward to the platform-dependant code
	return w.updateProps(data.(*Decoration))
}
