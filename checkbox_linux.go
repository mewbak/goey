package goey

import (
	"unsafe"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type mountedCheckbox struct {
	Control

	onChange func(bool)
	shClick  glib.SignalHandle
	onFocus  focusSlot
	onBlur   blurSlot
}

func (w *Checkbox) mount(parent Control) (Element, error) {
	control, err := gtk.CheckButtonNewWithLabel(w.Text)
	if err != nil {
		return nil, err
	}
	(*gtk.Container)(unsafe.Pointer(parent.handle)).Add(control)
	control.SetActive(w.Value)
	control.SetSensitive(!w.Disabled)

	retval := &mountedCheckbox{
		Control:  Control{&control.Widget},
		onChange: w.OnChange,
	}

	control.Connect("destroy", checkbox_onDestroy, retval)
	retval.shClick = setSignalHandler(&control.Widget, 0, w.OnChange != nil, "clicked", checkbox_onClick, retval)
	retval.onFocus.Set(&control.Widget, w.OnFocus)
	retval.onBlur.Set(&control.Widget, w.OnBlur)
	control.Show()

	return retval, nil
}

func checkbox_onClick(widget *gtk.CheckButton, mounted *mountedCheckbox) {
	if mounted.onChange == nil {
		return
	}

	mounted.onChange(widget.GetActive())
}

func checkbox_onDestroy(widget *gtk.CheckButton, mounted *mountedCheckbox) {
	mounted.handle = nil
}

func (w *mountedCheckbox) checkbutton() *gtk.CheckButton {
	return (*gtk.CheckButton)(unsafe.Pointer(w.handle))
}

func (w *mountedCheckbox) Props() Widget {
	checkbutton := w.checkbutton()
	text, err := checkbutton.GetLabel()
	if err != nil {
		panic("Could not get label: " + err.Error())
	}

	return &Checkbox{
		Value:    checkbutton.GetActive(),
		Text:     text,
		Disabled: !checkbutton.GetSensitive(),
		OnChange: w.onChange,
		OnFocus:  w.onFocus.callback,
		OnBlur:   w.onBlur.callback,
	}
}

func (w *mountedCheckbox) updateProps(data *Checkbox) error {
	checkbutton := w.checkbutton()

	w.onChange = nil // temporarily break OnChange to prevent event
	checkbutton.SetLabel(data.Text)
	checkbutton.SetActive(data.Value)
	checkbutton.SetSensitive(!data.Disabled)

	w.onChange = data.OnChange
	w.shClick = setSignalHandler(&checkbutton.Widget, w.shClick, data.OnChange != nil, "clicked", checkbox_onClick, w)
	w.onFocus.Set(&checkbutton.Widget, data.OnFocus)
	w.onBlur.Set(&checkbutton.Widget, data.OnBlur)

	return nil
}
