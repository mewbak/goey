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
	if w.ReadOnly {
		style = style | win.ES_READONLY
	}
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

func (w *mountedTextArea) Layout(bc Constraint) Size {
	width := w.MinIntrinsicWidth(0)
	height := w.MinIntrinsicHeight(0)
	return bc.Constrain(Size{width, height})
}

func (w *mountedTextArea) MinIntrinsicHeight(Length) Length {
	// https://msdn.microsoft.com/en-us/library/windows/desktop/dn742486.aspx#sizingandspacing
	const lineHeight = 16 * DIP
	return 23*DIP + lineHeight.Scale(w.minLines-1, 1)
}

func (w *mountedTextArea) Props() Widget {
	var buffer [80]uint16
	win.SendMessage(w.hWnd, win.EM_GETCUEBANNER, uintptr(unsafe.Pointer(&buffer[0])), 80)
	ndx := 0
	for i, v := range buffer {
		if v == 0 {
			ndx = i
			break
		}
	}
	placeholder := syscall.UTF16ToString(buffer[:ndx])

	return &TextArea{
		Value:       w.Control.Text(),
		Placeholder: placeholder,
		Disabled:    !win.IsWindowEnabled(w.hWnd),
		ReadOnly:    (win.GetWindowLong(w.hWnd, win.GWL_STYLE) & win.ES_READONLY) != 0,
		MinLines:    w.minLines,
		OnChange:    w.onChange,
		OnFocus:     w.onFocus,
		OnBlur:      w.onBlur,
	}
}

func (w *mountedTextArea) updateProps(data *TextArea) error {
	if data.Value != w.Text() {
		w.SetText(data.Value)
	}
	err := w.updatePlaceholder(data.Placeholder)
	if err != nil {
		return err
	}
	w.SetDisabled(data.Disabled)
	win.SendMessage(w.hWnd, win.EM_SETREADONLY, uintptr(win.BoolToBOOL(data.ReadOnly)), 0)

	w.minLines = minlinesDefault(data.MinLines)
	w.onChange = data.OnChange
	w.onFocus = data.OnFocus
	w.onBlur = data.OnBlur

	return nil
}
