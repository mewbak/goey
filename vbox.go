package goey

var (
	vboxKind = Kind{"bitbucket.org/rj/goey.VBox"}
)

// MainAxisAlign identifies the different types of alignment that is possible
// along the main axis for a vertical box or horizontal box layout.
type MainAxisAlign uint8

// Allowed values for alignment of the main axis in a vertical box (VBox) or
// horizontal box (HBox).
const (
	MainStart    = MainAxisAlign(iota) // Children will be packed together at the top or left of the box
	MainCenter                         // Children will be packed together and centered in the box.
	MainEnd                            // Children will be packed together at the bottom or right of the box
	SpaceAround                        // Children will be spaced apart
	SpaceBetween                       // Children will be spaced apart, but the first and last children will but the ends of the box.
	Homogeneous                        // Children will be allocated equal space.
)

// IsPacked returns true if the main axis alignment is a one where children
// will be packed together.
func (a MainAxisAlign) IsPacked() bool {
	return a <= MainEnd
}

// CrossAxisAlign identifies the different types of alignment that is possible
// along the cross axis for vertical box and horizontal box layouts.
type CrossAxisAlign uint8

// Allowed values for alignment of the cross axis in a vertical box (VBox) or
// horizontal box (HBox).
const (
	Stretch     = CrossAxisAlign(iota) // Children will be stretched so that the extend across box
	CrossStart                         // Children will be aligned to the left or top of the box
	CrossCenter                        // Children will be aligned in the center of the box
	CrossEnd                           // Children will be aligned to the right or bottom of the box
)

// VBox describes a layout widget that arranges its child widgets into a column.
// Children are positioned in order from the top towards the bottom.  The main
// axis for alignment is therefore vertical, with the cross axis for alignment is horiztonal.
//
// The size of the box will try to set a width sufficient to contain all of its
// children.  Extra space will be distributed according to the value of
// AlignMain.  Subject to the box constraints during layout, the height should
// match the largest minimum height of the child widgets.
type VBox struct {
	AlignMain  MainAxisAlign
	AlignCross CrossAxisAlign
	Children   []Widget
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*VBox) Kind() *Kind {
	return &vboxKind
}

// Mount creates a vertical layout for child widgets in the GUI.
// The newly created widget will be a child of the widget specified by parent.
func (w *VBox) Mount(parent Control) (Element, error) {
	c := make([]Element, 0, len(w.Children))

	for _, v := range w.Children {
		mountedChild, err := v.Mount(parent)
		if err != nil {
			CloseElements(c)
			return nil, err
		}
		c = append(c, mountedChild)
	}

	return &mountedVBox{
		parent:       parent,
		children:     c,
		alignMain:    w.AlignMain,
		alignCross:   w.AlignCross,
		childrenSize: make([]Size, len(c)),
	}, nil
}

func (*mountedVBox) Kind() *Kind {
	return &vboxKind
}

type mountedVBox struct {
	parent     Control
	children   []Element
	alignMain  MainAxisAlign
	alignCross CrossAxisAlign

	childrenSize []Size
	totalHeight  Length
	minimumSize  Size
}

func (w *mountedVBox) Close() {
	CloseElements(w.children)
	w.children = nil
}

func (w *mountedVBox) Layout(bc Constraint) Size {
	if len(w.children) == 0 {
		w.totalHeight = 0
		return Size{}
	}

	// Determine the constraints for layout of child elements.
	cbc := bc.LoosenHeight()
	if w.alignMain == Homogeneous {
		if count := len(w.children); count > 1 {
			gap := calculateVGap(nil, nil)
			cbc.TightenHeight(cbc.Max.Height.Scale(1, count) - gap.Scale(1, count-1))
		} else {
			cbc.TightenHeight(cbc.Max.Height.Scale(1, count))
		}
	}
	if w.alignCross == Stretch {
		if cbc.HasBoundedWidth() {
			cbc = cbc.TightenWidth(cbc.Max.Width)
		} else {
			size := w.MinimumSize()
			cbc = cbc.TightenWidth(max(cbc.Min.Width, size.Width))
		}
	} else {
		cbc = cbc.LoosenWidth()
	}

	height := Length(0)
	minWidth := Length(0)
	previous := Element(nil)
	for i, v := range w.children {
		if i > 0 {
			if w.alignMain.IsPacked() {
				height += calculateVGap(previous, v)
				previous = v
			}
		}
		w.childrenSize[i] = v.Layout(cbc)
		minWidth = max(minWidth, w.childrenSize[i].Width)
		height += w.childrenSize[i].Height
	}
	w.totalHeight = height

	if w.alignCross == Stretch {
		return bc.Constrain(Size{cbc.Min.Width, height})
	}
	return bc.Constrain(Size{minWidth, height})
}

