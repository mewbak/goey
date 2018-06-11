package goey

import (
	win2 "bitbucket.org/rj/goey/syscall"
	"github.com/lxn/win"
	"syscall"
	"unsafe"
)

var (
	editClassName     *uint16
	oldEditWindowProc uintptr
)

func init() {
	var err error
	editClassName, err = syscall.UTF16PtrFromString("EDIT")
	if err != nil {
		panic(err)
	}
}

func (w *TextInput) mount(parent NativeWidget) (Element, error) {
	text, err := syscall.UTF16PtrFromString(w.Value)
	if err != nil {
		return nil, err
	}

	style := uint32(win.WS_CHILD | win.WS_VISIBLE | win.WS_TABSTOP | win.ES_LEFT | win.ES_AUTOHSCROLL)
	if w.Password {
		style = style | win.ES_PASSWORD
	}
	if w.OnEnterKey != nil {
		style = style | win.ES_MULTILINE
	}
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
	subclassWindowProcedure(hwnd, &oldEditWindowProc, syscall.NewCallback(textinputWindowProc))

	// Create placeholder, if required.
	if w.Placeholder != "" {
		textPlaceholder, err := syscall.UTF16PtrFromString(w.Placeholder)
		if err != nil {
			win.DestroyWindow(hwnd)
			return nil, err
		}

		win.SendMessage(hwnd, win.EM_SETCUEBANNER, 0, uintptr(unsafe.Pointer(textPlaceholder)))
	}

	retval := &mountedTextInput{mountedTextInputBase{
		NativeWidget: NativeWidget{hwnd},
		onChange:     w.OnChange,
		onFocus:      w.OnFocus,
		onBlur:       w.OnBlur,
		onEnterKey:   w.OnEnterKey,
	}}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

type mountedTextInputBase struct {
	NativeWidget
	onChange   func(value string)
	onFocus    func()
	onBlur     func()
	onEnterKey func(value string)
}

type mountedTextInput struct {
	mountedTextInputBase
}

func (w *mountedTextInputBase) MeasureWidth() (Length, Length) {
	if paragraphMaxWidth == 0 {
		paragraphMeasureReflowLimits(w.hWnd)
	}

	// In the future, we should calculate the width based on the length of the text.
	return FromPixelsX(paragraphMinWidth), FromPixelsX(paragraphMaxWidth)
}

func (w *mountedTextInputBase) MeasureHeight(width Length) (Length, Length) {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	return 23 * DIP, 23 * DIP
}

func (w *mountedTextInputBase) updateProps(data *TextInput) error {
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

	w.onChange = data.OnChange
	w.onFocus = data.OnFocus
	w.onBlur = data.OnBlur
	w.onEnterKey = data.OnEnterKey

	return nil
}

func textinputWindowProc(hwnd win.HWND, msg uint32, wParam uintptr, lParam uintptr) (result uintptr) {
	switch msg {
	case win.WM_DESTROY:
		// Make sure that the data structure on the Go-side does not point to a non-existent
		// window.
		if w := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA); w != 0 {
			ptr := (*mountedTextInputBase)(unsafe.Pointer(w))
			ptr.hWnd = 0
		}
		// Defer to the old window proc

	case win.WM_SETFOCUS:
		if w := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA); w != 0 {
			ptr := (*mountedTextInputBase)(unsafe.Pointer(w))
			if ptr.onFocus != nil {
				ptr.onFocus()
			}
		}
		// Defer to the old window proc

	case win.WM_KILLFOCUS:
		if w := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA); w != 0 {
			ptr := (*mountedTextInputBase)(unsafe.Pointer(w))
			if ptr.onBlur != nil {
				ptr.onBlur()
			}
		}
		// Defer to the old window proc

	case win.WM_KEYDOWN:
		if wParam == win.VK_RETURN {
			if w := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA); w != 0 {
				ptr := (*mountedTextInputBase)(unsafe.Pointer(w))
				if ptr.onEnterKey != nil {
					ptr.onEnterKey(win2.GetWindowText(hwnd))
					return 0
				}
			}
		}
		// Defer to the old window proc

	case win.WM_COMMAND:
		notification := win.HIWORD(uint32(wParam))
		switch notification {
		case win.EN_UPDATE:
			if w := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA); w != 0 {
				ptr := (*mountedTextInputBase)(unsafe.Pointer(w))
				if ptr.onChange != nil {
					ptr.onChange(win2.GetWindowText(hwnd))
				}
			}
		}
		return 0

	}

	return win.CallWindowProc(oldEditWindowProc, hwnd, msg, wParam, lParam)
}
