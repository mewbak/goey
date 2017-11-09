package goey

type WidgetKind struct {
	name string
}

type Widget interface {
	Kind() *WidgetKind
	Mount(parent NativeWidget) (MountedWidget, error)
}

type MountedWidget interface {
	NativeMountedWidget

	Kind() *WidgetKind
	UpdateProps(data Widget) error
	Close()
}
