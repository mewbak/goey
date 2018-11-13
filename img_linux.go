package goey

import (
	"image"
	"image/draw"
	"unsafe"

	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/internal/syscall"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type imgElement struct {
	Control

	imageData []uint8
	width     base.Length
	height    base.Length
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

func pixbufToImage(pixbuf *gdk.Pixbuf) image.Image {
	cs := pixbuf.GetColorspace()
	bps := pixbuf.GetBitsPerSample()
	alpha := pixbuf.GetHasAlpha()

	if cs == gdk.COLORSPACE_RGB && alpha && bps == 8 {
		return &image.RGBA{
			Pix:    pixbuf.GetPixels(),
			Stride: pixbuf.GetRowstride(),
			Rect:   image.Rect(0, 0, pixbuf.GetWidth(), pixbuf.GetHeight()),
		}
	}

	return nil
}

func (w *Img) mount(parent base.Control) (base.Element, error) {
	// Create the bitmap
	pixbuf, buffer, err := imageToPixbuf(w.Image)
	if err != nil {
		return nil, err
	}

	handle, err := gtk.ImageNewFromPixbuf(pixbuf)
	if err != nil {
		return nil, err
	}
	parent.Handle.Add(handle)
	handle.Show()

	retval := &imgElement{Control{&handle.Widget}, buffer, w.Width, w.Height}
	handle.Connect("destroy", imgOnDestroy, retval)

	return retval, nil
}

func imgOnDestroy(widget *gtk.Label, mounted *imgElement) {
	mounted.handle = nil
}

func (w *imgElement) image() *gtk.Image {
	return (*gtk.Image)(unsafe.Pointer(w.handle))
}

func (w *imgElement) Props() base.Widget {
	return &Img{
		Image:  pixbufToImage(w.image().GetPixbuf()),
		Width:  w.width,
		Height: w.height,
	}
}

func (w *imgElement) updateProps(data *Img) error {
	w.width, w.height = data.Width, data.Height

	// Create the bitmap
	pixbuf, buffer, err := imageToPixbuf(data.Image)
	if err != nil {
		return err
	}
	w.imageData = buffer
	w.image().SetFromPixbuf(pixbuf)

	return nil
}
