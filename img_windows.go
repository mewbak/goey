package goey

import (
	"github.com/lxn/win"
	"image"
	"syscall"
	"unsafe"
)

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
	} else {
		panic("Unsupported image format.")
	}
}

func (w *Img) mount(parent NativeWidget) (MountedWidget, error) {
	// Create the bitmap
	hbitmap, buffer, err := imageToBitmap(w.Image)
	if err != nil {
		return nil, err
	}

	hwnd := win.CreateWindowEx(0, staticClassName, nil,
		win.WS_CHILD|win.WS_VISIBLE|win.SS_BITMAP|win.SS_LEFT,
		10, 10, 100, 100,
		parent.hWnd, 0, 0, nil)
	if hwnd == 0 {
		err := syscall.GetLastError()
		if err == nil {
			return nil, syscall.EINVAL
		}
		return nil, err
	}
	win.SendMessage(hwnd, STM_SETIMAGE, win.IMAGE_BITMAP, uintptr(hbitmap))

	retval := &mountedImg{NativeWidget: NativeWidget{hwnd}, imageData: buffer, width: w.Width, height: w.Height}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

type mountedImg struct {
	NativeWidget
	imageData     []uint8
	width, height DIP
}

func (w *mountedImg) MeasureWidth() (DIP, DIP) {
	return w.width, w.width
}

func (w *mountedImg) MeasureHeight(width DIP) (DIP, DIP) {
	return w.height, w.height
}

func (w *mountedImg) SetBounds(bounds image.Rectangle) {
	w.NativeWidget.SetBounds(bounds)

	// Not certain why this is required.  However, static controls don't
	// repaint when resized.  This forces a repaint.
	win.InvalidateRect(w.hWnd, nil, true)
}

func (w *mountedImg) updateProps(data *Img) error {
	w.width, w.height = data.Width, data.Height

	// Create the bitmap
	hbitmap, buffer, err := imageToBitmap(data.Image)
	if err != nil {
		return err
	}
	w.imageData = buffer
	win.SendMessage(w.hWnd, STM_SETIMAGE, win.IMAGE_BITMAP, uintptr(hbitmap))

	return nil
}
