package goey

var (
	emptyKind = WidgetKind{"empty"}
)

// Empty describes a widget that is either a horizontal or vertical gap.
type Empty struct {
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*Empty) Kind() *WidgetKind {
	return &emptyKind
}

// Mount creates a horiztonal layout for child widgets in the GUI.
// The newly created widget will be a child of the widget specified by parent.
func (w *Empty) Mount(parent NativeWidget) (MountedWidget, error) {
	// Forward to the platform-dependant code
	return w.mount(parent)
}

func (_ *mountedEmpty) Kind() *WidgetKind {
	return &emptyKind
}

func (w *mountedEmpty) UpdateProps(data Widget) error {
	// This widget does not have any properties, so there cannot be anything
	// to update.
	return nil
}
