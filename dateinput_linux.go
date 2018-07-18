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
	control, err := gtk.CalendarNew()
	if err != nil {
		return nil, err
	}
	(*gtk.Container)(unsafe.Pointer(parent.Handle)).Add(control)
	control.SelectMonth(uint(w.Value.Month())-1, uint(w.Value.Year()))
	control.SelectDay(uint(w.Value.Day()))

	retval := &dateinputElement{
		Control:  Control{&control.Widget},
		onChange: w.OnChange,
	}

	control.Connect("destroy", dateinputOnDestroy, retval)
	retval.shChange = setSignalHandler(&control.Widget, 0, retval.onChange != nil, "day-selected", dateinputOnChanged, retval)
	retval.onFocus.Set(&control.Widget, w.OnFocus)
	retval.onBlur.Set(&control.Widget, w.OnBlur)
	control.Show()

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
func (w *dateinputElement) updateProps(data *DateInput) error {
	handle := w.calendar()

	w.onChange = nil // temporarily break OnChange to prevent event
	handle.SelectMonth(uint(data.Value.Month())-1, uint(data.Value.Year()))
	handle.SelectDay(uint(data.Value.Day()))
	w.onChange = data.OnChange
	//w.shChange = setSignalHandler(&w.handle.Widget, w.shChange, data.OnChange != nil, "value-changed", intinputOnChanged, w)
	w.onFocus.Set(&handle.Widget, data.OnFocus)
	w.onBlur.Set(&handle.Widget, data.OnBlur)

	return nil
}
