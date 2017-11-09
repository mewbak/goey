package goey

var (
	vboxKind = WidgetKind{"vbox"}
)

type VBox struct {
	Children []Widget
}

func (_ *VBox) Kind() *WidgetKind {
	return &vboxKind
}

func (w *VBox) Mount(parent NativeWidget) (MountedWidget, error) {
	c := make([]MountedWidget, 0, len(w.Children))

	for _, v := range w.Children {
		mountedChild, err := v.Mount(parent)
		if err != nil {
			return nil, err
		}
		c = append(c, mountedChild)
	}

	return &MountedVBox{parent: parent, children: c}, nil
}

type MountedVBox struct {
	parent   NativeWidget
	children []MountedWidget
}

func (_ *MountedVBox) Kind() *WidgetKind {
	return &vboxKind
}

func (w *MountedVBox) Close() {
	// This widget does not hold a native control, so there are no resources
	// beyond memory to cleanup.
}

func (w *MountedVBox) UpdateProps(data_ Widget) error {
	data := data_.(*VBox)
	return w.SetChildren(data.Children)
}

func (w *MountedVBox) SetChildren(children []Widget) error {
	err := error(nil)
	w.children, err = diffChildren(w.parent, w.children, children)
	return err
}

func diffChildren(parent NativeWidget, lhs []MountedWidget, rhs []Widget) ([]MountedWidget, error) {

	// If the new tree does not contain any children, then we can trivially
	// match the tree by deleting the actual widgets.
	if len(lhs) == 0 {
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
