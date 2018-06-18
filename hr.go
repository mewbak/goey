package goey

var (
	hrKind = Kind{"bitbucket.org/rj/goey.HR"}
)

// HR describes a widget that is a horiztonal separator.
type HR struct {
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*HR) Kind() *Kind {
	return &hrKind
}

// Mount creates a horiztonal layout for child widgets in the GUI.
// The newly created widget will be a child of the widget specified by parent.
func (w *HR) Mount(parent Control) (Element, error) {
	// Forward to the platform-dependant code
	return w.mount(parent)
}

func (*mountedHR) Kind() *Kind {
	return &hrKind
}

func (w *mountedHR) Props() Widget {
	return &HR{}
}

func (w *mountedHR) UpdateProps(data Widget) error {
	// This widget does not have any properties, so there cannot be anything
	// to update.
	return nil
}
