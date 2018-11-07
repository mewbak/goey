package goey

import (
	"strconv"
	"syscall"
	"unsafe"

	"bitbucket.org/rj/goey/base"
	"github.com/lxn/win"
)

func (w *IntInput) mount(parent base.Control) (base.Element, error) {
	// Create the control
	style := uint32(win.WS_CHILD | win.WS_VISIBLE | win.WS_TABSTOP | win.ES_LEFT | win.ES_AUTOHSCROLL | win.ES_NUMBER)
	if w.OnEnterKey != nil {
		style = style | win.ES_MULTILINE
	}
	hwnd, _, err := createControlWindow(win.WS_EX_CLIENTEDGE, &edit.className[0], strconv.FormatInt(w.Value, 10), style, parent.HWnd)
	if err != nil {
		return nil, err
	}
	if w.Disabled {
		win.EnableWindow(hwnd, false)
	}

	// Subclass the window procedure
	subclassWindowProcedure(hwnd, &edit.oldWindowProc, textinputWindowProc)

	// Create placeholder, if required.
	if w.Placeholder != "" {
		textPlaceholder, err := syscall.UTF16PtrFromString(w.Placeholder)
		if err != nil {
			win.DestroyWindow(hwnd)
			return nil, err
		}

		win.SendMessage(hwnd, win.EM_SETCUEBANNER, 0, uintptr(unsafe.Pointer(textPlaceholder)))
	}

	retval := &intinputElement{
		textinputElementBase: textinputElementBase{
			Control: Control{hwnd},
			onFocus: w.OnFocus,
			onBlur:  w.OnBlur,
		},
		onChange:   w.OnChange,
		onEnterKey: w.OnEnterKey,
	}
	retval.setThunks()
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

type intinputElement struct {
	textinputElementBase
	onChange   func(int64)
	onEnterKey func(int64)
}

// We are delegating a lot of behaviour to the textinput element.  However,
// callbacks need to convert from string to int64.  This updates the callbacks
// in the textinput to thunks that will do the necessary conversions.
func (w *intinputElement) setThunks() {
	if w.onChange != nil {
		w.textinputElementBase.onChange = w.thunkOnChange
	} else {
		w.textinputElementBase.onChange = nil
	}
	if w.onEnterKey != nil {
		w.textinputElementBase.onEnterKey = w.thunkOnEnterKey
	} else {
		w.textinputElementBase.onEnterKey = nil
	}
}

func (w *intinputElement) thunkOnChange(value string) {
	// Convert text from control to an integer
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		// This case should not occur, as the control should prevent invalid
		// strings from being entered.
		// TODO:  What reporting should be done here?
		return
	}
	// With conversion completed, call original callback.
	w.onChange(i)
}

func (w *intinputElement) thunkOnEnterKey(value string) {
	// Convert text from control to an integer
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		// This case should not occur, as the control should prevent invalid
		// strings from being entered.
		// TODO:  What reporting should be done here?
		return
	}
	// With conversion completed, call original callback.
	w.onEnterKey(i)
}

func (w *intinputElement) Props() base.Widget {
	value, err := strconv.ParseInt(w.Control.Text(), 10, 64)
	if err != nil {
		return nil
	}

	return &IntInput{
		Value:       value,
		Placeholder: w.propsPlaceholder(),
		Disabled:    !win.IsWindowEnabled(w.hWnd),
		OnChange:    w.onChange,
		OnFocus:     w.onFocus,
		OnBlur:      w.onBlur,
		OnEnterKey:  w.onEnterKey,
	}
}

func (w *intinputElement) updateProps(data *IntInput) error {
	text := strconv.FormatInt(data.Value, 10)
	if text != w.Text() {
		w.SetText(text)
	}
	err := w.updatePlaceholder(data.Placeholder)
	if err != nil {
		return err
	}
	w.SetDisabled(data.Disabled)

	w.onChange = data.OnChange
	w.onFocus = data.OnFocus
	w.onBlur = data.OnBlur
	w.onEnterKey = data.OnEnterKey
	w.setThunks()

	return nil
}
