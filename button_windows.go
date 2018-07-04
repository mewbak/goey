package goey

import (
	"github.com/lxn/win"
	"syscall"
	"unsafe"
)

var (
	button struct {
		className     *uint16
		oldWindowProc uintptr
	}
)

func init() {
	var err error
	button.className, err = syscall.UTF16PtrFromString("BUTTON")
	if err != nil {
		panic(err)
	}
}

func (w *Button) mount(parent Control) (Element, error) {
	text, err := syscall.UTF16FromString(w.Text)
	if err != nil {
		return nil, err
	}

	style := uint32(win.WS_CHILD | win.WS_VISIBLE | win.WS_TABSTOP | win.BS_PUSHBUTTON | win.BS_TEXT | win.BS_NOTIFY)
	if w.Default {
		style = style | win.BS_DEFPUSHBUTTON
	}

	hwnd := win.CreateWindowEx(0, button.className, &text[0], style,
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
	subclassWindowProcedure(hwnd, &button.oldWindowProc, syscall.NewCallback(buttonWindowProc))

	retval := &mountedButton{
		Control: Control{hwnd},
		text:    text,
		onClick: w.OnClick,
		onFocus: w.OnFocus,
		onBlur:  w.OnBlur,
	}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

type mountedButton struct {
	Control
	text []uint16

	onClick func()
	onFocus func()
	onBlur  func()
}

func (w *mountedButton) Props() Widget {
	return &Button{
		Text:     w.Control.Text(),
		Disabled: !win.IsWindowEnabled(w.hWnd),
		Default:  (win.GetWindowLong(w.hWnd, win.GWL_STYLE) & win.BS_DEFPUSHBUTTON) != 0,
		OnClick:  w.onClick,
		OnFocus:  w.onFocus,
		OnBlur:   w.onBlur,
	}
}

func (w *mountedButton) Layout(bc Constraint) Size {
	width := w.MinIntrinsicWidth(0)
	height := w.MinIntrinsicHeight(0)
	return bc.Constrain(Size{width, height})
}

func (w *mountedButton) MinIntrinsicHeight(Length) Length {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	return 23 * DIP
}

func (w *mountedButton) MinIntrinsicWidth(Length) Length {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	width, _ := w.CalcRect(w.text)
	return max(
		75*DIP,
		FromPixelsX(int(width)+7),
	)
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
		buttonGetPtr(hwnd).hWnd = 0
		// Defer to the old window proc

	case win.WM_SETFOCUS:
		if w := buttonGetPtr(hwnd); w.onFocus != nil {
			w.onFocus()
		}
		// Defer to the old window proc

	case win.WM_KILLFOCUS:
		if w := buttonGetPtr(hwnd); w.onBlur != nil {
			w.onBlur()
		}
		// Defer to the old window proc

	case win.WM_COMMAND:
		notification := win.HIWORD(uint32(wParam))
		switch notification {
		case win.BN_CLICKED:
			if w := buttonGetPtr(hwnd); w.onClick != nil {
				w.onClick()
			}
		}
		return 0
	}

	return win.CallWindowProc(button.oldWindowProc, hwnd, msg, wParam, lParam)
}

func buttonGetPtr(hwnd win.HWND) *mountedButton {
	gwl := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA)
	if gwl == 0 {
		panic("Internal error.")
	}

	ptr := (*mountedButton)(unsafe.Pointer(gwl))
	if ptr.hWnd != hwnd && ptr.hWnd != 0 {
		panic("Internal error.")
	}

	return ptr
}
