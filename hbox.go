package goey

var (
	hboxKind = Kind{"bitbucket.org/rj/goey.HBox"}
)

// HBox describes a layout widget that arranges its child widgets into a horizontal row.
type HBox struct {
	AlignMain  MainAxisAlign
	AlignCross CrossAxisAlign
	Children   []Widget
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*HBox) Kind() *Kind {
	return &hboxKind
}

// Mount creates a horiztonal layout for child widgets in the GUI.
// The newly created widget will be a child of the widget specified by parent.
func (w *HBox) Mount(parent Control) (Element, error) {
	c := make([]Element, 0, len(w.Children))

	for _, v := range w.Children {
		mountedChild, err := v.Mount(parent)
		if err != nil {
			CloseElements(c)
			return nil, err
		}
		c = append(c, mountedChild)
	}

	return &mountedHBox{
		parent:       parent,
		children:     c,
		alignMain:    w.AlignMain,
		alignCross:   w.AlignCross,
		childrenSize: make([]Size, len(c)),
	}, nil
}

func (*mountedHBox) Kind() *Kind {
	return &hboxKind
}

type mountedHBox struct {
	parent     Control
	children   []Element
	alignMain  MainAxisAlign
	alignCross CrossAxisAlign

	childrenSize []Size
	totalWidth   Length
	minimumSize  Size
}

func (w *mountedHBox) Close() {
	CloseElements(w.children)
	w.children = nil
}

func (w *mountedHBox) Layout(bc Box) Size {
	if len(w.children) == 0 {
		w.totalWidth = 0
		return Size{}
	}

	// Determine the constraints for layout of child elements.
	cbc := bc.LoosenWidth()
	if w.alignCross == Stretch {
		if cbc.HasBoundedHeight() {
			cbc = cbc.TightenHeight(cbc.Max.Height)
		} else {
			size := w.MinimumSize()
			cbc = cbc.TightenHeight(max(cbc.Min.Height, size.Height))
		}
	} else {
		cbc = cbc.LoosenHeight()
	}

	width := Length(0)
	minHeight := Length(0)
	previous := Element(nil)
	for i, v := range w.children {
		if i > 0 {
			if w.alignMain.IsPacked() {
				width += calculateHGap(previous, v)
			}
			previous = v
		}
		w.childrenSize[i] = v.Layout(cbc)
		minHeight = max(minHeight, w.childrenSize[i].Height)
		width += w.childrenSize[i].Width
	}
	w.totalWidth = width

	if w.alignCross == Stretch {
		return bc.Constrain(Size{width, cbc.Min.Height})
	}
	return bc.Constrain(Size{width, minHeight})
}

func (w *mountedHBox) MinimumSize() Size {
	if len(w.children) == 0 {
		return Size{}
	}

	if !w.minimumSize.IsZero() {
		return w.minimumSize
	}

	size := w.children[0].MinimumSize()
	if w.alignMain.IsPacked() {
		previous := w.children[0]
		for _, v := range w.children[1:] {
			// Add the preferred gap between this pair of widgets
			size.Width += calculateHGap(previous, v)
			// Find minimum size for this widget, and update
			tmp := v.MinimumSize()
			size.Height = max(size.Height, tmp.Height)
			size.Width += tmp.Width
		}
	} else {
		for _, v := range w.children[1:] {
			// Find minimum size for this widget, and update
			tmp := v.MinimumSize()
			size.Height = max(size.Height, tmp.Height)
			size.Width += tmp.Width
		}

		// Add a minimum gap between the controls.
		if w.alignMain == SpaceBetween {
			size.Width += calculateHGap(nil, nil).Scale(len(w.children)-1, 1)
		} else {
			size.Width += calculateHGap(nil, nil).Scale(len(w.children)+1, 1)
		}
	}

	w.minimumSize = size
	return size
}

func (w *mountedHBox) SetBounds(bounds Rectangle) {
	// Adjust the bounds so that the minimum Y handles vertical alignment
	// of the controls.  We also calculate 'extraGap' which will adjust
	// spacing of the controls for non-packed alignments.
	extraGap := Length(0)
	switch w.alignMain {
	case MainStart:
		// Do nothing
	case MainCenter:
		bounds.Min.X += (bounds.Dx() - w.totalWidth) / 2
	case MainEnd:
		bounds.Min.X = bounds.Max.X - w.totalWidth
	case SpaceAround:
		extraGap = (bounds.Dx() - w.totalWidth).Scale(1, len(w.children)+1)
		bounds.Min.X += extraGap
	case SpaceBetween:
		if len(w.children) > 1 {
			extraGap = (bounds.Dx() - w.totalWidth).Scale(1, len(w.children)-1)
		} else {
			// There are no controls between which to put the extra space.
			// The following essentially convert SpaceBetween to SpaceAround
			bounds.Min.X += (bounds.Dx() - w.totalWidth) / 2
		}
	}

	// Position all of the child controls.
	posX := bounds.Min.X
	previous := Element(nil)
	for i, v := range w.children {
		if w.alignMain.IsPacked() {
			if i > 0 {
				posX += calculateHGap(previous, v)
			}
			previous = v
		}

		dx := w.childrenSize[i].Width
		dy := w.childrenSize[i].Height
		switch w.alignCross {
		case CrossStart:
			v.SetBounds(Rectangle{
				Point{posX, bounds.Min.Y},
				Point{posX + dx, bounds.Min.Y + dy},
			})
		case CrossCenter:
			v.SetBounds(Rectangle{
				Point{posX, bounds.Min.Y + (bounds.Dy()-dy)/2},
				Point{posX + dx, bounds.Min.Y + (bounds.Dy()+dy)/2},
			})
		case CrossEnd:
			v.SetBounds(Rectangle{
				Point{posX, bounds.Max.Y - dy},
				Point{posX + dx, bounds.Max.Y},
			})
		case Stretch:
			v.SetBounds(Rectangle{
				Point{posX, bounds.Min.Y},
				Point{posX + dx, bounds.Max.Y},
			})
		}
		posX += dx + extraGap
	}
}

func (w *mountedHBox) updateProps(data *HBox) (err error) {
	// Update properties
	w.alignMain = data.AlignMain
	w.alignCross = data.AlignCross
	w.children, err = DiffChildren(w.parent, w.children, data.Children)
	// Clear cached values
	w.childrenSize = make([]Size, len(w.children))
	w.totalWidth = 0
	w.minimumSize = Size{}
	return err
}

func (w *mountedHBox) UpdateProps(data Widget) error {
	return w.updateProps(data.(*HBox))
}
