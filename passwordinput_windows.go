package goey

import (
	"github.com/lxn/win"
	"syscall"
	"unsafe"
)

func (w *PasswordInput) mount(parent NativeWidget) (MountedWidget, error) {
	text, err := syscall.UTF16PtrFromString(w.Value)
	if err != nil {
		return nil, err
	}

	style := uint32(win.WS_CHILD | win.WS_VISIBLE | win.WS_TABSTOP | win.ES_LEFT | win.ES_AUTOHSCROLL | win.ES_PASSWORD)
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

	retval := &mountedPasswordInput{mountedTextInputBase{
		NativeWidget: NativeWidget{hwnd},
		onChange:     w.OnChange,
		onFocus:      w.OnFocus,
		onBlur:       w.OnBlur,
		onEnterKey:   w.OnEnterKey,
	}}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

type mountedPasswordInput struct {
	mountedTextInputBase
}

func (w *mountedPasswordInput) updateProps(data *PasswordInput) error {
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
