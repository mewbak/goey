package goey

import (
	"unsafe"

	"bitbucket.org/rj/goey/syscall"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type mountedSelectInput struct {
	handle *gtk.ComboBoxText

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
		handle:   control,
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
	mounted.onChange(widget.GetActive())
}

func selectinput_onDestroy(widget *gtk.ComboBoxText, mounted *mountedSelectInput) {
	mounted.handle = nil
}

func (w *mountedSelectInput) Close() {
	if w.handle != nil {
		w.handle.Destroy()
		w.handle = nil
	}
}

func (w *mountedSelectInput) Handle() *gtk.Widget {
	return &w.handle.Widget
}

func (w *mountedSelectInput) MeasureWidth() (Length, Length) {
	min, max := w.handle.GetPreferredWidth()
	return FromPixelsX(min), FromPixelsY(max)
}

func (w *mountedSelectInput) MeasureHeight(width Length) (Length, Length) {
	min, max := syscall.WidgetGetPreferredHeightForWidth(&w.handle.Widget, width.PixelsX())
	return FromPixelsY(min), FromPixelsY(max)
}

func (w *mountedSelectInput) SetBounds(bounds Rectangle) {
	pixels := bounds.Pixels()
	syscall.SetBounds(&w.handle.Widget, pixels.Min.X, pixels.Min.Y, pixels.Dx(), pixels.Dy())
}

func (w *mountedSelectInput) updateProps(data *SelectInput) error {
	// Todo, can we avoid rebuilding the list?
	w.handle.RemoveAll()
	for _, v := range data.Items {
		w.handle.AppendText(v)
	}
	w.handle.SetActive(data.Value)

	w.handle.SetSensitive(!data.Disabled)
	w.onChange = data.OnChange
	w.shChange = setSignalHandler(&w.handle.Widget, w.shChange, data.OnChange != nil, "changed", selectinput_onChanged, w)
	w.onFocus.Set(&w.handle.Widget, data.OnFocus)
	w.onBlur.Set(&w.handle.Widget, data.OnBlur)
	return nil
}
