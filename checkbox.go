package goey

var (
	checkboxKind = WidgetKind{"checkbox"}
)

type Checkbox struct {
	Text     string
	Value    bool
	Disabled bool
	OnChange func(value bool)
	OnFocus  func()
	OnBlur   func()
}

func (_ *Checkbox) Kind() *WidgetKind {
	return &checkboxKind
}

func (_ *mountedCheckbox) Kind() *WidgetKind {
	return &checkboxKind
}
