package goey

var (
	labelKind = Kind{"bitbucket.org/rj/goey.Label"}
)

// Label describes a widget that provides a descriptive label for other fields.
type Label struct {
	Text string // Text is the contents of the label
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*Label) Kind() *Kind {
	return &labelKind
}

// Mount creates a label in the GUI.  The newly created widget
// will be a child of the widget specified by parent.
func (w *Label) Mount(parent Control) (Element, error) {
	// Forward to the platform-dependant code
	return w.mount(parent)
}

func (*labelElement) Kind() *Kind {
	return &labelKind
}

func (w *labelElement) UpdateProps(data Widget) error {
	return w.updateProps(data.(*Label))
}
