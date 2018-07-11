package goey

import (
	"unsafe"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type mountedTextInput struct {
	Control

	onChange   func(string)
	shChange   glib.SignalHandle
	onFocus    focusSlot
	onBlur     blurSlot
	onEnterKey func(string)
	shEnterKey glib.SignalHandle
}

func (w *TextInput) mount(parent Control) (Element, error) {
	control, err := gtk.EntryNew()
	if err != nil {
		return nil, err
	}
	(*gtk.Container)(unsafe.Pointer(parent.handle)).Add(control)
	control.SetText(w.Value)
	control.SetPlaceholderText(w.Placeholder)
	control.SetSensitive(!w.Disabled)
	control.SetVisibility(!w.Password)
	control.SetEditable(!w.ReadOnly)

	retval := &mountedTextInput{
		Control:    Control{&control.Widget},
		onChange:   w.OnChange,
		onEnterKey: w.OnEnterKey,
	}

	control.Connect("destroy", textinputOnDestroy, retval)
	retval.shChange = setSignalHandler(&control.Widget, 0, retval.onChange != nil, "changed", textinputOnChanged, retval)
	retval.onFocus.Set(&control.Widget, w.OnFocus)
	retval.onBlur.Set(&control.Widget, w.OnBlur)
	retval.shEnterKey = setSignalHandler(&control.Widget, 0, retval.onEnterKey != nil, "activate", textinputOnActivate, retval)
	control.Show()

	return retval, nil
}

func textinputOnActivate(obj *glib.Object, mounted *mountedTextInput) {
	// Not sure why, but the widget comes into this callback as a glib.Object,
	// and not the gtk.Entry.  Need to wrap the value.  This pokes into the internals
	// of the gtk package.
	widget := gtk.Entry{gtk.Widget{glib.InitiallyUnowned{obj}}, gtk.Editable{obj}}
	text, err := widget.GetText()
	if err != nil {
		// TODO:  What is the correct reporting here
		return
	}
	mounted.onEnterKey(text)
}

func textinputOnChanged(widget *gtk.Entry, mounted *mountedTextInput) {
	if mounted.onChange == nil {
		return
	}

	text, err := widget.GetText()
	if err != nil {
		// TODO:  What is the correct reporting here
		return
	}
	mounted.onChange(text)
}

func textinputOnDestroy(widget *gtk.Entry, mounted *mountedTextInput) {
	mounted.handle = nil
}

func (w *mountedTextInput) entry() *gtk.Entry {
	return (*gtk.Entry)(unsafe.Pointer(w.handle))
}

func (w *mountedTextInput) Props() Widget {
	entry := w.entry()
	value, err := entry.GetText()
	if err != nil {
		panic("could not get text, " + err.Error())
	}
	placeholder, err := entry.GetPlaceholderText()
	if err != nil {
		panic("could not get placeholder text, " + err.Error())
	}

	return &TextInput{
		Value:       value,
		Disabled:    !entry.GetSensitive(),
		Placeholder: placeholder,
		Password:    !entry.GetVisibility(),
		ReadOnly:    !entry.GetEditable(),
		OnChange:    w.onChange,
		OnFocus:     w.onFocus.callback,
		OnBlur:      w.onBlur.callback,
		OnEnterKey:  w.onEnterKey,
	}
}

func (w *mountedTextInput) updateProps(data *TextInput) error {
	entry := w.entry()
	w.onChange = nil // temporarily break OnChange to prevent event
	entry.SetText(data.Value)
	entry.SetEditable(!data.ReadOnly)
	entry.SetPlaceholderText(data.Placeholder)
	entry.SetSensitive(!data.Disabled)
	entry.SetVisibility(!data.Password)
	w.onChange = data.OnChange
	w.shChange = setSignalHandler(&entry.Widget, w.shChange, data.OnChange != nil, "changed", textinputOnChanged, w)
	w.onFocus.Set(&entry.Widget, data.OnFocus)
	w.onBlur.Set(&entry.Widget, data.OnBlur)
	w.onEnterKey = data.OnEnterKey
	w.shEnterKey = setSignalHandler(&entry.Widget, w.shEnterKey, data.OnEnterKey != nil, "activate", textinputOnActivate, w)

	return nil
}
