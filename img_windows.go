package goey

import (
	"image"
	"image/draw"
	"unsafe"

	"bitbucket.org/rj/goey/base"
	win2 "bitbucket.org/rj/goey/syscall"
	"github.com/lxn/win"
)

func imageToIcon(prop image.Image) (win.HICON, []uint8, error) {
	// Create a mask for the icon.
	// Currently, we are using a straight white mask, but perhaps this
	// should be a copy of the alpha channel if the source image is
	// RGBA.
	bounds := prop.Bounds()
	imgMask := image.NewGray(prop.Bounds())
	draw.Draw(imgMask, bounds, image.White, image.Point{}, draw.Src)
	hmask, _, err := imageToBitmap(imgMask)
	if err != nil {
		return 0, nil, err
	}

	// Convert the image to a bitmap.
	hbitmap, buffer, err := imageToBitmap(prop)
	if err != nil {
		return 0, nil, err
	}

	// Create the icon
	iconinfo := win.ICONINFO{
		FIcon:    win.TRUE,
		HbmMask:  hmask,
		HbmColor: hbitmap,
	}
	hicon := win.CreateIconIndirect(&iconinfo)
	if hicon == 0 {
		panic("Error in CreateIconIndirect")
	}
	return hicon, buffer, nil
}

func imageToBitmap(prop image.Image) (win.HBITMAP, []uint8, error) {
	if img, ok := prop.(*image.RGBA); ok {
		// Create a copy of the backing for the pixel data
		buffer := append([]uint8(nil), img.Pix...)
		// Need to convert RGB to BGR
		for i := 0; i < len(buffer); i += 4 {
			buffer[i+0], buffer[i+2] = buffer[i+2], buffer[i+0]
		}

		// Create the bitmap
		hbitmap := win.CreateBitmap(int32(img.Rect.Dx()), int32(img.Rect.Dy()), 4, 8, unsafe.Pointer(&buffer[0]))
		if hbitmap == 0 {
			panic("Error in CreateBitmap")
		}
		return hbitmap, buffer, nil
	} else if img, ok := prop.(*image.Gray); ok {
		// Create a copy of the backing for the pixel data
		buffer := append([]uint8(nil), img.Pix...)
		// Create the bitmap
		hbitmap := win.CreateBitmap(int32(img.Rect.Dx()), int32(img.Rect.Dy()), 1, 8, unsafe.Pointer(&img.Pix[0]))
		if hbitmap == 0 {
			panic("Error in CreateBitmap")
		}
		return hbitmap, buffer, nil
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
	hbitmap := win.CreateBitmap(int32(img.Rect.Dx()), int32(img.Rect.Dy()), 4, 8, unsafe.Pointer(&img.Pix[0]))
	if hbitmap == 0 {
		panic("Error in CreateBitmap")
	}
	return hbitmap, img.Pix, nil
}

func bitmapToImage(hdc win.HDC, hbitmap win.HBITMAP) image.Image {
	bmi := win.BITMAPINFO{}
	bmi.BmiHeader.BiSize = uint32(unsafe.Sizeof(bmi))
	win.GetDIBits(hdc, hbitmap, 0, 0, nil, &bmi, 0)
	if bmi.BmiHeader.BiPlanes == 1 && bmi.BmiHeader.BiBitCount == 32 && bmi.BmiHeader.BiCompression == 3 /*BI_BITFIELDS*/ {
		// Get the pixel data
		buffer := make([]byte, bmi.BmiHeader.BiSizeImage)
		win.GetDIBits(hdc, hbitmap, 0, uint32(bmi.BmiHeader.BiHeight), &buffer[0], &bmi, 0)

		// Need to convert BGR to RGB
		for i := 0; i < len(buffer); i += 4 {
			buffer[i+0], buffer[i+2] = buffer[i+2], buffer[i+0]
		}
		return &image.RGBA{
			Pix:    buffer,
			Stride: int(bmi.BmiHeader.BiWidth * 4),
			Rect:   image.Rect(0, 0, int(bmi.BmiHeader.BiWidth), int(bmi.BmiHeader.BiHeight)),
		}
	}

	return nil
}

func (w *Img) mount(parent base.Control) (base.Element, error) {
	// Create the bitmap
	hbitmap, buffer, err := imageToBitmap(w.Image)
	if err != nil {
		return nil, err
	}

	// Create the control
	const STYLE = win.WS_CHILD | win.WS_VISIBLE | win.SS_BITMAP | win.SS_LEFT
	hwnd, _, err := createControlWindow(0, &staticClassName[0], "", STYLE, parent.HWnd)
	win.SendMessage(hwnd, win2.STM_SETIMAGE, win.IMAGE_BITMAP, uintptr(hbitmap))

	retval := &imgElement{
		Control:   Control{hwnd},
		imageData: buffer,
		width:     w.Width,
		height:    w.Height,
	}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

type imgElement struct {
	Control
	imageData []uint8
	width     base.Length
	height    base.Length
}

func (w *imgElement) Props() base.Widget {
	// Need to recreate the image from the HBITMAP
	hbitmap := win.HBITMAP(win.SendMessage(w.hWnd, win2.STM_GETIMAGE, 0 /*IMAGE_BITMAP*/, 0))
	if hbitmap == 0 {
		return &Img{
			Width:  w.width,
			Height: w.height,
		}
	}

	hdc := win.GetDC(w.hWnd)
	img := bitmapToImage(hdc, hbitmap)
	win.ReleaseDC(w.hWnd, hdc)

	return &Img{
		Image:  img,
		Width:  w.width,
		Height: w.height,
	}
}

func (w *imgElement) SetBounds(bounds base.Rectangle) {
	w.Control.SetBounds(bounds)

	// Not certain why this is required.  However, static controls don't
	// repaint when resized.  This forces a repaint.
	win.InvalidateRect(w.hWnd, nil, true)
}

func (w *imgElement) updateProps(data *Img) error {
	w.width, w.height = data.Width, data.Height

	// Create the bitmap
	hbitmap, buffer, err := imageToBitmap(data.Image)
	if err != nil {
		return err
	}
	w.imageData = buffer
	win.SendMessage(w.hWnd, win2.STM_SETIMAGE, win.IMAGE_BITMAP, uintptr(hbitmap))

	return nil
}
