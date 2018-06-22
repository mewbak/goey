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
	retval := &mountedEmpty{}

	return retval, nil
}

type mountedEmpty struct {
}

func (w *mountedEmpty) Close() {
	// Virtual control, so no resources to release
}

func (*mountedEmpty) Kind() *Kind {
	return &emptyKind
}

func (w *mountedEmpty) Props() Widget {
	return &Empty{}
}

func (w *mountedEmpty) Layout(bc Box) Size {
	// Determine ideal width.
	return bc.Constrain(Size{13 * DIP, 13 * DIP})
}

func (w *mountedEmpty) MinimumSize() Size {
	// Same as static text
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	return Size{13 * DIP, 13 * DIP}
}

func (w *mountedEmpty) SetBounds(bounds Rectangle) {
	// Virtual control, so no resource to resize
}

func (w *mountedEmpty) UpdateProps(data Widget) error {
	// This widget does not have any properties, so there cannot be anything
	// to update.
	return nil
}
