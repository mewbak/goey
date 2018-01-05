package goey

import (
	"bitbucket.org/rj/goey/syscall"
	"image"
	"image/draw"
	"unsafe"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type mountedImg struct {
	handle    *gtk.Image
	imageData []uint8
}

func imageToPixbuf(prop image.Image) (*gdk.Pixbuf, []uint8, error) {
	if img, ok := prop.(*image.RGBA); ok {
		// Need a copy of the pixel data to support the pixbuf
		buffer := append([]byte(nil), img.Pix...)
		pixbuf := syscall.PixbufNewFromBytes(buffer, gdk.COLORSPACE_RGB, true, 8, img.Rect.Dx(), img.Rect.Dy(), img.Stride)
		return pixbuf, buffer, nil
	}

	// Create a new image in RGBA format
	bounds := prop.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, prop, bounds.Min, draw.Src)

	// Create the bitmap
	pixbuf := syscall.PixbufNewFromBytes(img.Pix, gdk.COLORSPACE_RGB, true, 8, img.Rect.Dx(), img.Rect.Dy(), img.Stride)
	return pixbuf, img.Pix, nil
}

func (w *Img) mount(parent NativeWidget) (MountedWidget, error) {
	// Create the bitmap
	pixbuf, buffer, err := imageToPixbuf(w.Image)
	if err != nil {
		return nil, err
	}

	handle, err := gtk.ImageNewFromPixbuf(pixbuf)
	if err != nil {
		return nil, err
	}
	(*gtk.Container)(unsafe.Pointer(parent.handle)).Add(handle)
	handle.Show()

	retval := &mountedImg{handle, buffer}
	handle.Connect("destroy", img_onDestroy, retval)

	return retval, nil
}

func img_onDestroy(widget *gtk.Label, mounted *mountedImg) {
	mounted.handle = nil
}

func (w *mountedImg) Close() {
	if w.handle != nil {
		w.handle.Destroy()
		w.handle = nil
	}
}

func (w *mountedImg) Handle() *gtk.Widget {
	return &w.handle.Widget
}

func (w *mountedImg) updateProps(data *Img) error {
	// Create the bitmap
	pixbuf, buffer, err := imageToPixbuf(data.Image)
	if err != nil {
		return err
	}
	w.imageData = buffer
	w.handle.SetFromPixbuf(pixbuf)

	return nil
}
