package goey

var (
	hboxKind = WidgetKind{"hbox"}
)

type HBox struct {
	Children   []Widget
	AlignMain  MainAxisAlign
	AlignCross CrossAxisAlign
	Grid       bool
}

func (_ *HBox) Kind() *WidgetKind {
	return &hboxKind
}

func (_ *mountedHBox) Kind() *WidgetKind {
	return &hboxKind
}

func (w *mountedHBox) UpdateProps(data_ Widget) error {
	data := data_.(*HBox)
	return w.SetChildren(data.Children)
}
