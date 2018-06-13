package goey

var (
	hboxKind = Kind{"hbox"}
)

// HBox describes a layout widget that arranges its child widgets into a horizontal row.
type HBox struct {
	Children   []Widget
	AlignMain  MainAxisAlign
	AlignCross CrossAxisAlign
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*HBox) Kind() *Kind {
	return &hboxKind
}

// Mount creates a horiztonal layout for child widgets in the GUI.
// The newly created widget will be a child of the widget specified by parent.
func (w *HBox) Mount(parent Control) (Element, error) {
	// Forward to the platform-dependant code
	return w.mount(parent)
}

func (*mountedHBox) Kind() *Kind {
	return &hboxKind
}

type mountedHBox struct {
	parent     Control
	children   []Element
	alignMain  MainAxisAlign
	alignCross CrossAxisAlign

	minimumWidth Length
	maximumWidth Length
}

func (w *HBox) mount(parent Control) (Element, error) {
	c := make([]Element, 0, len(w.Children))

	for _, v := range w.Children {
		mountedChild, err := v.Mount(parent)
		if err != nil {
			return nil, err
		}
		c = append(c, mountedChild)
	}

	return &mountedHBox{parent: parent, children: c,
		alignMain:  w.AlignMain,
		alignCross: w.AlignCross,
	}, nil
}

func (w *mountedHBox) Close() {
	// On this platform, the mountedHBox handles layout, but does not actually
	// have an HWND, so there are no direct resources to release.

	// However, still need to free the children
	for _, v := range w.children {
		v.Close()
	}
	w.children = nil
}

func (w *mountedHBox) MeasureWidth() (Length, Length) {
	if len(w.children) == 0 {
		return 0, 0
	}

	previous := w.children[0]
	min, max := previous.MeasureWidth()
	for _, v := range w.children[1:] {
		gap := calculateHGap(previous, v)
		previous = v
		tmpMin, tmpMax := previous.MeasureWidth()

		min = min + tmpMin + gap
		max = max + tmpMax + gap
	}
	w.minimumWidth = min
	w.maximumWidth = max
	return min, max
}

func (w *mountedHBox) MeasureHeight(width Length) (Length, Length) {
	if len(w.children) == 0 {
		return 0, 0
	}

	if w.minimumWidth == 0 {
		w.MeasureWidth()
		if w.minimumWidth == 0 {
			return 0, 0
		}
	}

	scale1, scale2 := Length(0), Length(1)
	if width > w.minimumWidth && w.maximumWidth > w.minimumWidth {
		scale1, scale2 = width-w.minimumWidth, w.maximumWidth-w.minimumWidth
	}

	minWidth, maxWidth := w.children[0].MeasureWidth()
	childWidth := (minWidth + (maxWidth-minWidth)*scale1/scale2)
	min, max := w.children[0].MeasureHeight(childWidth)
	for _, v := range w.children[1:] {
		minWidth, maxWidth = v.MeasureWidth()
		childWidth := (minWidth + (maxWidth-minWidth)*scale1/scale2)
		tmpMin, tmpMax := v.MeasureHeight(childWidth)
		if tmpMin > min {
			min = tmpMin
		}
		if tmpMax > max {
			max = tmpMax
		}
	}
	return min, max
}

func (w *mountedHBox) SetBounds(bounds Rectangle) {
	if len(w.children) == 0 {
		return
	}

	posX := bounds.Min.X
	width := bounds.Dx()

	if w.minimumWidth == 0 {
		w.MeasureWidth()
		if w.minimumWidth == 0 {
			return
		}
	}

	// If there is more space than necessary, then we need to distribute the extra space.
	extraGap := Length(0)
	if width >= w.maximumWidth {
		switch w.alignMain {
		case MainStart:
			// No need to do any adjustment.  The algorithm below will lay out
			// controls aligned to the top.
		case MainCenter:
			// Adjust the starting position to align the contents.
			posX += (width - w.maximumWidth) / 2

		case MainEnd:
			// Adjust the starting position to align the contents.
			posX += width - w.maximumWidth

		case SpaceAround:
			extraGap = (width - w.maximumWidth).Scale(1, len(w.children)+1)
			posX += extraGap

		case SpaceBetween:
			if len(w.children) > 1 {
				extraGap = (width - w.maximumWidth).Scale(1, len(w.children)-1)
			} else {
				// There are no controls between which to put the extra space.
				// The following essentially convert SpaceBetween to SpaceAround
				extraGap = (width - w.maximumWidth).Scale(1, len(w.children)+1)
				posX += extraGap
			}
		}

		// Reduce available height
		width = w.maximumWidth
	}

	scale1, scale2 := Length(0), Length(1)
	if width > w.minimumWidth && w.maximumWidth > w.minimumWidth {
		scale1, scale2 = width-w.minimumWidth, w.maximumWidth-w.minimumWidth
	}

	previous := Element(nil)
	for _, v := range w.children {
		if previous != nil {
			posX += calculateHGap(previous, v)
		}
		minWidth, maxWidth := v.MeasureWidth()
		childWidth := (minWidth + (maxWidth-minWidth)*scale1/scale2)
		v.SetBounds(Rectangle{Point{posX, bounds.Min.Y}, Point{posX + childWidth, bounds.Max.Y}})
		posX += childWidth + extraGap
		previous = v
	}
}

func (w *mountedHBox) setChildren(children []Widget) error {
	err := error(nil)
	w.children, err = DiffChildren(w.parent, w.children, children)
	return err
}

func (w *mountedHBox) updateProps(data *HBox) error {
	w.alignMain = data.AlignMain
	w.alignCross = data.AlignCross
	return w.setChildren(data.Children)
}

func (w *mountedHBox) UpdateProps(data Widget) error {
	return w.updateProps(data.(*HBox))
}
