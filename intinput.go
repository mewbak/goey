package goey

var (
	intInputKind = Kind{"bitbucket.org/rj/goey.IntInput"}
)

// IntInput describes a widget that users input or update a single integer value.
// The model for the value is a int64.
type IntInput struct {
	Value       int64              // Values is the current string for the field
	Placeholder string             // Placeholder is a descriptive text that can be displayed when the field is empty
	Disabled    bool               // Disabled is a flag indicating that the user cannot interact with this field
	OnChange    func(value int64)  // OnChange will be called whenever the user changes the value for this field
	OnFocus     func()             // OnFocus will be called whenever the field receives the keyboard focus
	OnBlur      func()             // OnBlur will be called whenever the field loses the keyboard focus
	OnEnterKey  func(value string) // OnEnterKey will be called whenever the use hits the enter key
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*IntInput) Kind() *Kind {
	return &intInputKind
}

// Mount creates a text field in the GUI.  The newly created widget
// will be a child of the widget specified by parent.
func (w *IntInput) Mount(parent Control) (Element, error) {
	// Forward to the platform-dependant code
	return w.mount(parent)
}

func (*intinputElement) Kind() *Kind {
	return &intInputKind
}

func (w *intinputElement) UpdateProps(data Widget) error {
	return w.updateProps(data.(*IntInput))
}
