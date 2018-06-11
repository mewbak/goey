package goey

import (
	"unsafe"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type mountedIntInput struct {
	handle *gtk.SpinButton

	onChange func(int64)
	shChange glib.SignalHandle
	onFocus  focusSlot
	onBlur   blurSlot
}

func (w *IntInput) mount(parent NativeWidget) (Element, error) {
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
		handle:   control,
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
	text := widget.GetValue()
	mounted.onChange(int64(text))
}

func intinput_onDestroy(widget *gtk.SpinButton, mounted *mountedIntInput) {
	mounted.handle = nil
}

func (w *mountedIntInput) Close() {
	if w.handle != nil {
		w.handle.Destroy()
		w.handle = nil
	}
}

func (w *mountedIntInput) Handle() *gtk.Widget {
	return &w.handle.Widget
}

func (w *mountedIntInput) updateProps(data *IntInput) error {
	w.handle.SetValue(float64(data.Value))
	w.handle.SetPlaceholderText(data.Placeholder)
	w.onChange = data.OnChange
	w.shChange = setSignalHandler(&w.handle.Widget, w.shChange, data.OnChange != nil, "value-changed", intinput_onChanged, w)
	w.onFocus.Set(&w.handle.Widget, data.OnFocus)
	w.onBlur.Set(&w.handle.Widget, data.OnBlur)

	return nil
}
