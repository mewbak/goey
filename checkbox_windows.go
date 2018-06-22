package goey

import (
	"github.com/lxn/win"
	"syscall"
	"unsafe"
)

func (w *Checkbox) mount(parent Control) (Element, error) {
	text, err := syscall.UTF16FromString(w.Text)
	if err != nil {
		return nil, err
	}

	hwnd := win.CreateWindowEx(0, button.className, &text[0],
		win.WS_CHILD|win.WS_VISIBLE|win.WS_TABSTOP|win.BS_CHECKBOX|win.BS_TEXT|win.BS_NOTIFY,
		10, 10, 100, 100,
		parent.hWnd, win.HMENU(nextControlID()), 0, nil)
	if hwnd == 0 {
		err := syscall.GetLastError()
		if err == nil {
			return nil, syscall.EINVAL
		}
		return nil, err
	}
	if w.Value {
		win.SendMessage(hwnd, win.BM_SETCHECK, win.BST_CHECKED, 0)
	}

	// Set the font for the window
	if hMessageFont != 0 {
		win.SendMessage(hwnd, win.WM_SETFONT, uintptr(hMessageFont), 0)
	}

	if w.Disabled {
		win.EnableWindow(hwnd, false)
	}

	// Subclass the window procedure
	subclassWindowProcedure(hwnd, &button.oldWindowProc, syscall.NewCallback(checkboxWindowProc))

	retval := &mountedCheckbox{
		Control:  Control{hwnd},
		text:     text,
		onChange: w.OnChange,
		onFocus:  w.OnFocus,
		onBlur:   w.OnBlur,
	}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

type mountedCheckbox struct {
	Control
	text     []uint16
	onChange func(value bool)
	onFocus  func()
	onBlur   func()
}

func (w *mountedCheckbox) Props() Widget {
	return &Checkbox{
		Text:     w.Control.Text(),
		Value:    win.SendMessage(w.hWnd, win.BM_GETCHECK, 0, 0) == win.BST_CHECKED,
		Disabled: !win.IsWindowEnabled(w.hWnd),
		OnChange: w.onChange,
		OnFocus:  w.onFocus,
		OnBlur:   w.onBlur,
	}
}

func (w *mountedCheckbox) preferredWidth() Length {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing

	hdc := win.GetDC(w.hWnd)
	if hMessageFont != 0 {
		win.SelectObject(hdc, win.HGDIOBJ(hMessageFont))
	}
	rect := win.RECT{0, 0, 0xffff, 0xffff}
	win.DrawTextEx(hdc, &w.text[0], int32(len(w.text)), &rect, win.DT_CALCRECT, nil)
	win.ReleaseDC(w.hWnd, hdc)

	return FromPixelsX(int(rect.Right)) + 17
}

func (w *mountedCheckbox) Layout(bc Box) Size {
	// Determine ideal width.
	width := w.preferredWidth()
	height := 17 * DIP
	return bc.Constrain(Size{width, height})
}

func (w *mountedCheckbox) MinimumSize() Size {
	// Determine ideal width.
	width := w.preferredWidth()
	height := 17 * DIP
	return Size{width, height}
}

func (w *mountedCheckbox) updateProps(data *Checkbox) error {
	w.SetText(data.Text)
	w.SetDisabled(data.Disabled)
	if data.Value {
		win.SendMessage(w.hWnd, win.BM_SETCHECK, win.BST_CHECKED, 0)
	} else {
		win.SendMessage(w.hWnd, win.BM_SETCHECK, win.BST_UNCHECKED, 0)
	}

	w.onChange = data.OnChange
	w.onFocus = data.OnFocus
	w.onBlur = data.OnBlur

	return nil
}

func checkboxWindowProc(hwnd win.HWND, msg uint32, wParam uintptr, lParam uintptr) (result uintptr) {
	switch msg {
	case win.WM_DESTROY:
		// Make sure that the data structure on the Go-side does not point to a non-existent
		// window.
		if w := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA); w != 0 {
			ptr := (*mountedCheckbox)(unsafe.Pointer(w))
			ptr.hWnd = 0
		}
		// Defer to the old window proc

	case win.WM_SETFOCUS:
		if w := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA); w != 0 {
			ptr := (*mountedCheckbox)(unsafe.Pointer(w))
			if ptr.onFocus != nil {
				ptr.onFocus()
			}
		}
		// Defer to the old window proc

	case win.WM_KILLFOCUS:
		if w := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA); w != 0 {
			ptr := (*mountedCheckbox)(unsafe.Pointer(w))
			if ptr.onBlur != nil {
				ptr.onBlur()
			}
		}
		// Defer to the old window proc

	case win.WM_COMMAND:
		notification := win.HIWORD(uint32(wParam))
		switch notification {
		case win.BN_CLICKED:
			check := uintptr(win.BST_CHECKED)
			if win.SendMessage(hwnd, win.BM_GETCHECK, 0, 0) == win.BST_CHECKED {
				check = win.BST_UNCHECKED
			}
			win.SendMessage(hwnd, win.BM_SETCHECK, check, 0)
			if w := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA); w != 0 {
				ptr := (*mountedCheckbox)(unsafe.Pointer(w))
				if ptr.onChange != nil {
					ptr.onChange(check == win.BST_CHECKED)
				}
			}
		}
		return 0
	}

	return win.CallWindowProc(button.oldWindowProc, hwnd, msg, wParam, lParam)
}
