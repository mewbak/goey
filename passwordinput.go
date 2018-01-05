package goey

var (
	passwordInputKind = WidgetKind{"password"}
)

// TextInput describes a widget that users input or update a single line of text.
// The model for the value is a string value.
type PasswordInput struct {
	Value       string             // Values is the current string for the field
	Placeholder string             // Placeholder is a descriptive text that can be displayed when the field is empty
	Disabled    bool               // Disabled is a flag indicating that the user cannot interact with this field
	OnChange    func(value string) // OnChange will be called whenever the user changes the value for this field
	OnFocus     func()             // OnFocus will be called whenever the field receives the keyboard focus
	OnBlur      func()             // OnBlur will be called whenever the field loses the keyboard focus
	OnEnterKey  func(value string) // OnEnterKey will be called whenever the use hits the enter key
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*PasswordInput) Kind() *WidgetKind {
	return &passwordInputKind
}

// Mount creates a text field in the GUI.  The newly created widget
// will be a child of the widget specified by parent.
func (w *PasswordInput) Mount(parent NativeWidget) (MountedWidget, error) {
	// Forward to the platform-dependant code
	return w.mount(parent)
}

func (*mountedPasswordInput) Kind() *WidgetKind {
	return &passwordInputKind
}

func (w *mountedPasswordInput) UpdateProps(data Widget) error {
	return w.updateProps(data.(*PasswordInput))
}
