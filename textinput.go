package goey

var (
	textInputKind = WidgetKind{"textInput"}
)

type TextInput struct {
	Text        string
	Placeholder string
	Disabled    bool
	OnChange    func(value string)
	OnFocus     func()
	OnBlur      func()
	OnEnterKey  func(value string)
}

func (_ *TextInput) Kind() *WidgetKind {
	return &textInputKind
}

func (_ *MountedTextInput) Kind() *WidgetKind {
	return &textInputKind
}
