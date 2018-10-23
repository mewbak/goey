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

func imageToNSImage(prop image.Image) (unsafe.Pointer, []byte, error) {
	if img, ok := prop.(*image.RGBA); ok {
		// Create a copy of the backing for the pixel data
		buffer := append([]uint8(nil), img.Pix...)
		// Create the NSImage
		handle := C.imageNewFromRGBA(
			(*C.uint8_t)(unsafe.Pointer(&buffer[0])),
			C.int(img.Rect.Dx()), C.int(img.Rect.Dy()))
		return handle, buffer, nil
	} else if img, ok := prop.(*image.Gray); ok {
		// Create a copy of the backing for the pixel data
		buffer := append([]uint8(nil), img.Pix...)
		// Create the NSImage
		handle := C.imageNewFromGray(
			(*C.uint8_t)(unsafe.Pointer(&buffer[0])),
			C.int(img.Rect.Dx()), C.int(img.Rect.Dy()))
		return handle, buffer, nil
	}

	// Create a new image in RGBA format
	bounds := prop.Bounds()
	img := image.NewRGBA(bounds)
	draw.Draw(img, bounds, prop, bounds.Min, draw.Src)

	// Create the bitmap
	// Create a copy of the backing for the pixel data
	buffer := append([]uint8(nil), img.Pix...)
	// Create the NSImage
	handle := C.imageNewFromRGBA(
		(*C.uint8_t)(unsafe.Pointer(&buffer[0])),
		C.int(img.Rect.Dx()), C.int(img.Rect.Dy()))
	return handle, buffer, nil
}

func NewImageView(window *Window, prop image.Image) (*ImageView, []byte, error) {
	image, buffer, err := imageToNSImage(prop)
	if err != nil {
		return nil, nil, err
	}

	control := C.imageviewNew(unsafe.Pointer(window), image)
	C.imageClose(image)
	return (*ImageView)(control), buffer, nil
}

func (w *ImageView) Close() {
	C.controlClose(unsafe.Pointer(w))
}

func (w *ImageView) SetImage(prop image.Image) ([]byte, error) {
	image, buffer, err := imageToNSImage(prop)
	if err != nil {
		return nil, err
	}

	C.imageviewSetImage(unsafe.Pointer(w), image)
	return buffer, nil
}
