package goey

import (
	"github.com/lxn/win"
	"syscall"
	"unsafe"
)

var (
	buttonClassName     *uint16
	oldButtonWindowProc uintptr
)

func init() {
	var err error
	buttonClassName, err = syscall.UTF16PtrFromString("BUTTON")
	if err != nil {
		panic(err)
	}
}

func (w *Button) mount(parent NativeWidget) (MountedWidget, error) {
	text, err := syscall.UTF16FromString(w.Text)
	if err != nil {
		return nil, err
	}

	style := uint32(win.WS_CHILD | win.WS_VISIBLE | win.WS_TABSTOP | win.BS_PUSHBUTTON | win.BS_TEXT | win.BS_NOTIFY)
	if w.Default {
		style = style | win.BS_DEFPUSHBUTTON
	}

	hwnd := win.CreateWindowEx(0, buttonClassName, &text[0], style,
		10, 10, 100, 100,
		parent.hWnd, win.HMENU(nextControlID()), 0, nil)
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

	if w.Disabled {
		win.EnableWindow(hwnd, false)
	}

	// Subclass the window procedure
	subclassWindowProcedure(hwnd, &oldButtonWindowProc, syscall.NewCallback(buttonWindowProc))

	retval := &mountedButton{
		NativeWidget: NativeWidget{hwnd},
		text:         text,
		onClick:      w.OnClick,
		onFocus:      w.OnFocus,
		onBlur:       w.OnBlur,
	}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

type mountedButton struct {
	NativeWidget
	text []uint16

	onClick func()
	onFocus func()
	onBlur  func()
}

func (w *mountedButton) MeasureWidth() (DIP, DIP) {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing

	hdc := win.GetDC(w.hWnd)
	if hMessageFont != 0 {
		win.SelectObject(hdc, win.HGDIOBJ(hMessageFont))
	}
	rect := win.RECT{0, 0, 0xffff, 0xffff}
	win.DrawTextEx(hdc, &w.text[0], int32(len(w.text)), &rect, win.DT_CALCRECT, nil)
	win.ReleaseDC(w.hWnd, hdc)

	retval := ToDIPX(int(rect.Right) + 7)
	if retval < 75 {
		return 75, 75
	}

	return retval, retval
}

func (w *mountedButton) MeasureHeight(width DIP) (DIP, DIP) {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	return 23, 23
}

func (w *mountedButton) updateProps(data *Button) error {
	text, err := syscall.UTF16FromString(data.Text)
	if err != nil {
		return err
	}

	w.SetText(data.Text)
	w.text = text
	w.SetDisabled(data.Disabled)
	// TODO:  Update property .Default
	w.onClick = data.OnClick
	w.onFocus = data.OnFocus
	w.onBlur = data.OnBlur

	return nil
}

func buttonWindowProc(hwnd win.HWND, msg uint32, wParam uintptr, lParam uintptr) (result uintptr) {
	switch msg {
	case win.WM_DESTROY:
		// Make sure that the data structure on the Go-side does not point to a non-existent
		// window.
		if w := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA); w != 0 {
			ptr := (*mountedButton)(unsafe.Pointer(w))
			ptr.hWnd = 0
		}
		// Defer to the old window proc

	case win.WM_SETFOCUS:
		if w := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA); w != 0 {
			ptr := (*mountedButton)(unsafe.Pointer(w))
			if ptr.onFocus != nil {
				ptr.onFocus()
			}
		}
		// Defer to the old window proc

	case win.WM_KILLFOCUS:
		if w := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA); w != 0 {
			ptr := (*mountedButton)(unsafe.Pointer(w))
			if ptr.onBlur != nil {
				ptr.onBlur()
			}
		}
		// Defer to the old window proc

	case win.WM_COMMAND:
		notification := win.HIWORD(uint32(wParam))
		switch notification {
		case win.BN_CLICKED:
			if w := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA); w != 0 {
				ptr := (*mountedButton)(unsafe.Pointer(w))
				if ptr.onClick != nil {
					ptr.onClick()
				}
			}
		}
		return 0
	}

	return win.CallWindowProc(oldButtonWindowProc, hwnd, msg, wParam, lParam)
}
