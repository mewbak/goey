package goey

import (
	"image/color"
	"unsafe"

	"bitbucket.org/rj/goey/syscall"
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gtk"
)

func (w *Decoration) mount(parent Control) (Element, error) {
	control, err := gtk.DrawingAreaNew()
	if err != nil {
		return nil, err
	}
	(*gtk.Container)(unsafe.Pointer(parent.handle)).Add(control)

	retval := &mountedDecoration{
		handle:  control,
		fill:    w.Fill,
		padding: w.Padding,
	}

	control.Connect("destroy", decoration_onDestroy, retval)
	control.Connect("draw", decoration_onDraw, retval)
	control.Show()

	child, err := DiffChild(parent, nil, w.Child)
	if err != nil {
		control.Destroy()
		return nil, err
	}
	retval.child = child

	return retval, nil
}

type mountedDecoration struct {
	handle  *gtk.DrawingArea
	fill    color.RGBA
	padding Padding
	radius  Length
	child   Element
}

func decoration_onDestroy(widget *gtk.DrawingArea, mounted *mountedDecoration) {
	mounted.handle = nil
}

func decoration_onDraw(widget *gtk.DrawingArea, cr *cairo.Context, mounted *mountedDecoration) bool {
	a := mounted.handle.GetAllocation()
	if mounted.radius > 0 {
		panic("")
	} else {
		cr.Rectangle(0, 0, float64(a.GetWidth()), float64(a.GetHeight()))
		cr.SetSourceRGB(float64(mounted.fill.R)/0xFF, float64(mounted.fill.G)/0xFF, float64(mounted.fill.B)/0xFF)
		cr.Fill()
	}
	return false
}

func (w *mountedDecoration) Close() {
	if w.handle != nil {
		w.child.Close()
		w.child = nil
		w.handle.Destroy()
		w.handle = nil
	}
}

func (w *mountedDecoration) MeasureWidth() (Length, Length) {
	if w.child != nil {
		px := FromPixelsX(1)
		min, max := w.child.MeasureWidth()
		min += 2*px + w.padding.Left + w.padding.Right
		max += 2*px + w.padding.Left + w.padding.Right
		return min, max
	}

	return 13 * DIP, 13 * DIP
}

func (w *mountedDecoration) MeasureHeight(width Length) (Length, Length) {
	if w.child != nil {
		px := FromPixelsX(1)
		py := FromPixelsY(1)
		min, max := w.child.MeasureHeight(width - 2*px)
		min += 2*py + w.padding.Top + w.padding.Bottom
		max += 2*py + w.padding.Top + w.padding.Bottom
		return min, max
	}

	return 13 * DIP, 13 * DIP
}

func (w *mountedDecoration) SetBounds(bounds Rectangle) {
	pixels := bounds.Pixels()
	syscall.SetBounds(&w.handle.Widget, pixels.Min.X, pixels.Min.Y, pixels.Dx(), pixels.Dy())

	bounds.Min.X += w.padding.Left
	bounds.Min.Y += w.padding.Top
	bounds.Max.X -= w.padding.Right
	bounds.Max.Y -= w.padding.Bottom
	w.child.SetBounds(bounds)
}

func (w *mountedDecoration) updateProps(data *Decoration) error {
	w.fill = data.Fill

	parent, err := w.handle.GetParent()
	if err != nil {
		return err
	}
	w.child, err = DiffChild(Control{parent}, w.child, data.Child)
	if err != nil {
		return err
	}

	return nil
}
