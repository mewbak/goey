package goey

import (
	"unsafe"

	"github.com/gotk3/gotk3/gtk"
)

type mountedTextInput struct {
	NativeWidget

	onChange func(string)
	onFocus  func()
	onBlur   func()
}

func (w *TextInput) Mount(parent NativeWidget) (MountedWidget, error) {
	control, err := gtk.EntryNew()
	if err != nil {
		return nil, err
	}
	(*gtk.Container)(unsafe.Pointer(parent.handle)).Add(control)
	control.SetText(w.Text)
	control.SetPlaceholderText(w.Placeholder)

	retval := &mountedTextInput{
		NativeWidget: NativeWidget{&control.Widget},
		onChange:     w.OnChange,
		onFocus:      w.OnFocus,
		onBlur:       w.OnBlur,
	}

	if w.OnChange != nil {
		control.Connect("changed", textinput_onChanged, retval)
	}
	control.Connect("destroy", textinput_onDestroy, retval)
	control.Show()

	return retval, nil
}

func textinput_onChanged(widget *gtk.Entry, mounted *mountedTextInput) {
	text, err := widget.GetText()
	if err != nil {
		// TODO:  What is the correct reporting here
		return
	}
	mounted.onChange(text)
}

func textinput_onDestroy(widget *gtk.Entry, mounted *mountedTextInput) {
	mounted.handle = nil
}

func (w *mountedTextInput) UpdateProps(data Widget) error {
	panic("not implemented")
}
