package goey

var (
	emptyKind = Kind{"bitbucket.org/rj/goey.Empty"}
)

// Empty describes a widget that is either a horizontal or vertical gap.
type Empty struct {
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*Empty) Kind() *Kind {
	return &emptyKind
}

// Mount creates a horiztonal layout for child widgets in the GUI.
// The newly created widget will be a child of the widget specified by parent.
func (w *Empty) Mount(parent Control) (Element, error) {
	// Forward to the platform-dependant code
	return w.mount(parent)
}

func (*mountedEmpty) Kind() *Kind {
	return &emptyKind
}

func (w *mountedEmpty) MeasureWidth() (Length, Length) {
	return 13 * DIP, 13 * DIP
}

func (w *mountedEmpty) MeasureHeight(width Length) (Length, Length) {
	// Same as static text
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	return 13 * DIP, 13 * DIP
}

func (w *mountedEmpty) UpdateProps(data Widget) error {
	// This widget does not have any properties, so there cannot be anything
	// to update.
	return nil
}
