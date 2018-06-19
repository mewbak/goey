package goey

var (
	hboxKind = Kind{"bitbucket.org/rj/goey.HBox"}
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
	c := make([]Element, 0, len(w.Children))

	for _, v := range w.Children {
		mountedChild, err := v.Mount(parent)
		if err != nil {
			CloseElements(c)
			return nil, err
		}
		c = append(c, mountedChild)
	}

	return &mountedHBox{parent: parent, children: c,
		alignMain:  w.AlignMain,
		alignCross: w.AlignCross,
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

	minimumWidth Length
	maximumWidth Length
}

func (w *mountedHBox) Close() {
	CloseElements(w.children)
	w.children = nil
}

func (w *mountedHBox) MeasureWidth() (Length, Length) {
	if len(w.children) == 0 {
		return 0, 0
	}

	previous := w.children[0]
	min, max := previous.MeasureWidth()
	verifyLengthRange(min, max)
	for _, v := range w.children[1:] {
		gap := calculateHGap(previous, v)
		previous = v
		tmpMin, tmpMax := v.MeasureWidth()
		verifyLengthRange(tmpMin, tmpMax)

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

	scale1, scale2 := 0, 1
	if width > w.minimumWidth && w.maximumWidth > w.minimumWidth {
		scale1, scale2 = int(width-w.minimumWidth), int(w.maximumWidth-w.minimumWidth)
	}

	minWidth, maxWidth := w.children[0].MeasureWidth()
	childWidth := minWidth + (maxWidth-minWidth).Scale(scale1, scale2)
	min, max := w.children[0].MeasureHeight(childWidth)
	verifyLengthRange(min, max)
	for _, v := range w.children[1:] {
		minWidth, maxWidth = v.MeasureWidth()
		verifyLengthRange(minWidth, maxWidth)
		childWidth := minWidth + (maxWidth-minWidth).Scale(scale1, scale2)
		tmpMin, tmpMax := v.MeasureHeight(childWidth)
		verifyLengthRange(tmpMin, tmpMax)
		if tmpMin > min {
			min = tmpMin
		}
		if tmpMax > max {
			max = tmpMax
		}
	}
	println("hbox", "MeasureHeight", min.String(), max.String())
	return min, max
}

func (w *mountedHBox) SetBounds(bounds Rectangle) {
	println("hbox", "SetBounds", bounds.String())

	if len(w.children) == 0 {
		return
	}

	if w.minimumWidth == 0 {
		w.MeasureWidth()
		if w.minimumWidth == 0 {
			return
		}
	}

	extraGap, deltaX, scale1, scale2 := distributeVSpace(w.alignMain, len(w.children), bounds.Dx(), w.minimumWidth, w.maximumWidth)
	bounds.Min.X += deltaX

	previous := Element(nil)
	for _, v := range w.children {
		if previous != nil {
			bounds.Min.X += calculateHGap(previous, v) + extraGap
		}

		minWidth, maxWidth := v.MeasureWidth()
		childWidth := minWidth + (maxWidth-minWidth).Scale(scale1, scale2)
		v.SetBounds(Rectangle{bounds.Min, Point{bounds.Min.X + childWidth, bounds.Max.Y}})
		bounds.Min.X += childWidth
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
