package goey

import (
	"unsafe"

	"github.com/gotk3/gotk3/gtk"
)

type mountedPasswordInput struct {
	mountedTextInput
}

func (w *PasswordInput) mount(parent NativeWidget) (MountedWidget, error) {
	control, err := gtk.EntryNew()
	if err != nil {
		return nil, err
	}
	(*gtk.Container)(unsafe.Pointer(parent.handle)).Add(control)
	control.SetText(w.Value)
	control.SetPlaceholderText(w.Placeholder)
	control.SetVisibility(false)

	retval := &mountedPasswordInput{mountedTextInput{
		handle:   control,
		onChange: w.OnChange,
	}}

	control.Connect("destroy", textinput_onDestroy, &retval.mountedTextInput)
	retval.shChange = setSignalHandler(&control.Widget, 0, retval.onChange != nil, "changed", textinput_onChanged, &retval.mountedTextInput)
	retval.onFocus.Set(&control.Widget, w.OnFocus)
	retval.onBlur.Set(&control.Widget, w.OnBlur)
	control.Show()

	return retval, nil
}

func (w *mountedPasswordInput) updateProps(data *PasswordInput) error {
	w.handle.SetText(data.Value)
	w.handle.SetPlaceholderText(data.Placeholder)
	w.onChange = data.OnChange
	w.shChange = setSignalHandler(&w.handle.Widget, w.shChange, data.OnChange != nil, "changed", textinput_onChanged, &w.mountedTextInput)
	w.onFocus.Set(&w.handle.Widget, data.OnFocus)
	w.onBlur.Set(&w.handle.Widget, data.OnBlur)

	return nil
}
