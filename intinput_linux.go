package goey

import (
	"unsafe"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type mountedIntInput struct {
	Control

	onChange func(int64)
	shChange glib.SignalHandle
	onFocus  focusSlot
	onBlur   blurSlot
}

func (w *IntInput) mount(parent Control) (Element, error) {
	control, err := gtk.SpinButtonNew(nil, 1, 0)
	if err != nil {
		return nil, err
	}
	(*gtk.Container)(unsafe.Pointer(parent.handle)).Add(control)
	control.SetRange(-0x7FFFFFFFFFFFFFFF, 0x7FFFFFFFFFFFFFFF)
	control.SetValue(float64(w.Value))
	control.SetIncrements(1, 10)
	control.SetPlaceholderText(w.Placeholder)

	retval := &mountedIntInput{
		Control:  Control{&control.Widget},
		onChange: w.OnChange,
	}

	control.Connect("destroy", intinput_onDestroy, retval)
	retval.shChange = setSignalHandler(&control.Widget, 0, retval.onChange != nil, "value-changed", intinput_onChanged, retval)
	retval.onFocus.Set(&control.Widget, w.OnFocus)
	retval.onBlur.Set(&control.Widget, w.OnBlur)
	control.Show()

	return retval, nil
}

func intinput_onChanged(widget *gtk.SpinButton, mounted *mountedIntInput) {
	if mounted.onChange == nil {
		return
	}

	text := widget.GetValue()
	mounted.onChange(int64(text))
}

func intinput_onDestroy(widget *gtk.SpinButton, mounted *mountedIntInput) {
	mounted.handle = nil
}

func (w *mountedIntInput) spinbutton() *gtk.SpinButton {
	return (*gtk.SpinButton)(unsafe.Pointer(w.handle))
}

func (w *mountedIntInput) updateProps(data *IntInput) error {
	button := w.spinbutton()

	w.onChange = nil // break OnChange to prevent event
	button.SetValue(float64(data.Value))
	button.SetPlaceholderText(data.Placeholder)
	w.onChange = data.OnChange
	w.shChange = setSignalHandler(&button.Widget, w.shChange, data.OnChange != nil, "value-changed", intinput_onChanged, w)
	w.onFocus.Set(&button.Widget, data.OnFocus)
	w.onBlur.Set(&button.Widget, data.OnBlur)

	return nil
}
