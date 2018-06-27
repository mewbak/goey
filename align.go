package goey

var (
	alignKind = Kind{"bitbucket.org/rj/goey.Align"}
)

// Alignment represents the position of a widget along one dimension.
type Alignment int16

// Possible values of alignment.
const (
	AlignStart  = Alignment(-32768) // Widget is aligned at the start (left or top).
	AlignCenter = Alignment(0)      // Widget is aligned at the center.
	AlignEnd    = Alignment(0x7fff) // Widget is aligned at the end (right or bottom).
)

// Align describes a widget that aligns its child within its borders.
type Align struct {
	WidthFactor  float64   // If greater than zero, ratio of container width to child width.
	HeightFactor float64   // If greater than zero, ratio of container height to child height.
	HAlign       Alignment // Horizontal alignment of child widget.
	VAlign       Alignment // Vertical alignment of child widget.
	Child        Widget    // Child widget.
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*Align) Kind() *Kind {
	return &alignKind
}

// Mount creates a button in the GUI.  The newly created widget
// will be a child of the widget specified by parent.
func (w *Align) Mount(parent Control) (Element, error) {
	child, err := w.Child.Mount(parent)
	if err != nil {
		return nil, err
	}

	return &alignElement{
		parent:       parent,
		child:        child,
		widthFactor:  w.WidthFactor,
		heightFactor: w.HeightFactor,
		halign:       w.HAlign,
		valign:       w.VAlign,
	}, nil
}

type alignElement struct {
	parent       Control
	child        Element
	childSize    Size
	halign       Alignment
	valign       Alignment
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

func (w *alignElement) Layout(bc Constraint) Size {
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

	x := bounds.Min.X.Scale(int(w.halign)-int(AlignEnd), int(AlignStart)-int(AlignEnd)) +
		(bounds.Max.X-w.childSize.Width).Scale(int(w.halign)-int(AlignStart), int(AlignEnd)-int(AlignStart))
	y := bounds.Min.Y.Scale(int(w.valign)-int(AlignEnd), int(AlignStart)-int(AlignEnd)) +
		(bounds.Max.Y-w.childSize.Height).Scale(int(w.valign)-int(AlignStart), int(AlignEnd)-int(AlignStart))
	w.child.SetBounds(Rectangle{Point{x, y}, Point{x + w.childSize.Width, y + w.childSize.Height}})
}

func (w *alignElement) updateProps(data *Align) (err error) {
	w.child, err = DiffChild(w.parent, w.child, data.Child)
	w.widthFactor = data.WidthFactor
	w.heightFactor = data.HeightFactor
	w.halign = data.HAlign
	w.valign = data.VAlign
	return err
}

func (w *alignElement) UpdateProps(data Widget) error {
	return w.updateProps(data.(*Align))
}
