package goey

import (
	"github.com/lxn/win"
	"image"
	"syscall"
	"unsafe"
)

func imageToBitmap(prop image.Image) (win.HBITMAP, error) {
	if img, ok := prop.(*image.RGBA); ok {
		hbitmap := win.CreateBitmap(int32(img.Rect.Dx()), int32(img.Rect.Dy()), 4, 8, unsafe.Pointer(&img.Pix[0]))
		if hbitmap == 0 {
			panic("Error in CreateBitmap")
		}
		return hbitmap, nil
	} else {
		panic("Unsupported image format.")
	}
}

func (w *Img) mount(parent NativeWidget) (MountedWidget, error) {
	// Create the bitmap
	hbitmap, err := imageToBitmap(w.Image)
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

	retval := &mountedImg{NativeWidget: NativeWidget{hwnd}, width: w.Width, height: w.Height}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

type mountedImg struct {
	NativeWidget
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
	hbitmap, err := imageToBitmap(data.Image)
	if err != nil {
		return err
	}
	win.SendMessage(w.hWnd, STM_SETIMAGE, win.IMAGE_BITMAP, uintptr(hbitmap))

	return nil
}
