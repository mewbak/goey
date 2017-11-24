package goey

import (
	"unsafe"

	"github.com/gotk3/gotk3/gtk"
)

type mountedTextArea struct {
	NativeWidget

	onChange func(string)
	onFocus  func()
	onBlur   func()
}

func (w *TextArea) Mount(parent NativeWidget) (MountedWidget, error) {
	buffer, err := gtk.TextBufferNew(nil)
	if err != nil {
		return nil, err
	}
	buffer.SetText(w.Text)

	control, err := gtk.TextViewNewWithBuffer(buffer)
	if err != nil {
		buffer.Unref()
		return nil, err
	}
	buffer.Unref()

	//frame, err := gtk.FrameNew("")
	//if err != nil {
	//	control.Destroy()
	//	return nil, err
	//}
	//frame.Add(control)
	(*gtk.Container)(unsafe.Pointer(parent.handle)).Add(control)

	retval := &mountedTextArea{
		NativeWidget: NativeWidget{&control.Widget},
		onChange:     w.OnChange,
		onFocus:      w.OnFocus,
		onBlur:       w.OnBlur,
	}

	if w.OnChange != nil {
		control.Connect("changed", textarea_onChanged, retval)
	}
	control.Connect("destroy", textarea_onDestroy, retval)
	control.Show()

	return retval, nil
}

func textarea_onChanged(widget *gtk.TextView, mounted *mountedTextArea) {
	buffer, err := widget.GetBuffer()
	if err != nil {
		// TODO:  What is the correct reporting here
		return
	}
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

func (w *mountedTextArea) UpdateProps(data Widget) error {
	panic("not implemented")
}
