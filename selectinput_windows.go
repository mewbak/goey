package goey

import (
	"github.com/lxn/win"
	"syscall"
	"unsafe"
)

var (
	comboboxClassName     *uint16
	oldComboboxWindowProc uintptr
)

func init() {
	var err error
	comboboxClassName, err = syscall.UTF16PtrFromString("COMBOBOX")
	if err != nil {
		panic(err)
	}
}

func (w *SelectInput) mount(parent NativeWidget) (MountedWidget, error) {
	if w.Value >= len(w.Items) {
		w.Value = len(w.Items) - 1
	}
	hwnd := win.CreateWindowEx(win.WS_EX_CLIENTEDGE, comboboxClassName, nil,
		win.WS_CHILD|win.WS_VISIBLE|win.WS_TABSTOP|win.CBS_DROPDOWNLIST,
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

	// Add items to the control
	for _, v := range w.Items {
		text, err := syscall.UTF16PtrFromString(v)
		if err != nil {
			win.DestroyWindow(hwnd)
			return nil, err
		}
		win.SendMessage(hwnd, win.CB_ADDSTRING, 0, uintptr(unsafe.Pointer(text)))
	}
	if !w.Unset {
		win.SendMessage(hwnd, win.CB_SETCURSEL, uintptr(w.Value), 0)
	}

	// Subclass the window procedure
	if oldComboboxWindowProc == 0 {
		oldComboboxWindowProc = win.GetWindowLongPtr(hwnd, win.GWLP_WNDPROC)
	} else {
		oldWindowProc := win.GetWindowLongPtr(hwnd, win.GWLP_WNDPROC)
		if oldWindowProc != oldComboboxWindowProc {
			panic("Corrupted data")
		}
	}
	win.SetWindowLongPtr(hwnd, win.GWLP_WNDPROC, syscall.NewCallback(comboboxWindowProc))

	retval := &mountedSelectInput{
		NativeWidget: NativeWidget{hwnd},
		onChange:     w.OnChange,
		onFocus:      w.OnFocus,
		onBlur:       w.OnBlur,
	}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

type mountedSelectInput struct {
	NativeWidget
	onChange func(value int)
	onFocus  func()
	onBlur   func()
}

func (w *mountedSelectInput) MeasureWidth() (DIP, DIP) {
	// In the future, we should calculate the width based on the length of the text.
	return 160, 160
}

func (w *mountedSelectInput) MeasureHeight(width DIP) (DIP, DIP) {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	return 23, 23
}

func (w *mountedSelectInput) updateProps(data *SelectInput) error {
	// TODO:  Update the items in the combobox
	// TODO:  Update the selection based on Value
	// TODO:  Update the selection based on Unset.

	w.SetDisabled(data.Disabled)
	// TODO:  Update property .Default
	w.onChange = data.OnChange
	w.onFocus = data.OnFocus
	w.onBlur = data.OnBlur

	return nil
}

func comboboxWindowProc(hwnd win.HWND, msg uint32, wParam uintptr, lParam uintptr) (result uintptr) {
	switch msg {
	case win.WM_DESTROY:
		// Make sure that the data structure on the Go-side does not point to a non-existent
		// window.
		if w := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA); w != 0 {
			ptr := (*mountedSelectInput)(unsafe.Pointer(w))
			ptr.hWnd = 0
		}
		// Defer to the old window proc

	case win.WM_SETFOCUS:
		if w := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA); w != 0 {
			ptr := (*mountedSelectInput)(unsafe.Pointer(w))
			if ptr.onFocus != nil {
				ptr.onFocus()
			}
		}
		// Defer to the old window proc

	case win.WM_KILLFOCUS:
		if w := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA); w != 0 {
			ptr := (*mountedSelectInput)(unsafe.Pointer(w))
			if ptr.onBlur != nil {
				ptr.onBlur()
			}
		}
		// Defer to the old window proc

	case win.WM_COMMAND:
		notification := win.HIWORD(uint32(wParam))
		switch notification {
		case win.CBN_SELCHANGE:
			if w := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA); w != 0 {
				ptr := (*mountedSelectInput)(unsafe.Pointer(w))
				if ptr.onChange != nil {
					cursel := win.SendMessage(hwnd, win.CB_GETCURSEL, 0, 0)
					ptr.onChange(int(cursel))
				}
			}
		}
		// defer to old window proc
	}

	return win.CallWindowProc(oldComboboxWindowProc, hwnd, msg, wParam, lParam)
}
