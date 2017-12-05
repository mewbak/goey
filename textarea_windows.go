package goey

import (
	"github.com/lxn/win"
	"syscall"
	"unsafe"
)

func minlinesDefault(value int) int {
	if value < 1 {
		return 3
	}
	return value
}

func (w *TextArea) mount(parent NativeWidget) (MountedWidget, error) {
	text, err := syscall.UTF16PtrFromString(w.Value)
	if err != nil {
		return nil, err
	}

	style := uint32(win.WS_CHILD | win.WS_VISIBLE | win.WS_TABSTOP | win.ES_LEFT | win.ES_MULTILINE | win.ES_WANTRETURN | win.ES_AUTOVSCROLL)
	hwnd := win.CreateWindowEx(win.WS_EX_CLIENTEDGE, editClassName, text,
		style,
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
	if oldEditWindowProc == 0 {
		oldEditWindowProc = win.GetWindowLongPtr(hwnd, win.GWLP_WNDPROC)
	} else {
		oldWindowProc := win.GetWindowLongPtr(hwnd, win.GWLP_WNDPROC)
		if oldWindowProc != oldEditWindowProc {
			panic("Corrupted data")
		}
	}
	win.SetWindowLongPtr(hwnd, win.GWLP_WNDPROC, syscall.NewCallback(textareaWindowProc))

	// Create placeholder, if required.
	if w.Value == "" && w.Placeholder != "" {
		textPlaceholder, err := syscall.UTF16PtrFromString(w.Placeholder)
		if err != nil {
			win.DestroyWindow(hwnd)
			return nil, err
		}

		win.SendMessage(hwnd, win.EM_SETCUEBANNER, 0, uintptr(unsafe.Pointer(textPlaceholder)))
	}

	retval := &mountedTextArea{
		NativeWidget: NativeWidget{hwnd},
		minLines:     minlinesDefault(w.MinLines),
		onChange:     w.OnChange,
		onFocus:      w.OnFocus,
		onBlur:       w.OnBlur,
	}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

type mountedTextArea struct {
	NativeWidget
	minLines int
	onChange func(value string)
	onFocus  func()
	onBlur   func()
}

func (w *mountedTextArea) MeasureWidth() (DIP, DIP) {
	if paragraphMaxWidth == 0 {
		paragraphMeasureReflowLimits(w.hWnd)
	}

	// In the future, we should calculate the width based on the length of the text.
	return ToDIPX(paragraphMinWidth), ToDIPX(paragraphMaxWidth)
}

func (w *mountedTextArea) MeasureHeight(width DIP) (DIP, DIP) {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	return 23 + 16*DIP(w.minLines-1), 23 + 16*39
}

func (w *mountedTextArea) updateProps(data *TextArea) error {
	if data.Value != w.Text() {
		w.SetText(data.Value)
	}

	if data.Placeholder != "" {
		textPlaceholder, err := syscall.UTF16PtrFromString(data.Placeholder)
		if err != nil {
			return err
		}

		win.SendMessage(w.hWnd, win.EM_SETCUEBANNER, 0, uintptr(unsafe.Pointer(textPlaceholder)))
	} else {
		win.SendMessage(w.hWnd, win.EM_SETCUEBANNER, 0, 0)
	}

	w.minLines = minlinesDefault(data.MinLines)
	w.onChange = data.OnChange
	w.onFocus = data.OnFocus
	w.onBlur = data.OnBlur

	return nil
}

func textareaWindowProc(hwnd win.HWND, msg uint32, wParam uintptr, lParam uintptr) (result uintptr) {
	switch msg {
	case win.WM_DESTROY:
		// Make sure that the data structure on the Go-side does not point to a non-existent
		// window.
		if w := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA); w != 0 {
			ptr := (*mountedTextArea)(unsafe.Pointer(w))
			ptr.hWnd = 0
		}
		// Defer to the old window proc

	case win.WM_SETFOCUS:
		if w := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA); w != 0 {
			ptr := (*mountedTextArea)(unsafe.Pointer(w))
			if ptr.onFocus != nil {
				ptr.onFocus()
			}
		}
		// Defer to the old window proc

	case win.WM_KILLFOCUS:
		if w := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA); w != 0 {
			ptr := (*mountedTextArea)(unsafe.Pointer(w))
			if ptr.onBlur != nil {
				ptr.onBlur()
			}
		}
		// Defer to the old window proc

	case win.WM_COMMAND:
		notification := win.HIWORD(uint32(wParam))
		switch notification {
		case win.EN_UPDATE:
			if w := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA); w != 0 {
				ptr := (*mountedTextArea)(unsafe.Pointer(w))
				if ptr.onChange != nil {
					ptr.onChange(GetWindowText(hwnd))
				}
			}
		}
		return 0

	}

	return win.CallWindowProc(oldEditWindowProc, hwnd, msg, wParam, lParam)
}
