package goey

import (
	"image"
)

var (
	imgKind = WidgetKind{"img"}
)

// Img describes a widget that contains a bitmap image.
type Img struct {
	Image         image.Image
	Width, Height DIP
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*Img) Kind() *WidgetKind {
	return &imgKind
}

// Mount creates a horiztonal layout for child widgets in the GUI.
// The newly created widget will be a child of the widget specified by parent.
func (w *Img) Mount(parent NativeWidget) (MountedWidget, error) {
	// Forward to the platform-dependant code
	return w.mount(parent)
}

func (_ *mountedImg) Kind() *WidgetKind {
	return &imgKind
}

func (w *mountedImg) UpdateProps(data Widget) error {
	// Forward to the platform-dependant code
	return w.updateProps(data.(*Img))
}
