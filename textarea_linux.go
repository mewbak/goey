package goey

import (
	"unsafe"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type mountedTextArea struct {
	handle *gtk.TextView

	onChange func(string)
	shChange glib.SignalHandle
	onFocus  focusSlot
	onBlur   blurSlot
}

func (w *TextArea) mount(parent NativeWidget) (MountedWidget, error) {
	buffer, err := gtk.TextBufferNew(nil)
	if err != nil {
		return nil, err
	}
	buffer.SetText(w.Value)

	control, err := gtk.TextViewNewWithBuffer(buffer)
	if err != nil {
		buffer.Unref()
		return nil, err
	}
	control.SetWrapMode(gtk.WRAP_WORD)
	buffer.Unref()
	(*gtk.Container)(unsafe.Pointer(parent.handle)).Add(control)

	retval := &mountedTextArea{
		handle:   control,
		onChange: w.OnChange,
	}

	control.Connect("destroy", textarea_onDestroy, retval)
	if w.OnChange != nil {
		sh, err := buffer.Connect("changed", textarea_onChanged, retval)
		if err != nil {
			panic("Failed to connect 'changed' event")
		}
		retval.shChange = sh
	}
	retval.onFocus.Set(&control.Widget, w.OnFocus)
	retval.onBlur.Set(&control.Widget, w.OnBlur)
	control.Show()

	return retval, nil
}

func textarea_onChanged(buffer *gtk.TextBuffer, mounted *mountedTextArea) {
	text, err := buffer.GetText(buffer.GetStartIter(), buffer.GetEndIter(), true)
	if err != nil {
		// TODO:  What is the correct reporting here
		return
	}
	mounted.onChange(text)
}

func textarea_onDestroy(widget *gtk.TextView, mounted *mountedTextArea) {
	mounted.handle = nil
}

func (w *mountedTextArea) Close() {
	if w.handle != nil {
		w.handle.Destroy()
		w.handle = nil
	}
}

func (w *mountedTextArea) Handle() *gtk.Widget {
	return &w.handle.Widget
}

func (w *mountedTextArea) updateProps(data *TextArea) error {
	// TextView will send a 'changed' event, even if the new value is the
	// same.  To stop an infinite loop, we need to protect by checking
	// ourselves.
	buffer, err := w.handle.GetBuffer()
	if err != nil {
		return err
	}
	oldText, err := buffer.GetText(buffer.GetStartIter(), buffer.GetEndIter(), true)
	if err != nil {
		return err
	}
	if data.Value != oldText {
		buffer.SetText(data.Value)
	}

	w.onChange = data.OnChange
	if data.OnChange != nil && w.shChange == 0 {
		sh, err := buffer.Connect("changed", textinput_onChanged, w)
		if err != nil {
			panic("Failed to connect 'changed' event")
		}
		w.shChange = sh
	} else if data.OnChange == nil && w.shChange != 0 {
		buffer.HandlerDisconnect(w.shChange)
		w.shChange = 0
	}
	w.onFocus.Set(&w.handle.Widget, data.OnFocus)
	w.onBlur.Set(&w.handle.Widget, data.OnBlur)

	return nil
}
