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
		handle := C.imageNewFromRGBA(unsafe.Pointer(&img.Pix[0]), C.int(img.Rect.Dx()), C.int(img.Rect.Dy()))
		return handle, nil
	} else if img, ok := prop.(*image.Gray); ok {
		handle := C.imageNewFromGray(unsafe.Pointer(&img.Pix[0]), C.int(img.Rect.Dx()), C.int(img.Rect.Dy()))
		return handle, nil
	}

	// Create a new image in RGBA format
	bounds := prop.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, prop, bounds.Min, draw.Src)
	// Need to convert RGB to BGR
	for i := 0; i < len(img.Pix); i += 4 {
		img.Pix[i+0], img.Pix[i+2] = img.Pix[i+2], img.Pix[i+0]
	}

	// Create the bitmap
	handle := C.imageNewFromRGBA(unsafe.Pointer(&img.Pix[0]), C.int(img.Rect.Dx()), C.int(img.Rect.Dy()))
	return handle, nil
}

func NewImageView(window *Window, prop image.Image) (*ImageView, error) {
	image, err := imageToNSImage(prop)
	if err != nil {
		return nil, err
	}

	control := C.imageviewNew(unsafe.Pointer(window), image)
	return (*ImageView)(control), nil
}

func (w *ImageView) Close() {
	C.controlClose(unsafe.Pointer(w))
}
