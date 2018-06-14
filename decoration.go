package goey

import (
	"image/color"
)

var (
	decorationKind = Kind{"bitbucket.org/rj/goey.Decoration"}
)

// Padding describes the interior padding for a decoration widget.
type Padding struct {
	Top, Right, Bottom, Left Length
}

// EvenPadding is a utility to function to build a Padding with the same length for all sides.
func EvenPadding(p Length) Padding {
	return Padding{p, p, p, p}
}

// Decoration describes a widget that provides a border and background.
type Decoration struct {
	Fill    color.RGBA
	Stroke  color.RGBA
	Padding Padding
	Radius  Length
	Child   Widget
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

func (*mountedDecoration) Kind() *Kind {
	return &decorationKind
}

func (w *mountedDecoration) UpdateProps(data Widget) error {
	// Forward to the platform-dependant code
	return w.updateProps(data.(*Decoration))
}
