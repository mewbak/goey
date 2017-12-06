package goey

import (
	"unsafe"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type mountedButton struct {
	handle *gtk.Button

	onClick clickSlot
	onFocus focusSlot
	onBlur  blurSlot
}

func (w *Button) mount(parent NativeWidget) (MountedWidget, error) {
	control, err := gtk.ButtonNewWithLabel(w.Text)
	if err != nil {
		return nil, err
	}
	control.AddEvents(int(gdk.FOCUS_CHANGE_MASK))

	(*gtk.Container)(unsafe.Pointer(parent.handle)).Add(control)
	control.SetSensitive(!w.Disabled)
	control.SetCanDefault(true)
	if w.Default {
		control.GrabDefault()
	}

	retval := &mountedButton{
		handle: control,
	}

	control.Connect("destroy", button_onDestroy, retval)
	err = retval.onClick.Set(&control.Widget, w.OnClick)
	if err != nil {
		control.Destroy()
		return nil, err
	}
	err = retval.onFocus.Set(&control.Widget, w.OnFocus)
	if err != nil {
		control.Destroy()
		return nil, err
	}
	err = retval.onBlur.Set(&control.Widget, w.OnBlur)
	if err != nil {
		control.Destroy()
		return nil, err
	}
	control.Show()

	return retval, nil
}

func button_onDestroy(widget *gtk.Button, mounted *mountedButton) {
	mounted.handle = nil
}

func (w *mountedButton) Close() {
	if w.handle != nil {
		w.onClick.Close(&w.handle.Widget)
		w.onFocus.Close(&w.handle.Widget)
		w.onBlur.Close(&w.handle.Widget)
		w.handle.Destroy()
		w.handle = nil
	}
}

func (w *mountedButton) updateProps(data *Button) error {
	label_, err := w.handle.GetChild()
	if err != nil {
		return err
	}

	(*gtk.Label)(unsafe.Pointer(label_)).SetText(data.Text)
	w.handle.SetSensitive(!data.Disabled)

	if data.Default {
		w.handle.GrabDefault()
	}
	err = w.onClick.Set(&w.handle.Widget, data.OnClick)
	if err != nil {
		return err
	}
	err = w.onFocus.Set(&w.handle.Widget, data.OnFocus)
	if err != nil {
		return err
	}
	err = w.onBlur.Set(&w.handle.Widget, data.OnBlur)
	if err != nil {
		return err
	}

	return nil
}
