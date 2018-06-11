package goey

var (
	hboxKind = Kind{"hbox"}
)

// HBox describes a layout widget that arranges its child widgets into a horizontal row.
type HBox struct {
	Children   []Widget
	AlignMain  MainAxisAlign
	AlignCross CrossAxisAlign
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*HBox) Kind() *Kind {
	return &hboxKind
}

// Mount creates a horiztonal layout for child widgets in the GUI.
// The newly created widget will be a child of the widget specified by parent.
func (w *HBox) Mount(parent NativeWidget) (Element, error) {
	// Forward to the platform-dependant code
	return w.mount(parent)
}

func (*mountedHBox) Kind() *Kind {
	return &hboxKind
}

func (w *mountedHBox) UpdateProps(data_ Widget) error {
	data := data_.(*HBox)
	return w.SetChildren(data.Children)
}
