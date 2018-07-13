package goey

import (
	"unsafe"

	"github.com/gotk3/gotk3/gtk"
)

type progressElement struct {
	Control
	min, max int
}

func (w *Progress) mount(parent Control) (Element, error) {
	control, err := gtk.ProgressBarNew()
	if err != nil {
		return nil, err
	}

	(*gtk.Container)(unsafe.Pointer(parent.handle)).Add(control)
	control.SetFraction(float64(w.Value-w.Min) / float64(w.Max-w.Min))

	retval := &progressElement{
		Control: Control{&control.Widget},
		min:     w.Min,
		max:     w.Max,
	}

	control.Connect("destroy", progressOnDestroy, retval)
	control.Show()

	return retval, nil
}

func progressOnDestroy(widget *gtk.ProgressBar, mounted *progressElement) {
	mounted.handle = nil
}

func (w *progressElement) progressbar() *gtk.ProgressBar {
	return (*gtk.ProgressBar)(unsafe.Pointer(w.handle))
}

func (w *progressElement) Props() Widget {
	if w.min == w.max {
		return &Progress{
			Value: w.min,
			Min:   w.min,
			Max:   w.max,
		}
	}

	pb := w.progressbar()
	value := pb.GetFraction()
	return &Progress{
		Value: w.min + int(float64(w.max-w.min)*value),
		Min:   w.min,
		Max:   w.max,
	}
}

func (w *progressElement) updateProps(data *Progress) error {
	pb := w.progressbar()
	w.min = data.Min
	w.max = data.Max
	pb.SetFraction(float64(data.Value-data.Min) / float64(data.Max-data.Min))
	return nil
}
