package goey

import (
	"unsafe"

	"github.com/gotk3/gotk3/gtk"
)

type mountedSelectInput struct {
	NativeWidget

	onChange func(int)
	onFocus  func()
	onBlur   func()
}

func (w *SelectInput) mount(parent NativeWidget) (MountedWidget, error) {
	control, err := gtk.ComboBoxTextNew()
	if err != nil {
		return nil, err
	}
	(*gtk.Container)(unsafe.Pointer(parent.handle)).Add(control)
	for _, v := range w.Items {
		control.AppendText(v)
	}

	retval := &mountedSelectInput{
		NativeWidget: NativeWidget{&control.Widget},
		onChange:     w.OnChange,
		onFocus:      w.OnFocus,
		onBlur:       w.OnBlur,
	}

	if w.OnChange != nil {
		control.Connect("changed", selectinput_onChanged, retval)
	}
	control.Connect("destroy", selectinput_onDestroy, retval)
	control.Show()

	return retval, nil
}

func selectinput_onChanged(widget *gtk.ComboBoxText, mounted *mountedSelectInput) {
	mounted.onChange(widget.GetActive())
}

func selectinput_onDestroy(widget *gtk.ComboBoxText, mounted *mountedSelectInput) {
	mounted.handle = nil
}

func (w *mountedSelectInput) updateProps(data *SelectInput) error {
	panic("not implemented")
}
