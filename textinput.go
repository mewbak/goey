package goey

var (
	textInputKind = WidgetKind{"textinput"}
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

func (_ *mountedTextInput) Kind() *WidgetKind {
	return &textInputKind
}
