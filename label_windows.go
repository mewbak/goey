package goey

import (
	"github.com/lxn/win"
	"image"
	"syscall"
	"unsafe"
)

var (
	staticClassName     *uint16
	oldStaticWindowProc uintptr
)

func init() {
	var err error
	staticClassName, err = syscall.UTF16PtrFromString("STATIC")
	if err != nil {
		panic(err)
	}
}

func (w *Label) Mount(parent NativeWidget) (MountedWidget, error) {
	text, err := syscall.UTF16FromString(w.Text)
	if err != nil {
		return nil, err
	}

	hwnd := win.CreateWindowEx(0, staticClassName, &text[0],
		win.WS_CHILD|win.WS_VISIBLE|win.SS_LEFT,
		10, 10, 100, 100,
		parent.hWnd, 0, 0, nil)
	if hwnd == 0 {
		err := syscall.GetLastError()
		if err == nil {
			return nil, syscall.EINVAL
		}
		return nil, err
	}

	// Set the font for the window
	if hMessageFont != 0 {
		win.SendMessage(hwnd, win.WM_SETFONT, uintptr(hMessageFont), 0)
	}

	retval := &MountedLabel{NativeWidget: NativeWidget{hwnd}, text: text}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

type MountedLabel struct {
	NativeWidget
	text []uint16
}

func (w *MountedLabel) MinimumWidth() DP {
	hdc := win.GetDC(w.hWnd)
	if hMessageFont != 0 {
		win.SelectObject(hdc, win.HGDIOBJ(hMessageFont))
	}
	rect := win.RECT{0, 0, 0xffff, 0xffff}
	win.DrawTextEx(hdc, &w.text[0], int32(len(w.text)), &rect, win.DT_CALCRECT, nil)
	win.ReleaseDC(w.hWnd, hdc)

	return DP(int(rect.Right) * 96 / dpi.X)
}

func (w *MountedLabel) CalculateHeight(width DP) DP {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	return 13
}

func (w *MountedLabel) SetBounds(bounds image.Rectangle) {
	w.NativeWidget.SetBounds(bounds)

	// Not certain why this is required.  However, static controls don't
	// repaint when resized.  This forces a repaint.
	win.InvalidateRect(w.hWnd, nil, true)
}

func (w *MountedLabel) UpdateProps(data_ Widget) error {
	data := data_.(*Label)

	text, err := syscall.UTF16FromString(data.Text)
	if err != nil {
		return err
	}
	w.text = text
	SetWindowText(w.hWnd, &text[0])
	// TODO:  Update alignment

	return nil
}
