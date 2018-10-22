package goey

import (
	"unsafe"

	"bitbucket.org/rj/goey/base"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type checkboxElement struct {
	Control

	onChange func(bool)
	shClick  glib.SignalHandle
	onFocus  focusSlot
	onBlur   blurSlot
}

func (w *Checkbox) mount(parent base.Control) (base.Element, error) {
	control, err := gtk.CheckButtonNewWithLabel(w.Text)
	if err != nil {
		return nil, err
	}
	parent.Handle.Add(control)
	control.SetActive(w.Value)
	control.SetSensitive(!w.Disabled)

	retval := &checkboxElement{
		Control:  Control{&control.Widget},
		onChange: w.OnChange,
	}

	control.Connect("destroy", checkboxOnDestroy, retval)
	retval.shClick = setSignalHandler(&control.Widget, 0, w.OnChange != nil, "clicked", checkboxOnClick, retval)
	retval.onFocus.Set(&control.Widget, w.OnFocus)
	retval.onBlur.Set(&control.Widget, w.OnBlur)
	control.Show()

	return retval, nil
}

func checkboxOnClick(widget *gtk.CheckButton, mounted *checkboxElement) {
	if mounted.onChange == nil {
		return
	}

	mounted.onChange(widget.GetActive())
}

func checkboxOnDestroy(widget *gtk.CheckButton, mounted *checkboxElement) {
	mounted.handle = nil
}

func (w *checkboxElement) checkbutton() *gtk.CheckButton {
	return (*gtk.CheckButton)(unsafe.Pointer(w.handle))
}

func (w *checkboxElement) Click() {
	w.checkbutton().Clicked()
}

func (w *checkboxElement) Props() base.Widget {
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

func (w *checkboxElement) updateProps(data *Checkbox) error {
	checkbutton := w.checkbutton()

	w.onChange = nil // temporarily break OnChange to prevent event
	checkbutton.SetLabel(data.Text)
	checkbutton.SetActive(data.Value)
	checkbutton.SetSensitive(!data.Disabled)

	w.onChange = data.OnChange
	w.shClick = setSignalHandler(&checkbutton.Widget, w.shClick, data.OnChange != nil, "clicked", checkboxOnClick, w)
	w.onFocus.Set(&checkbutton.Widget, data.OnFocus)
	w.onBlur.Set(&checkbutton.Widget, data.OnBlur)

	return nil
}
