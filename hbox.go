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
