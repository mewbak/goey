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

func (_ *mountedP) Kind() *WidgetKind {
	return &paragraphKind
}
