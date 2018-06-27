package goey

var (
	checkboxKind = Kind{"bitbucket.org/rj/goey.Checkbox"}
)

// Checkbox describes a widget that users input or update a flag.
// The model for the value is a boolean value.
type Checkbox struct {
	Text     string           // Text is a caption for the checkbox.
	Value    bool             // Is the checkbox checked?
	Disabled bool             // Disabled is a flag indicating that the user cannot interact with this checkbox.
	OnChange func(value bool) // OnChange will be called whenever the value (checked or unchcked) changes.
	OnFocus  func()           // OnFocus will be called whenever the checkbox receives the keyboard focus.
	OnBlur   func()           // OnBlur will be called whenever the checkbox loses the keyboard focus.
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*Checkbox) Kind() *Kind {
	return &checkboxKind
}

// Mount creates a checkbox in the GUI.  The newly created widget
// will be a child of the widget specified by parent.
func (w *Checkbox) Mount(parent Control) (Element, error) {
	// Forward to the platform-dependant code
	return w.mount(parent)
}

func (*mountedCheckbox) Kind() *Kind {
	return &checkboxKind
}

func (w *mountedCheckbox) UpdateProps(data Widget) error {
	return w.updateProps(data.(*Checkbox))
}
