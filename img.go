package goey

import (
	"errors"
	"image"
)

var (
	imgKind         = Kind{"bitbucket.org/rj/goey.Img"}
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
func (*Img) Kind() *Kind {
	return &imgKind
}

// Mount creates a horiztonal layout for child widgets in the GUI.
// The newly created widget will be a child of the widget specified by parent.
func (w *Img) Mount(parent Control) (Element, error) {
	// It is an error use to a null image
	if w.Image == nil {
		return nil, ErrInvalidImage
	}

	// Fill in the height and width if they are left at zero.
	w.UpdateDimensions()

	// Forward to the platform-dependant code
	return w.mount(parent)
}

func (_ *mountedImg) Kind() *Kind {
	return &imgKind
}

func (w *mountedImg) Layout(bc Constraint) Size {
	// Determine ideal width.
	return bc.ConstrainAndAttemptToPreserveAspectRatio(Size{w.width, w.height})
}

func (w *mountedImg) MinimumSize() Size {
	// Determine ideal width.
	return Size{w.width, w.height}
}

// UpdateDimensions calculates default values for Width and Height if either
// or zero based on the image dimensions.  The member Image cannot be nil.
func (w *Img) UpdateDimensions() {
	if w.Width == 0 && w.Height == 0 {
		bounds := w.Image.Bounds()
		// Assume that images are at 92 pixels per inch
		w.Width = ((1 * DIP) * 92 / 96).Scale(bounds.Dx(), 1)
		w.Height = ((1 * DIP) * 92 / 96).Scale(bounds.Dy(), 1)
	} else if w.Width == 0 {
		bounds := w.Image.Bounds()
		w.Width = w.Height.Scale(bounds.Dx(), bounds.Dy())
	} else if w.Height == 0 {
		bounds := w.Image.Bounds()
		w.Height = w.Width.Scale(bounds.Dy(), bounds.Dx())
	}
}

func (w *mountedImg) UpdateProps(data Widget) error {
	// Forward to the platform-dependant code
	return w.updateProps(data.(*Img))
}
