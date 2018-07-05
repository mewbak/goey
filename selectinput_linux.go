package goey

import (
	"unsafe"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type mountedSelectInput struct {
	Control

	onChange func(int)
	shChange glib.SignalHandle
	onFocus  focusSlot
	onBlur   blurSlot
}

func (w *SelectInput) mount(parent Control) (Element, error) {
	control, err := gtk.ComboBoxTextNew()
	if err != nil {
		return nil, err
	}
	(*gtk.Container)(unsafe.Pointer(parent.handle)).Add(control)
	for _, v := range w.Items {
		control.AppendText(v)
	}
	control.SetActive(w.Value)
	control.SetCanFocus(true)
	control.SetSensitive(!w.Disabled)

	retval := &mountedSelectInput{
		Control:  Control{&control.Widget},
		onChange: w.OnChange,
	}

	control.Connect("destroy", selectinput_onDestroy, retval)
	retval.shChange = setSignalHandler(&control.Widget, 0, w.OnChange != nil, "changed", selectinput_onChanged, retval)
	retval.onFocus.Set(&control.Widget, w.OnFocus)
	retval.onBlur.Set(&control.Widget, w.OnBlur)
	control.Show()

	return retval, nil
}

func selectinput_onChanged(widget *gtk.ComboBoxText, mounted *mountedSelectInput) {
	if mounted.onChange == nil {
		return
	}

	mounted.onChange(widget.GetActive())
}

func selectinput_onDestroy(widget *gtk.ComboBoxText, mounted *mountedSelectInput) {
	mounted.handle = nil
}

func (w *mountedSelectInput) comboboxtext() *gtk.ComboBoxText {
	return (*gtk.ComboBoxText)(unsafe.Pointer(w.handle))
}

func (w *mountedSelectInput) updateProps(data *SelectInput) error {
	cbt := w.comboboxtext()

	w.onChange = nil // temporarily break OnChange to prevent event
	// Todo, can we avoid rebuilding the list?
	cbt.RemoveAll()
	for _, v := range data.Items {
		cbt.AppendText(v)
	}
	cbt.SetActive(data.Value)

	cbt.SetSensitive(!data.Disabled)
	w.onChange = data.OnChange
	w.shChange = setSignalHandler(&cbt.Widget, w.shChange, data.OnChange != nil, "changed", selectinput_onChanged, w)
	w.onFocus.Set(&cbt.Widget, data.OnFocus)
	w.onBlur.Set(&cbt.Widget, data.OnBlur)
	return nil
}
