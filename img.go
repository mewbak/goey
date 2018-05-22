package goey

import (
	"errors"
	"image"
)

var (
	imgKind         = WidgetKind{"img"}
	ErrInvalidImage = errors.New("invalid image property")
)

// Img describes a widget that contains a bitmap image.
// Width and Height may be left as zero.  If both are zero,
// then the image will be sized as if its resolution is 92 DPI.
// Otherwise, if only one is zero, it will be calculated to
// maintain the aspect ratio of the image.
type Img struct {
	Image         image.Image
	Width, Height Length
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (*Img) Kind() *WidgetKind {
	return &imgKind
}

// Mount creates a horiztonal layout for child widgets in the GUI.
// The newly created widget will be a child of the widget specified by parent.
func (w *Img) Mount(parent NativeWidget) (MountedWidget, error) {
	// It is an error use to a null image
	if w.Image == nil {
		return nil, ErrInvalidImage
	}

	// Fill in the height and width if they are left at zero.
	if w.Width == 0 {
		bounds := w.Image.Bounds()
		if w.Height == 0 {
			w.Width = Length(bounds.Dx()*96) / 92
			w.Height = Length(bounds.Dy()*96) / 92
		} else {
			w.Width = (w.Height * Length(bounds.Dx())) / Length(bounds.Dy())
		}
	} else {
		if w.Height == 0 {
			bounds := w.Image.Bounds()
			w.Height = (w.Width * Length(bounds.Dy())) / Length(bounds.Dx())
		}
	}
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
