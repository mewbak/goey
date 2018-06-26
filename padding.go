package goey

var (
	paddingKind = Kind{"bitbucket.org/rj/goey.Padding"}
)

type Insets struct {
	Top    Length
	Right  Length
	Bottom Length
	Left   Length
}

func DefaultInsets() Insets {
	const padding = 11 * DIP
	return Insets{padding, padding, padding, padding}
}

func UniformInset(l Length) Insets {
	return Insets{l, l, l, l}
}

type Padding struct {
	Insets Insets
	Child  Widget
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*Padding) Kind() *Kind {
	return &paddingKind
}

// Mount creates a button in the GUI.  The newly created widget
// will be a child of the widget specified by parent.
func (w *Padding) Mount(parent Control) (Element, error) {
	child, err := w.Child.Mount(parent)
	if err != nil {
		return nil, err
	}

	return &paddingElement{
		parent: parent,
		child:  child,
		insets: w.Insets,
	}, nil
}

type paddingElement struct {
	parent    Control
	child     Element
	childSize Size
	insets    Insets
}

func (w *paddingElement) Close() {
	if w.child != nil {
		w.child.Close()
		w.child = nil
	}
}

func (*paddingElement) Kind() *Kind {
	return &paddingKind
}

func (w *paddingElement) Layout(bc Box) Size {
	hinset := w.insets.Left + w.insets.Right
	vinset := w.insets.Top + w.insets.Bottom

	if w.child == nil {
		return bc.Constrain(Size{hinset, vinset})
	}

	innerConstraints := bc.Deflate(hinset, vinset)
	w.childSize = w.child.Layout(innerConstraints)
	return Size{
		w.childSize.Width + hinset,
		w.childSize.Height + vinset,
	}
}

func (w *paddingElement) MinimumSize() Size {
	hinset := w.insets.Left + w.insets.Right
	vinset := w.insets.Top + w.insets.Bottom

	if w.child == nil {
		return Size{hinset, vinset}
	}

	size := w.child.MinimumSize()
	size.Width += hinset
	size.Height += vinset
	return size
}

func (w *paddingElement) SetBounds(bounds Rectangle) {
	if w.child == nil {
		return
	}

	bounds.Min.X += w.insets.Left
	bounds.Min.Y += w.insets.Top
	bounds.Max.X -= w.insets.Right
	bounds.Max.Y -= w.insets.Bottom
	w.child.SetBounds(bounds)
}

func (w *paddingElement) updateProps(data *Padding) (err error) {
	w.child, err = DiffChild(w.parent, w.child, data.Child)
	w.insets = data.Insets
	return err
}

func (w *paddingElement) UpdateProps(data Widget) error {
	return w.updateProps(data.(*Padding))
}
