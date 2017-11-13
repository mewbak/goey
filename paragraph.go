package goey

var (
	paragraphKind = WidgetKind{"paragraph"}
)

type Alignment uint8

const (
	DefaultAlign = Alignment(iota)
	Left
	Center
	Right
	Justify
)

type P struct {
	Text  string
	Align Alignment
}

func (_ *P) Kind() *WidgetKind {
	return &paragraphKind
}

func (_ *MountedP) Kind() *WidgetKind {
	return &paragraphKind
}
