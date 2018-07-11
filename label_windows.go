package goey

import (
	"syscall"
	"unsafe"

	win2 "bitbucket.org/rj/goey/syscall"
	"github.com/lxn/win"
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

func (w *Label) mount(parent Control) (Element, error) {
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

	retval := &labelElement{Control: Control{hwnd}, text: text}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

type labelElement struct {
	Control
	text []uint16
}

func (w *labelElement) Props() Widget {
	return &Label{
		Text: w.Control.Text(),
	}
}

func (w *labelElement) Layout(bc Constraint) Size {
	width := w.MinIntrinsicWidth(0)
	height := w.MinIntrinsicHeight(0)
	return bc.Constrain(Size{width, height})
}

func (w *labelElement) MinIntrinsicHeight(Length) Length {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	return 13 * DIP
}

func (w *labelElement) MinIntrinsicWidth(Length) Length {
	width, _ := w.CalcRect(w.text)
	return FromPixelsX(int(width))
}

func (w *labelElement) SetBounds(bounds Rectangle) {
	w.Control.SetBounds(bounds)

	// Not certain why this is required.  However, static controls don't
	// repaint when resized.  This forces a repaint.
	win.InvalidateRect(w.hWnd, nil, true)
}

func (w *labelElement) updateProps(data *Label) error {
	text, err := syscall.UTF16FromString(data.Text)
	if err != nil {
		return err
	}
	w.text = text
	win2.SetWindowText(w.hWnd, &text[0])
	// TODO:  Update alignment

	return nil
}
