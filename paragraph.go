package goey

var (
	paragraphKind = WidgetKind{"paragraph"}
)

type TextAlignment uint8

const (
	DefaultAlign = TextAlignment(iota)
	Left
	Center
	Right
	Justify
)

type P struct {
	Text  string
	Align TextAlignment
}

func (_ *P) Kind() *WidgetKind {
	return &paragraphKind
}

func (_ *mountedP) Kind() *WidgetKind {
	return &paragraphKind
}
