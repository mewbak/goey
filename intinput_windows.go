package goey

import (
	"strconv"
	"syscall"
	"unsafe"

	"bitbucket.org/rj/goey/base"
	"github.com/lxn/win"
)

func (w *IntInput) mount(parent base.Control) (base.Element, error) {
	text, err := syscall.UTF16PtrFromString(strconv.FormatInt(w.Value, 10))
	if err != nil {
		return nil, err
	}

	style := uint32(win.WS_CHILD | win.WS_VISIBLE | win.WS_TABSTOP | win.ES_LEFT | win.ES_AUTOHSCROLL | win.ES_NUMBER)
	if w.OnEnterKey != nil {
		style = style | win.ES_MULTILINE
	}
	hwnd := win.CreateWindowEx(win.WS_EX_CLIENTEDGE, &edit.className[0], text,
		style,
		10, 10, 100, 100,
		parent.HWnd, win.HMENU(nextControlID()), 0, nil)
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
	subclassWindowProcedure(hwnd, &edit.oldWindowProc, syscall.NewCallback(textinputWindowProc))

	// Create placeholder, if required.
	if w.Placeholder != "" {
		textPlaceholder, err := syscall.UTF16PtrFromString(w.Placeholder)
		if err != nil {
			win.DestroyWindow(hwnd)
			return nil, err
		}

		win.SendMessage(hwnd, win.EM_SETCUEBANNER, 0, uintptr(unsafe.Pointer(textPlaceholder)))
	}

	retval := &intinputElement{textinputElementBase{
		Control:    Control{hwnd},
		onChange:   w.wrapOnChange(),
		onFocus:    w.OnFocus,
		onBlur:     w.OnBlur,
		onEnterKey: w.OnEnterKey,
	}}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

func (w *IntInput) wrapOnChange() func(string) {
	if w.OnChange == nil {
		return nil
	}

	return func(value string) {
		// Convert text from control to an integer
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			// This case should not occur, as the control should prevent invalid
			// strings from being entered.
			// TODO:  What reporting should be done here?
			return
		}
		// With conversion completed, call original callback.
		w.OnChange(i)
	}
}

type intinputElement struct {
	textinputElementBase
}

func (w *intinputElement) updateProps(data *IntInput) error {
	text := strconv.FormatInt(data.Value, 10)
	if text != w.Text() {
		w.SetText(text)
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

	w.onChange = data.wrapOnChange()
	w.onFocus = data.OnFocus
	w.onBlur = data.OnBlur
	w.onEnterKey = data.OnEnterKey

	return nil
}
