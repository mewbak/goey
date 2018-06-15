package goey

import (
	"unsafe"

	"bitbucket.org/rj/goey/syscall"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type mountedCheckbox struct {
	handle *gtk.CheckButton

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
		handle:   control,
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
	mounted.onChange(widget.GetActive())
}

func checkbox_onDestroy(widget *gtk.CheckButton, mounted *mountedCheckbox) {
	mounted.handle = nil
}

func (w *mountedCheckbox) Close() {
	if w.handle != nil {
		w.handle.Destroy()
		w.handle = nil
	}
}

func (w *mountedCheckbox) Props() Widget {
	label, err := w.handle.GetChild()
	if err != nil {
		panic("Could not get child: " + err.Error())
	}
	text, err := (*gtk.Label)(unsafe.Pointer(label)).GetText()
	if err != nil {
		panic("Could not get text: " + err.Error())
	}

	return &Checkbox{
		Value:    w.handle.GetActive(),
		Text:     text,
		Disabled: !w.handle.GetSensitive(),
		OnChange: w.onChange,
		OnFocus:  w.onFocus.callback,
		OnBlur:   w.onBlur.callback,
	}
}

func (w *mountedCheckbox) Handle() *gtk.Widget {
	return &w.handle.Widget
}

func (w *mountedCheckbox) MeasureWidth() (Length, Length) {
	min, max := w.handle.GetPreferredWidth()
	return FromPixelsX(min), FromPixelsY(max)
}

func (w *mountedCheckbox) MeasureHeight(width Length) (Length, Length) {
	min, max := syscall.WidgetGetPreferredHeightForWidth(&w.handle.Widget, width.PixelsX())
	return FromPixelsY(min), FromPixelsY(max)
}

func (w *mountedCheckbox) SetBounds(bounds Rectangle) {
	pixels := bounds.Pixels()
	syscall.SetBounds(&w.handle.Widget, pixels.Min.X, pixels.Min.Y, pixels.Dx(), pixels.Dy())
}

func (w *mountedCheckbox) updateProps(data *Checkbox) error {
	label, err := w.handle.GetChild()
	if err != nil {
		return err
	}

	(*gtk.Label)(unsafe.Pointer(label)).SetText(data.Text)
	w.handle.SetActive(data.Value)
	w.handle.SetSensitive(!data.Disabled)

	w.onChange = data.OnChange
	w.shClick = setSignalHandler(&w.handle.Widget, w.shClick, data.OnChange != nil, "clicked", checkbox_onClick, w)
	w.onFocus.Set(&w.handle.Widget, data.OnFocus)
	w.onBlur.Set(&w.handle.Widget, data.OnBlur)

	return nil
}
