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

// Kind returns the concrete type for use in the Widget interface.
func (*Button) Kind() *WidgetKind {
	return &buttonKind
}

func (*mountedButton) Kind() *WidgetKind {
	return &buttonKind
}
