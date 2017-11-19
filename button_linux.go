package goey

import (
	"unsafe"

	"github.com/gotk3/gotk3/gtk"
)

type mountedButton struct {
	handle *gtk.Button

	onClick func()
	onFocus func()
	onBlur  func()
}

func (w *Button) Mount(parent NativeWidget) (MountedWidget, error) {
	control, err := gtk.ButtonNewWithLabel(w.Text)
	if err != nil {
		return nil, err
	}
	(*gtk.Container)(unsafe.Pointer(parent.handle)).Add(control)
	control.SetSensitive(!w.Disabled)

	retval := &mountedButton{
		handle:  control,
		onClick: w.OnClick,
		onFocus: nil,
		onBlur:  nil,
	}

	if w.OnClick != nil {
		control.Connect("clicked", button_onClick, retval)
	}
	control.Connect("destroy", button_onDestroy, retval)
	control.Show()

	return retval, nil
}

func button_onClick(widget *gtk.Button, mounted *mountedButton) {
	mounted.onClick()
}

func button_onDestroy(widget *gtk.Button, mounted *mountedButton) {
	mounted.handle = nil
}

func (w *mountedButton) Close() {
	if w.handle != nil {
		w.handle.Destroy()
		w.handle = nil
	}
}

func (w *mountedButton) UpdateProps(data_ Widget) error {
	data := data_.(*Button)

	label_, err := w.handle.GetChild()
	if err != nil {
		return err
	}

	(*gtk.Label)(unsafe.Pointer(label_)).SetText(data.Text)
	w.handle.SetSensitive(!data.Disabled)
	// TODO:  Update property .Default
	w.onClick = data.OnClick
	w.onFocus = data.OnFocus
	w.onBlur = data.OnBlur

	return nil
}
