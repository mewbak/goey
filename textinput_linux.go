package goey

import (
	"unsafe"

	"bitbucket.org/rj/goey/syscall"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type mountedTextInput struct {
	handle *gtk.Entry

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
	if w.Password {
		control.SetVisibility(false)
	}
	if w.ReadOnly {
		control.SetEditable(false)
	}

	retval := &mountedTextInput{
		handle:     control,
		onChange:   w.OnChange,
		onEnterKey: w.OnEnterKey,
	}

	control.Connect("destroy", textinput_onDestroy, retval)
	retval.shChange = setSignalHandler(&control.Widget, 0, retval.onChange != nil, "changed", textinput_onChanged, retval)
	retval.onFocus.Set(&control.Widget, w.OnFocus)
	retval.onBlur.Set(&control.Widget, w.OnBlur)
	retval.shEnterKey = setSignalHandler(&control.Widget, 0, retval.onEnterKey != nil, "activate", textinput_onActivate, retval)
	control.Show()

	return retval, nil
}

func textinput_onActivate(obj *glib.Object, mounted *mountedTextInput) {
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

func textinput_onChanged(widget *gtk.Entry, mounted *mountedTextInput) {
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

func textinput_onDestroy(widget *gtk.Entry, mounted *mountedTextInput) {
	mounted.handle = nil
}

func (w *mountedTextInput) Close() {
	if w.handle != nil {
		w.handle.Destroy()
		w.handle = nil
	}
}

func (w *mountedTextInput) Handle() *gtk.Widget {
	return &w.handle.Widget
}

func (w *mountedTextInput) MeasureWidth() (Length, Length) {
	min, max := w.handle.GetPreferredWidth()
	return FromPixelsX(min), FromPixelsX(max)
}

func (w *mountedTextInput) MeasureHeight(width Length) (Length, Length) {
	min, max := syscall.WidgetGetPreferredHeightForWidth(&w.handle.Widget, width.PixelsX())
	return FromPixelsY(min), FromPixelsY(max)
}

func (w *mountedTextInput) SetBounds(bounds Rectangle) {
	pixels := bounds.Pixels()
	syscall.SetBounds(&w.handle.Widget, pixels.Min.X, pixels.Min.Y, pixels.Dx(), pixels.Dy())
}

func (w *mountedTextInput) updateProps(data *TextInput) error {
	w.onChange = nil // temporarily break OnChange to prevent event
	w.handle.SetText(data.Value)
	w.handle.SetEditable(!data.ReadOnly)
	w.handle.SetPlaceholderText(data.Placeholder)
	w.handle.SetVisibility(!data.Password)
	w.onChange = data.OnChange
	w.shChange = setSignalHandler(&w.handle.Widget, w.shChange, data.OnChange != nil, "changed", textinput_onChanged, w)
	w.onFocus.Set(&w.handle.Widget, data.OnFocus)
	w.onBlur.Set(&w.handle.Widget, data.OnBlur)
	w.onEnterKey = data.OnEnterKey
	w.shEnterKey = setSignalHandler(&w.handle.Widget, w.shEnterKey, data.OnEnterKey != nil, "activate", textinput_onActivate, w)

	return nil
}
