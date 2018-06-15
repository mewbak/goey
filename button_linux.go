package goey

import (
	"unsafe"

	"bitbucket.org/rj/goey/syscall"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type mountedButton struct {
	handle *gtk.Button

	onClick clickSlot
	onFocus focusSlot
	onBlur  blurSlot
}

func (w *Button) mount(parent Control) (Element, error) {
	control, err := gtk.ButtonNewWithLabel(w.Text)
	if err != nil {
		return nil, err
	}
	control.AddEvents(int(gdk.FOCUS_CHANGE_MASK))

	(*gtk.Container)(unsafe.Pointer(parent.handle)).Add(control)
	control.SetSensitive(!w.Disabled)
	control.SetCanDefault(true)
	if w.Default {
		control.GrabDefault()
	}

	retval := &mountedButton{
		handle: control,
	}

	control.Connect("destroy", button_onDestroy, retval)
	retval.onClick.Set(&control.Widget, w.OnClick)
	retval.onFocus.Set(&control.Widget, w.OnFocus)
	retval.onBlur.Set(&control.Widget, w.OnBlur)
	control.Show()

	return retval, nil
}

func button_onDestroy(widget *gtk.Button, mounted *mountedButton) {
	mounted.handle = nil
}

func (w *mountedButton) Close() {
	if w.handle != nil {
		w.handle.Destroy()
		w.handle = nil
	}
}

func (w *mountedButton) Props() Widget {
	label, err := w.handle.GetChild()
	if err != nil {
		panic("Could not get child: " + err.Error())
	}
	text, err := (*gtk.Label)(unsafe.Pointer(label)).GetText()
	if err != nil {
		panic("Could not get text: " + err.Error())
	}

	return &Button{
		Text:     text,
		Disabled: !w.handle.GetSensitive(),
		Default:  w.handle.HasDefault(),
		OnClick:  w.onClick.callback,
		OnFocus:  w.onFocus.callback,
		OnBlur:   w.onBlur.callback,
	}
}

func (w *mountedButton) Handle() *gtk.Widget {
	return &w.handle.Widget
}

func (w *mountedButton) MeasureWidth() (Length, Length) {
	min, max := w.handle.GetPreferredWidth()
	return FromPixelsX(min), FromPixelsY(max)
}

func (w *mountedButton) MeasureHeight(width Length) (Length, Length) {
	min, max := syscall.WidgetGetPreferredHeightForWidth(&w.handle.Widget, width.PixelsX())
	return FromPixelsY(min), FromPixelsY(max)
}

func (w *mountedButton) SetBounds(bounds Rectangle) {
	pixels := bounds.Pixels()
	syscall.SetBounds(&w.handle.Widget, pixels.Min.X, pixels.Min.Y, pixels.Dx(), pixels.Dy())
}

func (w *mountedButton) updateProps(data *Button) error {
	label, err := w.handle.GetChild()
	if err != nil {
		return err
	}

	(*gtk.Label)(unsafe.Pointer(label)).SetText(data.Text)
	w.handle.SetSensitive(!data.Disabled)

	if data.Default {
		w.handle.GrabDefault()
	}
	w.onClick.Set(&w.handle.Widget, data.OnClick)
	w.onFocus.Set(&w.handle.Widget, data.OnFocus)
	w.onBlur.Set(&w.handle.Widget, data.OnBlur)

	return nil
}
