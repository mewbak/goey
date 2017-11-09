package goey

var (
	hboxKind = WidgetKind{"hbox"}
)

type HBox struct {
	Children []Widget
	Align    Alignment
}

func (_ *HBox) Kind() *WidgetKind {
	return &hboxKind
}

func (w *HBox) Mount(parent NativeWidget) (MountedWidget, error) {
	c := make([]MountedWidget, 0, len(w.Children))

	for _, v := range w.Children {
		mountedChild, err := v.Mount(parent)
		if err != nil {
			return nil, err
		}
		c = append(c, mountedChild)
	}

	align := w.Align
	if align == DefaultAlign {
		align = Justify
	}

	return &MountedHBox{parent: parent, children: c, align: align}, nil
}

type MountedHBox struct {
	parent         NativeWidget
	children       []MountedWidget
	align          Alignment
	preferredWidth int
}

func (_ *MountedHBox) Kind() *WidgetKind {
	return &hboxKind
}

func (w *MountedHBox) Close() {
	// This widget does not hold a native control, so there are no resources
	// beyond memory to cleanup.
}

func (w *MountedHBox) UpdateProps(data_ Widget) error {
	panic("not implemented")
}
