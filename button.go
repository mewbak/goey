package goey

var (
	buttonKind = WidgetKind{"button"}
)

// Button describes a widget that users can click to initiate an action.
type Button struct {
	Text     string
	Disabled bool
	Default  bool
	OnClick  func()
	OnFocus  func()
	OnBlur   func()
}

func (_ *Button) Kind() *WidgetKind {
	return &buttonKind
}

func (_ *mountedButton) Kind() *WidgetKind {
	return &buttonKind
}
