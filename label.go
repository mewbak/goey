package goey

var (
	labelKind = WidgetKind{"label"}
)

type Alignment uint8

const (
	DefaultAlign = Alignment(iota)
	Left
	Center
	Right
	Justify
)

type Label struct {
	Text  string
	Align Alignment
}

func (_ *Label) Kind() *WidgetKind {
	return &labelKind
}

func (_ *MountedLabel) Kind() *WidgetKind {
	return &labelKind
}
