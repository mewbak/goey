package cocoa

/*
#include "cocoa.h"
#include <stdlib.h>
*/
import "C"
import "image"
import "image/draw"
import "unsafe"

// ImageView is a wrapper for a NSImageView.
type ImageView struct {
	Control
	private int
}

func imageToNSImage(prop image.Image) (unsafe.Pointer, error) {
	if img, ok := prop.(*image.RGBA); ok {
		// Create the NSImage
		handle := C.imageNewFromRGBA(
			(*C.uint8_t)(unsafe.Pointer(&img.Pix[0])),
			C.int(img.Rect.Dx()), C.int(img.Rect.Dy()), C.int(img.Stride))
		return handle, nil
	} else if img, ok := prop.(*image.Gray); ok {
		// Create the NSImage
		handle := C.imageNewFromGray(
			(*C.uint8_t)(unsafe.Pointer(&img.Pix[0])),
			C.int(img.Rect.Dx()), C.int(img.Rect.Dy()), C.int(img.Stride))
		return handle, nil
	}

	// Create a new image in RGBA format
	bounds := prop.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, prop, bounds.Min, draw.Src)

	// Create the NSImage
	handle := C.imageNewFromRGBA(
		(*C.uint8_t)(unsafe.Pointer(&img.Pix[0])),
		C.int(img.Rect.Dx()), C.int(img.Rect.Dy()), C.int(img.Stride))
	return handle, nil
}

func NewImageView(window *View, prop image.Image) (*ImageView, error) {
	image, err := imageToNSImage(prop)
	if err != nil {
		return nil, err
	}

	control := C.imageviewNew(unsafe.Pointer(window), image)
	C.imageClose(image)
	return (*ImageView)(control), nil
}

func (w *ImageView) Image() image.Image {
	return nil
}

func (w *ImageView) SetImage(prop image.Image) error {
	image, err := imageToNSImage(prop)
	if err != nil {
		return err
	}

	C.imageviewSetImage(unsafe.Pointer(w), image)
	C.imageClose(image)
	return nil
}
