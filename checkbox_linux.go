package goey

import (
	"unsafe"

	"github.com/gotk3/gotk3/gtk"
)

type mountedCheckbox struct {
	NativeWidget

	onChange func(bool)
	onFocus  func()
	onBlur   func()
}

func (w *Checkbox) Mount(parent NativeWidget) (MountedWidget, error) {
	control, err := gtk.CheckButtonNewWithLabel(w.Text)
	if err != nil {
		return nil, err
	}
	(*gtk.Container)(unsafe.Pointer(parent.handle)).Add(control)

	retval := &mountedCheckbox{
		NativeWidget: NativeWidget{&control.Bin.Container.Widget},
		onChange:     w.OnChange,
		onFocus:      nil,
		onBlur:       nil,
	}

	if w.OnChange != nil {
		control.Connect("clicked", checkbox_onClicked, retval)
	}
	control.Connect("destroy", checkbox_onDestroy, retval)
	control.Show()

	return retval, nil
}

func checkbox_onClicked(widget *gtk.CheckButton, mounted *mountedCheckbox) {
	mounted.onChange(widget.GetActive())
}

func checkbox_onDestroy(widget *gtk.CheckButton, mounted *mountedCheckbox) {
	mounted.handle = nil
}

func (w *mountedCheckbox) UpdateProps(data Widget) error {
	panic("not implemented")
}
