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
	if !w.Unset {
		control.SetActive(w.Value)
	}
	control.SetCanFocus(true)
	control.SetSensitive(!w.Disabled)

	retval := &selectinputElement{
		Control:  Control{&control.Widget},
		onChange: w.OnChange,
	}

	control.Connect("destroy", selectinputOnDestroy, retval)
	retval.shChange = setSignalHandler(&control.Widget, 0, w.OnChange != nil, "changed", selectinputOnChanged, retval)
	if child, err := control.GetChild(); err == nil {
		child.SetCanFocus(true)
		retval.onFocus.Set(child, w.OnFocus)
		retval.onBlur.Set(child, w.OnBlur)
	}
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
	if unset {
		value = 0
	}

	return &SelectInput{
		Items:    w.propsItems(),
		Value:    int(value),
		Unset:    unset,
		Disabled: !w.comboboxtext().GetSensitive(),
		OnChange: w.onChange,
		OnFocus:  w.onFocus.callback,
		OnBlur:   w.onBlur.callback,
	}

}

func (w *selectinputElement) propsItems() []string {
	// Get the model for the combobox, which contains the list of items.
	model, err := w.comboboxtext().GetModel()
	if err != nil {
		return nil
	}

	// Iterate through the list.  The model can in principle hold a tree, but
	// that won't occur within the combobox.
	items := []string{}
	for iter, ok := model.GetIterFirst(); ok; ok = model.IterNext(iter) {
		v, err := model.GetValue(iter, 0)
		if err != nil {
			return nil
		}
		vs, err := v.GetString()
		if err != nil {
			return nil
		}
		items = append(items, vs)
	}

	return items
}

func (w *selectinputElement) TakeFocus() bool {
	widget, err := w.comboboxtext().GetChild()
	if err != nil {
		return false
	}

	control := Control{widget}
	return control.TakeFocus()
}

func (w *selectinputElement) updateProps(data *SelectInput) error {
	cbt := w.comboboxtext()

	w.onChange = nil // temporarily break OnChange to prevent event
	// Todo, can we avoid rebuilding the list?
	cbt.RemoveAll()
	for _, v := range data.Items {
		cbt.AppendText(v)
	}
	if !data.Unset {
		cbt.SetActive(data.Value)
	}

	cbt.SetSensitive(!data.Disabled)
	w.onChange = data.OnChange
	w.shChange = setSignalHandler(&cbt.Widget, w.shChange, data.OnChange != nil, "changed", selectinputOnChanged, w)
	if child, err := cbt.GetChild(); err == nil {
		w.onFocus.Set(child, data.OnFocus)
		w.onBlur.Set(child, data.OnBlur)
	}
	return nil
}
