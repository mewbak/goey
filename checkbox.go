package goey

var (
	checkboxKind = WidgetKind{"checkbox"}
)

// Checkbox describes a widget that users input or update a flag.
// The model for the value is a boolean value.
type Checkbox struct {
	Text     string
	Value    bool
	Disabled bool
	OnChange func(value bool)
	OnFocus  func()
	OnBlur   func()
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*Checkbox) Kind() *WidgetKind {
	return &checkboxKind
}

// Mount creates a checkbox in the GUI.  The newly created widget
// will be a child of the widget specified by parent.
func (w *Checkbox) Mount(parent NativeWidget) (MountedWidget, error) {
	// Forward to the platform-dependant code
	return w.mount(parent)
}

func (*mountedCheckbox) Kind() *WidgetKind {
	return &checkboxKind
}

func (w *mountedCheckbox) UpdateProps(data Widget) error {
	return w.updateProps(data.(*Checkbox))
}
