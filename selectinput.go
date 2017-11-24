package goey

var (
	selectKind = WidgetKind{"selectinput"}
)

type SelectInput struct {
	Items    []string
	Value    int
	Unset    bool
	Disabled bool
	OnChange func(value int)
	OnFocus  func()
	OnBlur   func()
}

func (_ *SelectInput) Kind() *WidgetKind {
	return &selectKind
}

func (_ *mountedSelectInput) Kind() *WidgetKind {
	return &selectKind
}
