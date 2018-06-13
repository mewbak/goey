package goey

var (
	buttonKind = Kind{"button"}
)

// Button describes a widget that users can click to initiate an action.
type Button struct {
	Text     string // Text is a caption for the button.
	Disabled bool   // Disabled is a flag indicating that the user cannot interact with this button
	Default  bool   // Default is a flag indicating that the button represents the default action for the interface
	OnClick  func() // OnClick will be called whenever the user presses the button
	OnFocus  func() // OnFocus will be called whenever the button receives the keyboard focus
	OnBlur   func() // OnBlur will be called whenever the button loses the keyboard focus
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*Button) Kind() *Kind {
	return &buttonKind
}

// Mount creates a button in the GUI.  The newly created widget
// will be a child of the widget specified by parent.
func (w *Button) Mount(parent Control) (Element, error) {
	// Forward to the platform-dependant code
	return w.mount(parent)
}

func (*mountedButton) Kind() *Kind {
	return &buttonKind
}

func (w *mountedButton) UpdateProps(data Widget) error {
	// Forward to the platform-dependant code
	return w.updateProps(data.(*Button))
}
