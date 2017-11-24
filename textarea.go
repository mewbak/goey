package goey

var (
	textareaKind = WidgetKind{"textarea"}
)

type TextArea struct {
	Text        string
	Placeholder string
	Disabled    bool
	MinLines    int
	OnChange    func(value string)
	OnFocus     func()
	OnBlur      func()
}

func (_ *TextArea) Kind() *WidgetKind {
	return &textareaKind
}

func (_ *mountedTextArea) Kind() *WidgetKind {
	return &textareaKind
}
