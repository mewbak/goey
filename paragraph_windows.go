package goey

import (
	win2 "bitbucket.org/rj/goey/syscall"
	"github.com/lxn/win"
	"image"
	"syscall"
	"unsafe"
)

var (
	paragraphMinWidth int
	paragraphMaxWidth int
)

func (w *P) calcStyle() uint32 {
	style := uint32(win.WS_CHILD | win.WS_VISIBLE | win.SS_LEFT)
	if w.Align == Center {
		style = style | win.SS_CENTER
	} else if w.Align == Right {
		style = style | win.SS_RIGHT
	}
	return style
}

func (w *P) mount(parent NativeWidget) (MountedWidget, error) {
	text, err := syscall.UTF16FromString(w.Text)
	if err != nil {
		return nil, err
	}

	hwnd := win.CreateWindowEx(0, staticClassName, &text[0],
		w.calcStyle(),
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

	retval := &mountedP{NativeWidget: NativeWidget{hwnd}, text: text}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

type mountedP struct {
	NativeWidget
	text []uint16
}

func paragraphMeasureReflowLimits(hwnd win.HWND) {
	hdc := win.GetDC(hwnd)
	if hMessageFont != 0 {
		win.SelectObject(hdc, win.HGDIOBJ(hMessageFont))
	}
	// Calculate the width of a single 'm' (find the em width)
	rect := win.RECT{0, 0, 0xffff, 0xffff}
	caption := [...]uint16{'m'}
	win.DrawTextEx(hdc, &caption[0], 1, &rect, win.DT_CALCRECT, nil)
	win.ReleaseDC(hwnd, hdc)
	paragraphMinWidth = int(rect.Right) * 20
	paragraphMaxWidth = int(rect.Right) * 80
}

func (w *mountedP) MeasureWidth() (DIP, DIP) {
	// If the printed text will be more than 80 characters wide, it will start
	// to impact readability.  We want to force reflow in this case, so we limit
	// the width
	//
	// See the following for the conversion from characters to relative pixels.
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	if paragraphMaxWidth == 0 {
		paragraphMeasureReflowLimits(w.hWnd)
	}

	hdc := win.GetDC(w.hWnd)
	if hMessageFont != 0 {
		win.SelectObject(hdc, win.HGDIOBJ(hMessageFont))
	}
	rect := win.RECT{0, 0, 0xffff, 0xffff}
	win.DrawTextEx(hdc, &w.text[0], int32(len(w.text)), &rect, win.DT_CALCRECT|win.DT_WORDBREAK, nil)
	win.ReleaseDC(w.hWnd, hdc)

	// For reflow if the text is more than 60 characters wide
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	if int(rect.Right) > paragraphMinWidth {
		return ToDIPX(paragraphMinWidth), ToDIPX(paragraphMaxWidth)
	}

	// Not enough text for reflow.
	retval := ToDIPX(int(rect.Right))
	return retval, retval
}

func (w *mountedP) MeasureHeight(width DIP) (DIP, DIP) {
	hdc := win.GetDC(w.hWnd)
	if hMessageFont != 0 {
		win.SelectObject(hdc, win.HGDIOBJ(hMessageFont))
	}
	rect := win.RECT{0, 0, int32(width), 0xffff}
	win.DrawTextEx(hdc, &w.text[0], int32(len(w.text)), &rect, win.DT_CALCRECT|win.DT_WORDBREAK, nil)
	win.ReleaseDC(w.hWnd, hdc)

	retval := ToDIPY(int(rect.Bottom))
	return retval, retval
}

func (w *mountedP) SetBounds(bounds image.Rectangle) {
	w.NativeWidget.SetBounds(bounds)

	// Not certain why this is required.  However, static controls don't
	// repaint when resized.  This forces a repaint.
	win.InvalidateRect(w.hWnd, nil, true)
}

func (w *mountedP) updateProps(data *P) error {
	text, err := syscall.UTF16FromString(data.Text)
	if err != nil {
		return err
	}
	w.text = text
	win2.SetWindowText(w.hWnd, &text[0])
	win.SetWindowLongPtr(w.hWnd, win.GWL_STYLE, uintptr(data.calcStyle()))

	return nil
}
