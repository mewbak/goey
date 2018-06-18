package goey

import (
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

func (w *TextArea) mount(parent Control) (Element, error) {
	text, err := syscall.UTF16PtrFromString(w.Value)
	if err != nil {
		return nil, err
	}

	style := uint32(win.WS_CHILD | win.WS_VISIBLE | win.WS_TABSTOP | win.ES_LEFT | win.ES_MULTILINE | win.ES_WANTRETURN | win.ES_AUTOVSCROLL)
	hwnd := win.CreateWindowEx(win.WS_EX_CLIENTEDGE, edit.className, text,
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
	subclassWindowProcedure(hwnd, &edit.oldWindowProc, syscall.NewCallback(textinputWindowProc))

	// Create placeholder, if required.
	if w.Value == "" && w.Placeholder != "" {
		textPlaceholder, err := syscall.UTF16PtrFromString(w.Placeholder)
		if err != nil {
			win.DestroyWindow(hwnd)
			return nil, err
		}

		win.SendMessage(hwnd, win.EM_SETCUEBANNER, 0, uintptr(unsafe.Pointer(textPlaceholder)))
	}

	retval := &mountedTextArea{mountedTextInputBase{
		Control:  Control{hwnd},
		onChange: w.OnChange,
		onFocus:  w.OnFocus,
		onBlur:   w.OnBlur,
	},
		minlinesDefault(w.MinLines),
	}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

type mountedTextArea struct {
	mountedTextInputBase
	minLines int
}

func (w *mountedTextArea) MeasureHeight(width Length) (Length, Length) {
	const lineHeight = 16 * DIP
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	return 23*DIP + lineHeight.Scale(w.minLines-1, 1), 23*DIP + lineHeight.Scale(39, 1)
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