func (w *mountedVBox) MinimumSize() Size {
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
			size.Height += calculateVGap(previous, v)
			// Find minimum size for this widget, and update
			tmp := v.MinimumSize()
			size.Width = max(size.Width, tmp.Width)
			size.Height += tmp.Height
		}
	} else if w.alignMain == Homogeneous {
		for _, v := range w.children[1:] {
			// Find minimum size for this widget, and update
			tmp := v.MinimumSize()
			size.Height = max(size.Height, tmp.Height)
			size.Width = max(size.Width, tmp.Width)
		}

		// Add a minimum gap between the controls.
		size.Height = size.Height.Scale(len(w.children), 1) + calculateVGap(nil, nil).Scale(len(w.children)-1, 1)
	} else {
		for _, v := range w.children[1:] {
			// Find minimum size for this widget, and update
			tmp := v.MinimumSize()
			size.Width = max(size.Width, tmp.Width)
			size.Height += tmp.Height
		}

		// Add a minimum gap between the controls.
		if w.alignMain == SpaceBetween {
			size.Height += calculateVGap(nil, nil).Scale(len(w.children)-1, 1)
		} else {
			size.Height += calculateVGap(nil, nil).Scale(len(w.children)+1, 1)
		}
	}

	w.minimumSize = size
	return size
}

func (w *mountedVBox) Props() Widget {
	children := []Widget(nil)
	if len(w.children) != 0 {
		children = make([]Widget, 0, len(w.children))
		for _, v := range w.children {
			children = append(children, v.(Proper).Props())
		}
	}

	return &VBox{
		AlignMain:  w.alignMain,
		AlignCross: w.alignCross,
		Children:   children,
	}
}

func (w *mountedVBox) SetBounds(bounds Rectangle) {
	if len(w.children) == 0 {
		return
	}

	if w.alignMain == Homogeneous {
		gap := calculateVGap(nil, nil)
		dy := bounds.Dy() + gap
		count := len(w.children)

		for i, v := range w.children {
			y1 := bounds.Min.Y + dy.Scale(i, count)
			y2 := bounds.Min.Y + dy.Scale(i+1, count) - gap
			w.setBoundsForChild(i, v, bounds.Min.X, y1, bounds.Max.X, y2)
		}
		return
	}

	// Adjust the bounds so that the minimum Y handles vertical alignment
	// of the controls.  We also calculate 'extraGap' which will adjust
	// spacing of the controls for non-packed alignments.
	extraGap := Length(0)
	switch w.alignMain {
	case MainStart:
		// Do nothing
	case MainCenter:
		bounds.Min.Y += (bounds.Dy() - w.totalHeight) / 2
	case MainEnd:
		bounds.Min.Y = bounds.Max.Y - w.totalHeight
	case SpaceAround:
		extraGap = (bounds.Dy() - w.totalHeight).Scale(1, len(w.children)+1)
		bounds.Min.Y += extraGap
	case SpaceBetween:
		if len(w.children) > 1 {
			extraGap = (bounds.Dy() - w.totalHeight).Scale(1, len(w.children)-1)
		} else {
			// There are no controls between which to put the extra space.
			// The following essentially convert SpaceBetween to SpaceAround
			bounds.Min.Y += (bounds.Dy() - w.totalHeight) / 2
		}
	}

	// Position all of the child controls.
	posY := bounds.Min.Y
	previous := Element(nil)
	for i, v := range w.children {
		if w.alignMain.IsPacked() {
			if i > 0 {
				posY += calculateVGap(previous, v)
			}
			previous = v
		}

		dy := w.childrenSize[i].Height
		w.setBoundsForChild(i, v, bounds.Min.X, posY, bounds.Max.X, posY+dy)
		posY += dy + extraGap
	}
}

func (w *mountedVBox) setBoundsForChild(i int, v Element, posX, posY, posX2, posY2 Length) {
	dx := w.childrenSize[i].Width
	switch w.alignCross {
	case CrossStart:
		v.SetBounds(Rectangle{
			Point{posX, posY},
			Point{posX + dx, posY2},
		})
	case CrossCenter:
		v.SetBounds(Rectangle{
			Point{posX + (posX2-posX-dx)/2, posY},
			Point{posX + (posX2-posX+dx)/2, posY2},
		})
	case CrossEnd:
		v.SetBounds(Rectangle{
			Point{posX2 - dx, posY},
			Point{posX2, posY2},
		})
	case Stretch:
		v.SetBounds(Rectangle{
			Point{posX, posY},
			Point{posX2, posY2},
		})
	}
}

func (w *mountedVBox) updateProps(data *VBox) (err error) {
	// Update properties
	w.alignMain = data.AlignMain
	w.alignCross = data.AlignCross
	w.children, err = DiffChildren(w.parent, w.children, data.Children)
	// Clear cached values
	w.childrenSize = make([]Size, len(w.children))
	w.totalHeight = 0
	w.minimumSize = Size{}
	return err
}

func (w *mountedVBox) UpdateProps(data Widget) error {
	return w.updateProps(data.(*VBox))
}
