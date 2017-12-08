package goey

var (
	vboxKind = WidgetKind{"vbox"}
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
func (*VBox) Kind() *WidgetKind {
	return &vboxKind
}

// Mount creates a vertical layout for child widgets in the GUI.
// The newly created widget will be a child of the widget specified by parent.
func (w *VBox) Mount(parent NativeWidget) (MountedWidget, error) {
	// Forward to the platform-dependant code
	return w.mount(parent)
}

func (*mountedVBox) Kind() *WidgetKind {
	return &vboxKind
}

func (w *mountedVBox) UpdateProps(data Widget) error {
	return w.updateProps(data.(*VBox))
}

func diffChildren(parent NativeWidget, lhs []MountedWidget, rhs []Widget) ([]MountedWidget, error) {

	// If the new tree does not contain any children, then we can trivially
	// match the tree by deleting the actual widgets.
	if len(rhs) == 0 {
		for _, v := range lhs {
			v.Close()
		}
		return nil, nil
	}

	// If the old tree does not contain any children, then we can trivially
	// match the tree by mounting all of the widgets.
	if len(lhs) == 0 && len(rhs) > 0 {
		c := make([]MountedWidget, 0, len(rhs))

		for _, v := range rhs {
			mountedChild, err := v.Mount(parent)
			if err != nil {
				return nil, err
			}
			c = append(c, mountedChild)
		}

		return c, nil
	}

	// Delete excessive children
	if len(lhs) > len(rhs) {
		for _, v := range lhs[len(rhs):] {
			v.Close()
		}
		lhs = lhs[:len(rhs)]
	}

	// Update kind (if necessary) and properties for all of the currently
	// existing children.
	for i := range lhs {
		if kind1, kind2 := lhs[i].Kind(), rhs[i].Kind(); kind1 == kind2 {
			err := lhs[i].UpdateProps(rhs[i])
			if err != nil {
				return lhs, err
			}
		} else {
			mountedWidget, err := rhs[i].Mount(parent)
			if err != nil {
				return lhs, err
			}
			lhs[i].Close()
			lhs[i] = mountedWidget
		}
	}

	// Mount any remaining children.
	for _, v := range rhs[len(lhs):] {
		mountedWidget, err := v.Mount(parent)
		if err != nil {
			return lhs, err
		}
		lhs = append(lhs, mountedWidget)
	}

	return lhs, nil
}
