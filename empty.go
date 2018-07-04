package goey

var (
	emptyKind = Kind{"bitbucket.org/rj/goey.Empty"}
)

// Empty describes a widget that is either a horizontal or vertical gap.
//
// The size of the control will be a (perhaps platform dependent) spacing
// between controls.  This applies to both the width and height.
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

func (w *mountedEmpty) Layout(bc Constraint) Size {
	width := w.MinIntrinsicWidth(0)
	height := w.MinIntrinsicHeight(0)
	return bc.Constrain(Size{width,height})
}

func (w *mountedEmpty) MinIntrinsicHeight(width Length) Length {
	// Same as static text
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	return 13 * DIP
}

func (w *mountedEmpty) MinIntrinsicWidth(height Length) Length {
	// Same as static text
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	return 13 * DIP
}

func (w *mountedEmpty) SetBounds(bounds Rectangle) {
	// Virtual control, so no resource to resize
}

func (w *mountedEmpty) UpdateProps(data Widget) error {
	// This widget does not have any properties, so there cannot be anything
	// to update.
	return nil
}
