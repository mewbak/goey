package goey

import (
	"strconv"
	"unsafe"

	"bitbucket.org/rj/goey/base"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type intinputElement struct {
	Control

	onChange   func(int64)
	shChange   glib.SignalHandle
	onFocus    focusSlot
	onBlur     blurSlot
	onEnterKey func(int64)
	shEnterKey glib.SignalHandle
}

func (w *IntInput) mount(parent base.Control) (base.Element, error) {
	// Create the control
	control, err := gtk.SpinButtonNew(nil, 1, 0)
	if err != nil {
		return nil, err
	}
	parent.Handle.Add(control)

	// Update properties on the control
	control.SetRange(float64(w.Min), float64(w.Max))
	control.SetValue(float64(w.Value))
	control.SetIncrements(1, 10)
	control.SetPlaceholderText(w.Placeholder)
	control.SetSensitive(!w.Disabled)

	// Create the element
	retval := &intinputElement{
		Control:    Control{&control.Widget},
		onChange:   w.OnChange,
		onEnterKey: w.OnEnterKey,
	}

	// Connect all callbacks for the events
	control.Connect("destroy", intinputOnDestroy, retval)
	retval.shChange = setSignalHandler(&control.Widget, 0, retval.onChange != nil, "value-changed", intinputOnChanged, retval)
	retval.onFocus.Set(&control.Widget, w.OnFocus)
	retval.onBlur.Set(&control.Widget, w.OnBlur)
	retval.shEnterKey = setSignalHandler(&control.Widget, 0, retval.onEnterKey != nil, "activate", intinputOnActivate, retval)
	control.Show()

	return retval, nil
}

func intinputOnActivate(obj *glib.Object, mounted *intinputElement) {
	// Not sure why, but the widget comes into this callback as a glib.Object,
	// and not the gtk.SpinButton.  Need to wrap the value.  This pokes into the internals
	// of the gtk package.
	widget := gtk.SpinButton{gtk.Entry{gtk.Widget{glib.InitiallyUnowned{obj}}, gtk.Editable{obj}}}
	text, _ := widget.GetText()
	value, _ := strconv.ParseInt(text, 10, 64)
	// What should be done with a parsing error.  The control should prevent
	// that occurring.
	mounted.onEnterKey(value)
}

func intinputOnChanged(widget *gtk.SpinButton, mounted *intinputElement) {
	if mounted.onChange == nil {
		return
	}

	text, _ := widget.GetText()
	value, _ := strconv.ParseInt(text, 10, 64)
	// What should be done with a parsing error.  The control should prevent
	// that occurring.
	mounted.onChange(value)
}

func intinputOnDestroy(widget *gtk.SpinButton, mounted *intinputElement) {
	mounted.handle = nil
}

func (w *intinputElement) spinbutton() *gtk.SpinButton {
	return (*gtk.SpinButton)(unsafe.Pointer(w.handle))
}

// Because GTK uses double (or float64 in Go) to store the range, it cannot
// keep full precision for values near the minimum or maximum of the int64
// range.  This function is just for Props, and is used to adjust the int64
// to get a match.
func toInt64(scale float64) int64 {
	a := int64(scale)
	if float64(a) == scale {
		return a
	}
	if float64(a-1) == scale {
		return a - 1
	}
	println("mismatch", a, scale)
	return a
}

func (w *intinputElement) Props() base.Widget {
	button := w.spinbutton()

	placeholder, err := button.GetPlaceholderText()
	if err != nil {
		panic("Could not get placeholder text: " + err.Error())
	}

	return &IntInput{
		Value:       int64(button.GetValue()),
		Placeholder: placeholder,
		Disabled:    !button.GetSensitive(),
		Min:         toInt64(button.GetAdjustment().GetLower()),
		Max:         toInt64(button.GetAdjustment().GetUpper()),
		OnChange:    w.onChange,
		OnFocus:     w.onFocus.callback,
		OnBlur:      w.onBlur.callback,
		OnEnterKey:  w.onEnterKey,
	}
}

func (w *intinputElement) updateProps(data *IntInput) error {
	button := w.spinbutton()

	w.onChange = nil // break OnChange to prevent event
	button.SetRange(float64(data.Min), float64(data.Max))
	button.SetValue(float64(data.Value))
	button.SetPlaceholderText(data.Placeholder)
	button.SetSensitive(!data.Disabled)
	w.onChange = data.OnChange
	w.shChange = setSignalHandler(&button.Widget, w.shChange, data.OnChange != nil, "value-changed", intinputOnChanged, w)
	w.onFocus.Set(&button.Widget, data.OnFocus)
	w.onBlur.Set(&button.Widget, data.OnBlur)
	w.onEnterKey = data.OnEnterKey
	w.shEnterKey = setSignalHandler(&button.Widget, w.shEnterKey, data.OnEnterKey != nil, "activate", intinputOnActivate, w)

	return nil
}
