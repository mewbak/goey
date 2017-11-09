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

func (w *Button) Mount(parent NativeWidget) (MountedWidget, error) {
	text, err := syscall.UTF16PtrFromString(w.Text)
	if err != nil {
		return nil, err
	}

	style := uint32(win.WS_CHILD | win.WS_VISIBLE | win.WS_TABSTOP | win.BS_PUSHBUTTON | win.BS_TEXT | win.BS_NOTIFY)
	if w.Default {
		style = style | win.BS_DEFPUSHBUTTON
	}

	hwnd := win.CreateWindowEx(0, buttonClassName, text, style,
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
	if oldButtonWindowProc == 0 {
		oldButtonWindowProc = win.GetWindowLongPtr(hwnd, win.GWLP_WNDPROC)
	} else {
		oldWindowProc := win.GetWindowLongPtr(hwnd, win.GWLP_WNDPROC)
		if oldWindowProc != oldButtonWindowProc {
			panic("Corrupted data")
		}
	}
	win.SetWindowLongPtr(hwnd, win.GWLP_WNDPROC, syscall.NewCallback(buttonWindowProc))

	retval := &MountedButton{
		NativeWidget: NativeWidget{hwnd},
		onClick:      w.OnClick,
		onFocus:      w.OnFocus,
		onBlur:       w.OnBlur,
	}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

type MountedButton struct {
	NativeWidget
	onClick func()
	onFocus func()
	onBlur  func()
}

func (w *MountedButton) PreferredWidth() int {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	// In the future, we should calculate the width based on the length of the text.
	return 50
}

func (w *MountedButton) CalculateHeight(width int) int {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	return 23
}

func (w *MountedButton) UpdateProps(data_ Widget) error {
	data := data_.(*Button)

	w.SetText(data.Text)
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
			ptr := (*MountedButton)(unsafe.Pointer(w))
			ptr.hWnd = 0
		}
		// Defer to the old window proc

	case win.WM_SETFOCUS:
		if w := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA); w != 0 {
			ptr := (*MountedButton)(unsafe.Pointer(w))
			if ptr.onFocus != nil {
				ptr.onFocus()
			}
		}
		// Defer to the old window proc

	case win.WM_KILLFOCUS:
		if w := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA); w != 0 {
			ptr := (*MountedButton)(unsafe.Pointer(w))
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
				ptr := (*MountedButton)(unsafe.Pointer(w))
				if ptr.onClick != nil {
					ptr.onClick()
				}
			}
		}
		return 0
	}

	return win.CallWindowProc(oldButtonWindowProc, hwnd, msg, wParam, lParam)
}
