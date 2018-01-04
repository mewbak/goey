package goey

var (
	paragraphKind = WidgetKind{"paragraph"}
)

// TextAlignment identifies the different types of text alignment that are possible.
type TextAlignment uint8

const (
	Left    = TextAlignment(iota) // Text aligned to the left (ragged right)
	Center                        // Center alignment
	Right                         // Text aligned to the right (ragged left)
	Justify                       // Text justified so that both left and right are flush
)

// P describes a widget that contains significant text, which can reflow if necessary.
type P struct {
	Text  string        // Text is the content of the paragraph
	Align TextAlignment // Align is the text alignment for the paragraph
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*P) Kind() *WidgetKind {
	return &paragraphKind
}

// Mount creates a paragraph in the GUI.  The newly created widget
// will be a child of the widget specified by parent.
func (w *P) Mount(parent NativeWidget) (MountedWidget, error) {
	// Forward to the platform-dependant code
	return w.mount(parent)
}

func (*mountedP) Kind() *WidgetKind {
	return &paragraphKind
}

func (w *mountedP) UpdateProps(data Widget) error {
	return w.updateProps(data.(*P))
}
