package goey

// Kind identifies the different kind of widgets.  Most widgets have two
// concrete types associated with their behaviour.  First, there is a type with data
// to describe the widget when unmounted.  Second, there is a type with a handle
// to the windowing system when mounted.  Automatic reconciliation of two widget
// trees relies on Kind to match the unmounted and mounted widgets.
type Kind struct {
	name string
}

// Widget is an interface that wraps any type describing a GUI widget or control.
// A Widget can be 'mounted' to instantiate a widget or control in the GUI.
type Widget interface {
	// Kind returns the concrete type's Kind.  The returned value should
	// be constant, and the same for all instances of a concrete type.
	// Users should not need to use this method directly.
	Kind() *Kind
	// Mount creates a widget or control in the GUI.  The newly created widget
	// will be a child of the widget specified by parent.  If non-nil, the returned
	// Element must have a matching kind.
	Mount(parent NativeWidget) (Element, error)
}

// Element is an interface that wrap any type representing an existing
// widget or control in the GUI.
type Element interface {
	// NativeMountedWidget provides platform dependent methods.  These should
	// not be used by client libraries, but exist for the internal implementation
	// of platform dependent code.
	NativeMountedWidget

	// Close removes the widget from the GUI, and frees any associated resources.
	Close()
	// Kind returns the concrete type for the Element.
	// Users should not need to use this method directly.
	Kind() *Kind
	// MeasureWidth returns the minimum and maximum comfortable widths for
	// this widget.
	MeasureWidth() (min Length, max Length)
	// MeasureHeight returns the minimum and maximum comfortable heights for
	// this widget.
	MeasureHeight(width Length) (min Length, max Length)
	// SetBounds updates the position of the widget.
	SetBounds(bounds Rectangle)
	// UpdateProps will update the properties of the widget.  The Kind for
	// the parameter data must match the Kind for the interface.
	UpdateProps(data Widget) error
}
