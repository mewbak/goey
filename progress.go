package goey

var (
	progressKind = Kind{"bitbucket.org/rj/goey.Progress"}
)

// Progress describes a widget that shows a progress bar.
// The model for the value is an int.
type Progress struct {
	Value    int
	Min, Max int
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*Progress) Kind() *Kind {
	return &progressKind
}

// Mount creates a text area control in the GUI.  The newly created widget
// will be a child of the widget specified by parent.
func (w *Progress) Mount(parent Control) (Element, error) {
	// Forward to the platform-dependant code
	return w.mount(parent)
}

func (*progressElement) Kind() *Kind {
	return &progressKind
}

func (w *progressElement) UpdateProps(data Widget) error {
	return w.updateProps(data.(*Progress))
}
