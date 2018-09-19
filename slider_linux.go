package goey

import (
	"unsafe"

	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/syscall"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type sliderElement struct {
	Control
	value    float64
	min, max float64

	onChange func(float64)
	shChange glib.SignalHandle
	onFocus  focusSlot
	onBlur   blurSlot
}

func (w *Slider) mount(parent base.Control) (base.Element, error) {
	control, err := gtk.ScaleNewWithRange(gtk.ORIENTATION_HORIZONTAL, w.Min, w.Max, (w.Max-w.Min)/10)
	if err != nil {
		return nil, err
	}

	parent.Handle.Add(control)
	control.SetValue(w.Value)
	control.SetDrawValue(false)
	control.SetSensitive(!w.Disabled)

	retval := &sliderElement{
		Control:  Control{&control.Widget},
		value:    w.Value,
		min:      w.Min,
		max:      w.Max,
		onChange: w.OnChange,
	}

	control.Connect("destroy", sliderOnDestroy, retval)
	retval.shChange = setSignalHandler(&control.Widget, 0, retval.onChange != nil, "change-value", sliderOnChangeValue, retval)
	retval.onFocus.Set(&control.Widget, w.OnFocus)
	retval.onBlur.Set(&control.Widget, w.OnBlur)
	control.Show()

	return retval, nil
}

func sliderOnDestroy(widget *gtk.Scale, mounted *sliderElement) {
	mounted.handle = nil
}

func sliderOnChangeValue(widget *gtk.Scale, scroll int, value float64, mounted *sliderElement) bool {
	if value < mounted.min {
		value = mounted.min
	} else if value > mounted.max {
		value = mounted.max
	}
	if value != mounted.value {
		mounted.value = value
		widget.SetValue(value)
		mounted.onChange(value)
		widget.QueueDraw()
	}
	return true
}

func (w *sliderElement) Props() base.Widget {
	pb := w.scale()

	return &Slider{
		Value:    pb.GetValue(),
		Min:      w.min,
		Max:      w.max,
		Disabled: !pb.GetSensitive(),
		OnFocus:  w.onFocus.callback,
		OnBlur:   w.onBlur.callback,
	}
}

// Layout determines the best size for an element that satisfies the
// constraints.
func (w *sliderElement) Layout(bc base.Constraints) base.Size {
	if !bc.HasBoundedWidth() && !bc.HasBoundedHeight() {
		// No need to worry about breaking the constraints.  We can take as
		// much space as desired.
		width := w.MinIntrinsicWidth(base.Inf)
		_, height := w.handle.GetPreferredHeight()
		// Dimensions may need to be increased to meet minimums.
		return bc.Constrain(base.Size{width, base.FromPixelsY(height)})
	}
	if !bc.HasBoundedHeight() {
		// No need to worry about height.  Find the width that best meets the
		// widgets preferred width.
		width := bc.ConstrainWidth(w.MinIntrinsicWidth(base.Inf))
		// Get the best height for this width.
		_, height := syscall.WidgetGetPreferredHeightForWidth(w.handle, width.PixelsX())
		// Height may need to be increased to meet minimum.
		return base.Size{width, bc.ConstrainHeight(base.FromPixelsY(height))}
	}

	// Not clear the following is the best general approach given GTK layout
	// model.
	_, height2 := w.handle.GetPreferredHeight()
	if height := base.FromPixelsY(height2); height < bc.Max.Height {
		width := w.MinIntrinsicWidth(height)
		return bc.Constrain(base.Size{width, height})
	}

	height := base.FromPixelsY(height2)
	width := w.MinIntrinsicWidth(height)
	return bc.Constrain(base.Size{width, height})
}

// MinIntrinsicWidth returns the minimum width that this element requires
// to be correctly displayed.
func (w *sliderElement) MinIntrinsicWidth(base.Length) base.Length {
	width, _ := w.handle.GetPreferredWidth()
	if limit := base.FromPixelsX(width); limit < 160*DIP {
		return 160 * DIP
	}
	return base.FromPixelsX(width)
}

func (w *sliderElement) scale() *gtk.Scale {
	return (*gtk.Scale)(unsafe.Pointer(w.handle))
}

func (w *sliderElement) updateProps(data *Slider) error {
	pb := w.scale()
	w.value = data.Value
	w.min = data.Min
	w.max = data.Max
	pb.SetRange(data.Min, data.Max)
	pb.SetValue(data.Value)
	pb.SetSensitive(!data.Disabled)
	w.onChange = data.OnChange
	w.shChange = setSignalHandler(w.handle, w.shChange, w.onChange != nil, "change-value", sliderOnChangeValue, w)
	w.onFocus.Set(w.handle, data.OnFocus)
	w.onBlur.Set(w.handle, data.OnBlur)
	return nil
}
