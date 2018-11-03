package goey

import (
	"time"
	"unsafe"

	"bitbucket.org/rj/goey/base"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type dateinputElement struct {
	Control

	onChange func(time.Time)
	shChange glib.SignalHandle
	onFocus  focusSlot
	onBlur   blurSlot
}

func (w *DateInput) mount(parent base.Control) (base.Element, error) {
	// Create the control
	control, err := gtk.CalendarNew()
	if err != nil {
		return nil, err
	}
	parent.Handle.Add(control)

	// Update properties on the control
	control.SelectMonth(uint(w.Value.Month())-1, uint(w.Value.Year()))
	control.SelectDay(uint(w.Value.Day()))
	control.SetSensitive(!w.Disabled)
	control.Show()

	// Create the element
	retval := &dateinputElement{
		Control:  Control{&control.Widget},
		onChange: w.OnChange,
	}

	// Connect all callbacks for the events
	control.Connect("destroy", dateinputOnDestroy, retval)
	retval.shChange = setSignalHandler(&control.Widget, 0, retval.onChange != nil, "day-selected", dateinputOnChanged, retval)
	retval.onFocus.Set(&control.Widget, w.OnFocus)
	retval.onBlur.Set(&control.Widget, w.OnBlur)

	return retval, nil
}

func dateinputOnChanged(widget *gtk.Calendar, mounted *dateinputElement) {
	if mounted.onChange == nil {
		return
	}

	y, m, d := widget.GetDate()
	mounted.onChange(time.Date(int(y), time.Month(m+1), int(d), 0, 0, 0, 0, time.Local))
}

func dateinputOnDestroy(widget *gtk.Calendar, mounted *dateinputElement) {
	mounted.handle = nil
}

func (w *dateinputElement) calendar() *gtk.Calendar {
	return (*gtk.Calendar)(unsafe.Pointer(w.handle))
}

func (w *dateinputElement) Props() base.Widget {
	control := w.calendar()
	year, month, day := control.GetDate()

	return &DateInput{
		Value:    time.Date(int(year), time.Month(month+1), int(day), 0, 0, 0, 0, time.Local),
		Disabled: !control.GetSensitive(),
		OnChange: w.onChange,
		OnFocus:  w.onFocus.callback,
		OnBlur:   w.onBlur.callback,
	}

}

func (w *dateinputElement) updateProps(data *DateInput) error {
	handle := w.calendar()

	w.onChange = nil // temporarily break OnChange to prevent event
	handle.SelectMonth(uint(data.Value.Month())-1, uint(data.Value.Year()))
	handle.SelectDay(uint(data.Value.Day()))
	handle.SetSensitive(!data.Disabled)
	w.onChange = data.OnChange
	w.shChange = setSignalHandler(w.handle, w.shChange, data.OnChange != nil, "day-selected", dateinputOnChanged, w)
	w.onFocus.Set(w.handle, data.OnFocus)
	w.onBlur.Set(w.handle, data.OnBlur)

	return nil
}
