package goey

import (
	"time"
	"unsafe"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type mountedDateInput struct {
	handle *gtk.Calendar

	onChange func(time.Time)
	shChange glib.SignalHandle
	onFocus  focusSlot
	onBlur   blurSlot
}

func (w *DateInput) mount(parent NativeWidget) (MountedWidget, error) {
	control, err := gtk.CalendarNew()
	if err != nil {
		return nil, err
	}
	(*gtk.Container)(unsafe.Pointer(parent.handle)).Add(control)
	control.SelectMonth(uint(w.Value.Month())-1, uint(w.Value.Year()))
	control.SelectDay(uint(w.Value.Day()))

	retval := &mountedDateInput{
		handle:   control,
		onChange: w.OnChange,
	}

	control.Connect("destroy", dateinput_onDestroy, retval)
	retval.shChange = setSignalHandler(&control.Widget, 0, retval.onChange != nil, "day-selected", dateinput_onChanged, retval)
	retval.onFocus.Set(&control.Widget, w.OnFocus)
	retval.onBlur.Set(&control.Widget, w.OnBlur)
	control.Show()

	return retval, nil
}

func dateinput_onChanged(widget *gtk.Calendar, mounted *mountedDateInput) {
	y, m, d := widget.GetDate()
	mounted.onChange(time.Date(int(y), time.Month(m+1), int(d), 0, 0, 0, 0, time.Local))
}

func dateinput_onDestroy(widget *gtk.Calendar, mounted *mountedDateInput) {
	mounted.handle = nil
}

func (w *mountedDateInput) Close() {
	if w.handle != nil {
		w.handle.Destroy()
		w.handle = nil
	}
}

func (w *mountedDateInput) Handle() *gtk.Widget {
	return &w.handle.Widget
}

func (w *mountedDateInput) updateProps(data *DateInput) error {
	w.handle.SelectMonth(uint(data.Value.Month())-1, uint(data.Value.Year()))
	w.handle.SelectDay(uint(data.Value.Day()))
	w.onChange = data.OnChange
	w.shChange = setSignalHandler(&w.handle.Widget, w.shChange, data.OnChange != nil, "value-changed", intinput_onChanged, w)
	w.onFocus.Set(&w.handle.Widget, data.OnFocus)
	w.onBlur.Set(&w.handle.Widget, data.OnBlur)

	return nil
}
