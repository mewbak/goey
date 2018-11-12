package goey

import (
	"unsafe"

	"bitbucket.org/rj/goey/base"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type selectinputElement struct {
	Control

	onChange func(int)
	shChange glib.SignalHandle
	onFocus  focusSlot
	onBlur   blurSlot
}

func (w *SelectInput) mount(parent base.Control) (base.Element, error) {
	control, err := gtk.ComboBoxTextNew()
	if err != nil {
		return nil, err
	}
	parent.Handle.Add(control)
	for _, v := range w.Items {
		control.AppendText(v)
	}
	control.SetActive(w.Value)
	control.SetCanFocus(true)
	control.SetSensitive(!w.Disabled)

	retval := &selectinputElement{
		Control:  Control{&control.Widget},
		onChange: w.OnChange,
	}

	control.Connect("destroy", selectinputOnDestroy, retval)
	retval.shChange = setSignalHandler(&control.Widget, 0, w.OnChange != nil, "changed", selectinputOnChanged, retval)
	retval.onFocus.Set(&control.Widget, w.OnFocus)
	retval.onBlur.Set(&control.Widget, w.OnBlur)
	control.Show()

	return retval, nil
}

func selectinputOnChanged(widget *gtk.ComboBoxText, mounted *selectinputElement) {
	if mounted.onChange == nil {
		return
	}

	mounted.onChange(widget.GetActive())
}

func selectinputOnDestroy(widget *gtk.ComboBoxText, mounted *selectinputElement) {
	mounted.handle = nil
}

func (w *selectinputElement) comboboxtext() *gtk.ComboBoxText {
	return (*gtk.ComboBoxText)(unsafe.Pointer(w.handle))
}

func (w *selectinputElement) Props() base.Widget {

	value := w.comboboxtext().GetActive()
	unset := value < 0

	return &SelectInput{
		Items:    nil,
		Value:    int(value),
		Unset:    unset,
		Disabled: !w.comboboxtext().GetSensitive(),
		OnChange: w.onChange,
		OnFocus:  w.onFocus.callback,
		OnBlur:   w.onBlur.callback,
	}

}

func (w *selectinputElement) updateProps(data *SelectInput) error {
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
	w.shChange = setSignalHandler(&cbt.Widget, w.shChange, data.OnChange != nil, "changed", selectinputOnChanged, w)
	w.onFocus.Set(&cbt.Widget, data.OnFocus)
	w.onBlur.Set(&cbt.Widget, data.OnBlur)
	return nil
}
