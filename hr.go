package goey

var (
	hrKind = WidgetKind{"hr"}
)

type HR struct {
}

func (_ *HR) Kind() *WidgetKind {
	return &hrKind
}

func (_ *MountedHR) Kind() *WidgetKind {
	return &hrKind
}
