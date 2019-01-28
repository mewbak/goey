package goey

import (
	"strconv"
	"syscall"
	"unsafe"

	"bitbucket.org/rj/goey/base"
	"github.com/lxn/win"
)

var (
	intinput struct {
		className     []uint16
		oldWindowProc uintptr
	}
)

func init() {
	intinput.className = []uint16{'m', 's', 'c', 't', 'l', 's', '_', 'u', 'p', 'd', 'o', 'w', 'n', '3', '2', 0}
}

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

	// Create the updown control.
	// Range for the updown control is is only int32, not int64.
	hwndUpDown := win.HWND(0)
	if w.useUpDownControl() {
		hwndUpDown, _, err = createControlWindow(win.WS_EX_LEFT|win.WS_EX_LTRREADING,
			&intinput.className[0],
			"",
			win.WS_CHILDWINDOW|win.WS_VISIBLE|win.UDS_SETBUDDYINT|win.UDS_ARROWKEYS|win.UDS_HOTTRACK|win.UDS_NOTHOUSANDS,
			parent.HWnd)
		if err != nil {
			return nil, err
		}
		win.SendMessage(hwndUpDown, win.UDM_SETRANGE32, uintptr(w.Min), uintptr(w.Max))
		win.SendMessage(hwndUpDown, win.UDM_SETPOS32, 0, uintptr(w.Value))
		win.SendMessage(hwndUpDown, win.UDM_SETBUDDY, uintptr(hwnd), 0)
	}

	if w.Disabled {
		win.EnableWindow(hwnd, false)
	}

	// Subclass the window procedure
	if hwndUpDown != 0 {
		subclassWindowProcedure(hwnd, &intinput.oldWindowProc, textinputWindowProc)
	} else {
		subclassWindowProcedure(hwnd, &edit.oldWindowProc, textinputWindowProc)
	}

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
		hwndUpDown: hwndUpDown,
		min:        w.Min,
		max:        w.Max,
		onChange:   w.OnChange,
		onEnterKey: w.OnEnterKey,
	}
	retval.setThunks()
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(retval)))

	return retval, nil
}

func (w *IntInput) useUpDownControl() bool {
	// Range for the updown control is is only int32, not int64.
	// Need to make sure that we can properly set the range for the updown
	// control in order to include it in the GUI.
	return w.Min >= -2147483648 && w.Max <= 2147483647
}

type intinputElement struct {
	textinputElementBase

	hwndUpDown win.HWND
	min        int64
	max        int64
	onChange   func(int64)
	onEnterKey func(int64)
}

func (w *intinputElement) Close() {
	if w.hwndUpDown != 0 {
		win.DestroyWindow(w.hwndUpDown)
		w.hwndUpDown = 0
	}
	if w.hWnd != 0 {
		win.DestroyWindow(w.hWnd)
		w.hWnd = 0
	}
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
	// Clamp the value
	if i < w.min {
		i = w.min
	} else if i > w.max {
		i = w.max
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
	// Clamp the value
	if i < w.min {
		i = w.min
	} else if i > w.max {
		i = w.max
	}
	// With conversion completed, call original callback.
	w.onEnterKey(i)
}

func (w *intinputElement) Props() base.Widget {
	value := int64(0)
	if w.hwndUpDown != 0 {
		value = int64(win.SendMessage(w.hwndUpDown, win.UDM_GETPOS32, 0, 0))
	} else {
		value, _ = strconv.ParseInt(w.Control.Text(), 10, 64)
	}

	return &IntInput{
		Value:       value,
		Placeholder: w.propsPlaceholder(),
		Disabled:    !win.IsWindowEnabled(w.hWnd),
		Min:         w.min,
		Max:         w.max,
		OnChange:    w.onChange,
		OnFocus:     w.onFocus,
		OnBlur:      w.onBlur,
		OnEnterKey:  w.onEnterKey,
	}
}

func (w *intinputElement) SetBounds(bounds base.Rectangle) {
	buddyWidth := (23 * DIP) * 2 / 3

	if w.hwndUpDown == 0 {
		win.MoveWindow(w.hWnd, int32(bounds.Min.X.PixelsX()), int32(bounds.Min.Y.PixelsY()), int32(bounds.Dx().PixelsX()), int32(bounds.Dy().PixelsY()), false)
		return
	}

	if bounds.Dx() >= 4*buddyWidth {
		win.MoveWindow(w.hWnd, int32(bounds.Min.X.PixelsX()), int32(bounds.Min.Y.PixelsY()), int32((bounds.Dx() - buddyWidth).PixelsX()), int32(bounds.Dy().PixelsY()), false)
		win.MoveWindow(w.hwndUpDown, int32((bounds.Max.X - buddyWidth).PixelsX()), int32(bounds.Min.Y.PixelsY()), int32(buddyWidth.PixelsX()), int32(bounds.Dy().PixelsY()), false)
		win.ShowWindow(w.hwndUpDown, win.SW_SHOW)
	} else {
		win.MoveWindow(w.hWnd, int32(bounds.Min.X.PixelsX()), int32(bounds.Min.Y.PixelsY()), int32(bounds.Dx().PixelsX()), int32(bounds.Dy().PixelsY()), false)
		win.ShowWindow(w.hwndUpDown, win.SW_HIDE)
	}
}

func (w *intinputElement) updateProps(data *IntInput) error {
	// Remove the updown control is the range is too large.
	if w.hwndUpDown != 0 && !data.useUpDownControl() {
		win.DestroyWindow(w.hwndUpDown)
		w.hwndUpDown = 0
	}

	text := strconv.FormatInt(data.Value, 10)
	if text != w.Text() {
		w.SetText(text)
	}
	if w.hwndUpDown != 0 {
		win.SendMessage(w.hwndUpDown, win.UDM_SETRANGE32, uintptr(data.Min), uintptr(data.Max))
		win.SendMessage(w.hwndUpDown, win.UDM_SETPOS32, 0, uintptr(data.Value))
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
