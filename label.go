package goey

var (
	labelKind = WidgetKind{"label"}
)

type Label struct {
	Text string
}

func (_ *Label) Kind() *WidgetKind {
	return &labelKind
}

func (_ *MountedLabel) Kind() *WidgetKind {
	return &labelKind
}
