package goey

var (
	alignKind = Kind{"bitbucket.org/rj/goey.Align"}
)

type Center struct {
	WidthFactor  float64
	HeightFactor float64
	Child        Widget
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*Center) Kind() *Kind {
	return &alignKind
}

// Mount creates a button in the GUI.  The newly created widget
// will be a child of the widget specified by parent.
func (w *Center) Mount(parent Control) (Element, error) {
	child, err := w.Child.Mount(parent)
	if err != nil {
		return nil, err
	}

	return &alignElement{
		parent:       parent,
		child:        child,
		widthFactor:  w.WidthFactor,
		heightFactor: w.HeightFactor,
	}, nil
}

type alignElement struct {
	parent       Control
	child        Element
	childSize    Size
	widthFactor  float64
	heightFactor float64
}

func (w *alignElement) Close() {
	if w.child != nil {
		w.child.Close()
		w.child = nil
	}
}

func (*alignElement) Kind() *Kind {
	return &alignKind
}

func (w *alignElement) Layout(bc Box) Size {
	shrinkWrapWidth := w.widthFactor > 0 || !bc.HasBoundedWidth()
	shrinkWrapHeight := w.heightFactor > 0 || !bc.HasBoundedHeight()

	if w.child == nil {

		size := Size{}
		if !shrinkWrapWidth {
			size.Width = Inf
		}
		if !shrinkWrapHeight {
			size.Height = Inf
		}
		return bc.Constrain(size)
	}

	size := w.child.Layout(bc.Loosen())
	w.childSize = size
	if shrinkWrapWidth && w.widthFactor > 0 {
		size.Width = Length(float64(size.Width) * w.widthFactor)
	}
	if shrinkWrapHeight && w.heightFactor > 0 {
		size.Height = Length(float64(size.Height) * w.heightFactor)
	}
	return bc.Constrain(size)
}

func (w *alignElement) MinimumSize() Size {
	if w.child == nil {
		return Size{}
	}

	return w.child.MinimumSize()
}

func (w *alignElement) SetBounds(bounds Rectangle) {
	if w.child == nil {
		return
	}

	x := bounds.Min.X + (bounds.Dx()-w.childSize.Width)/2
	y := bounds.Min.Y + (bounds.Dy()-w.childSize.Height)/2
	w.child.SetBounds(Rectangle{Point{x, y}, Point{x + w.childSize.Width, y + w.childSize.Height}})
}

func (w *alignElement) updateProps(data *Center) (err error) {
	w.child, err = DiffChild(w.parent, w.child, data.Child)
	w.widthFactor = data.WidthFactor
	w.heightFactor = data.HeightFactor
	return err
}

func (w *alignElement) UpdateProps(data Widget) error {
	return w.updateProps(data.(*Center))
}
