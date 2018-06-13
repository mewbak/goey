package goey

var (
	vboxKind = Kind{"vbox"}
)

// MainAxisAlign identifies the different types of alignment that is possible
// along the main axis for a vertical box or horizontal box layout.
type MainAxisAlign uint8

const (
	MainStart    = MainAxisAlign(iota) // Children will be packed at the top or left of the box
	MainCenter                         // Children will be packed together and centered in the box.
	MainEnd                            // Children will be packed together at the bottom or right of the box
	SpaceAround                        // Children will be spaced apart
	SpaceBetween                       // Children will be spaced apart, but the first and last children will but the ends of the box.
)

// CrossAxisAlign identifies the different types of alignment that is possible
// along the cross axis for vertical box and horizontal box layouts.
type CrossAxisAlign uint8

const (
	Stretch     = CrossAxisAlign(iota) // Children will be stretched so that the extend across box
	CrossStart                         // Children will be aligned to the left or top of the box
	CrossCenter                        // Children will be aligned in the center of the box
	CrossEnd                           // Children will be aligned to the right or bottom of the box
)

// VBox describes a vertical layout.  Children are positioned in order from
// the top towards the bottom.  The main axis for alignment is therefore vertical,
// with the cross axis for alignment horiztonal.
type VBox struct {
	Children   []Widget
	AlignMain  MainAxisAlign
	AlignCross CrossAxisAlign
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*VBox) Kind() *Kind {
	return &vboxKind
}

// Mount creates a vertical layout for child widgets in the GUI.
// The newly created widget will be a child of the widget specified by parent.
func (w *VBox) Mount(parent NativeWidget) (Element, error) {
	// Forward to the platform-dependant code
	return w.mount(parent)
}

func (*mountedVBox) Kind() *Kind {
	return &vboxKind
}

type mountedVBox struct {
	parent     NativeWidget
	children   []Element
	alignMain  MainAxisAlign
	alignCross CrossAxisAlign
}

func (w *VBox) mount(parent NativeWidget) (Element, error) {
	c := make([]Element, 0, len(w.Children))

	for _, v := range w.Children {
		mountedChild, err := v.Mount(parent)
		if err != nil {
			return nil, err
		}
		c = append(c, mountedChild)
	}

	return &mountedVBox{parent: parent, children: c}, nil
}

func (w *mountedVBox) Close() {
	// Need to free the children
	for _, v := range w.children {
		v.Close()
	}
	w.children = nil
}

func (w *mountedVBox) MeasureWidth() (Length, Length) {
	if len(w.children) == 0 {
		return 0, 0
	}

	min, max := w.children[0].MeasureWidth()
	for _, v := range w.children[1:] {
		tmpMin, tmpMax := v.MeasureWidth()
		if tmpMin > min {
			min = tmpMin
		}
		if tmpMax > max {
			max = tmpMax
		}
	}
	return min, max
}

func (w *mountedVBox) MeasureHeight(width Length) (Length, Length) {
	if len(w.children) == 0 {
		return 0, 0
	}

	previous := w.children[0]
	min, max := previous.MeasureHeight(width)
	for _, v := range w.children[1:] {
		tmpMin, tmpMax := v.MeasureHeight(width)
		gap := calculateVGap(previous, v)
		min += tmpMin + gap
		max += tmpMax + gap
		previous = v
	}
	return min, max
}

func (w *mountedVBox) SetBounds(bounds Rectangle) {
	if len(w.children) == 0 {
		return
	}

	width := bounds.Dx()
	minTotal, maxTotal := w.MeasureHeight(width)

	extraGap, deltaY, scale1, scale2 := distributeVSpace(w.alignMain, len(w.children), bounds.Dy(), minTotal, maxTotal)
	bounds.Min.Y += deltaY

	// Assuming that height of bounds is sufficient
	previous := Element(nil)
	for _, v := range w.children {
		if previous != nil {
			bounds.Min.Y += calculateVGap(previous, v) + extraGap
		}

		deltaY := setBoundsWithAlign(v, bounds, w.alignCross, scale1, scale2)
		bounds.Min.Y += deltaY
		previous = v
	}
}

func (w *mountedVBox) setChildren(children []Widget) error {
	err := error(nil)
	w.children, err = diffChildren(w.parent, w.children, children)
	return err
}

func (w *mountedVBox) updateProps(data *VBox) error {
	w.alignMain = data.AlignMain
	w.alignCross = data.AlignCross
	return w.setChildren(data.Children)
}

func (w *mountedVBox) UpdateProps(data Widget) error {
	return w.updateProps(data.(*VBox))
}
