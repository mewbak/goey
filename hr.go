package goey

var (
	hrKind = WidgetKind{"hr"}
)

// HR describes a widget that is a horiztonal separator.
type HR struct {
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*HR) Kind() *WidgetKind {
	return &hrKind
}

// Mount creates a horiztonal layout for child widgets in the GUI.
// The newly created widget will be a child of the widget specified by parent.
func (w *HR) Mount(parent NativeWidget) (MountedWidget, error) {
	// Forward to the platform-dependant code
	return w.mount(parent)
}

func (_ *mountedHR) Kind() *WidgetKind {
	return &hrKind
}

func (w *mountedHR) UpdateProps(data Widget) error {
	// This widget does not have any properties, so there cannot be anything
	// to update.
	return nil
}
