package goey

import (
	"github.com/lxn/win"
	"syscall"
	"unsafe"
)

func (w *Checkbox) Mount(parent NativeWidget) (MountedWidget, error) {
	text, err := syscall.UTF16PtrFromString(w.Text)
	if err != nil {
		return nil, err
	}

	hwnd := win.CreateWindowEx(0, buttonClassName, text,
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
	if oldButtonWindowProc == 0 {
		oldButtonWindowProc = win.GetWindowLongPtr(hwnd, win.GWLP_WNDPROC)
	} else {
		oldWindowProc := win.GetWindowLongPtr(hwnd, win.GWLP_WNDPROC)
		if oldWindowProc != oldButtonWindowProc {
			panic("Corrupted data")
		}
	}
	win.SetWindowLongPtr(hwnd, win.GWLP_WNDPROC, syscall.NewCallback(checkboxWindowProc))

	retval := &MountedCheckbox{NativeWidget: NativeWidget{hwnd}, onChange: w.OnChange}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

type MountedCheckbox struct {
	NativeWidget
	onChange func(value bool)
	onFocus  func()
	onBlur   func()
}

func (w *MountedCheckbox) MinimumWidth() DP {
	// In the future, we should calculate the width based on the length of the text.
	return 160
}

func (w *MountedCheckbox) CalculateHeight(width DP) DP {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	return 17
}

func (w *MountedCheckbox) UpdateProps(data_ Widget) error {
	data := data_.(*Checkbox)

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
			ptr := (*MountedCheckbox)(unsafe.Pointer(w))
			ptr.hWnd = 0
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
				ptr := (*MountedCheckbox)(unsafe.Pointer(w))
				if ptr.onChange != nil {
					ptr.onChange(check == win.BST_CHECKED)
				}
			}
		}
		return 0
	}

	return win.CallWindowProc(oldButtonWindowProc, hwnd, msg, wParam, lParam)
}
