package goey

import (
	"goey/syscall"
	"image"
	"unsafe"

	"github.com/gotk3/gotk3/gtk"

	"github.com/gotk3/gotk3/gdk"
)

type mountedImg struct {
	handle *gtk.Image
}

func imageToPixbuf(prop image.Image) (*gdk.Pixbuf, error) {
	if img, ok := prop.(*image.RGBA); ok {
		pixbuf := syscall.PixbufNewFromBytes(img.Pix, gdk.COLORSPACE_RGB, true, 8, img.Rect.Dx(), img.Rect.Dy(), img.Stride)
		return pixbuf, nil
	} else {
		panic("Unsupported image format.")
	}
}

func (w *Img) mount(parent NativeWidget) (MountedWidget, error) {
	// Create the bitmap
	pixbuf, err := imageToPixbuf(w.Image)
	if err != nil {
		return nil, err
	}

	handle, err := gtk.ImageNewFromPixbuf(pixbuf)
	if err != nil {
		return nil, err
	}
	(*gtk.Container)(unsafe.Pointer(parent.handle)).Add(handle)
	handle.Show()

	retval := &mountedImg{handle}
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
	pixbuf, err := imageToPixbuf(data.Image)
	if err != nil {
		return err
	}
	w.handle.SetFromPixbuf(pixbuf)

	return nil
}
