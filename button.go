package goey

var (
	buttonKind = WidgetKind{"button"}
)

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

func (_ *MountedButton) Kind() *WidgetKind {
	return &buttonKind
}
