package goey

var (
	columnKind = Kind{"column"}
)

// HBox describes a layout widget that arranges its child widgets into a horizontal row.
type Column struct {
	Children [][]Widget
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*Column) Kind() *Kind {
	return &columnKind
}

// Mount creates a horiztonal layout for child widgets in the GUI.
// The newly created widget will be a child of the widget specified by parent.
func (w *Column) Mount(parent NativeWidget) (Element, error) {
	// Forward to the platform-dependant code
	return w.mount(parent)
}

func (*mountedColumn) Kind() *Kind {
	return &columnKind
}

func (w *mountedColumn) UpdateProps(data_ Widget) error {
	data := data_.(*Column)
	return w.SetChildren(data.Children)
}
