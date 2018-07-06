package goey

import (
	"bitbucket.org/rj/goey/syscall"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

// Control is an opaque type used as a platform-specific handle to a control
// created using the platform GUI.  As an example, this will refer to a HWND
// when targeting Windows, but a *GtkWidget when targeting GTK.
//
// Unless developping new widgets, users should not need to use this type.
//
// Any method's on this type will be platform specific.
type Control struct {
	handle *gtk.Widget
}

func (w *Control) Close() {
	if w.handle != nil {
		w.handle.Destroy()
		w.handle = nil
	}
}

func (w *Control) Handle() *gtk.Widget {
	return w.handle
}

func (w *Control) Layout(bc Constraint) Size {
	if !bc.HasBoundedWidth() && !bc.HasBoundedHeight() {
		// No need to worry about breaking the constraints.  We can take as 
		// much space as desired.
		_, width := w.handle.GetPreferredWidth()
		_, height := w.handle.GetPreferredHeight()
		// Dimensions may need to be increased to meet minimums.
		return bc.Constrain(Size{FromPixelsX(width), FromPixelsY(height)})
	}
	if !bc.HasBoundedHeight() {
		// No need to worry about height.  Find the width that best meets the
		// widgets preferred width.
		_, width1 := w.handle.GetPreferredWidth()
		width := bc.ConstrainWidth(FromPixelsX(width1))
		// Get the best height for this width.
		_, height := syscall.WidgetGetPreferredHeightForWidth(w.handle, width.PixelsX())
		// Height may need to be increased to meet minimum.
		return Size{width, bc.ConstrainHeight(FromPixelsY(height)}))
	}

	// Not clear the following is the best general approach given GTK layout
	// model.  
	height1, height2 := w.handle.GetPreferredHeight()
	if height := FromPixelsY(height2); height < bc.Max.Height {
		_, width := w.handle.GetPreferredWidth()
		return bc.Constrain(Size{FromPixelsX(width), height})
	}

	_, width := w.handle.GetPreferredWidth()
	return bc.Constrain(Size{FromPixelsX(width), FromPixelsX(height1)})
}

func (w *Control) MinIntrinsicHeight(width Length) Length {
	if width != Inf {
		height, _ := syscall.WidgetGetPreferredHeightForWidth(w.handle, width.PixelsX())
		return FromPixelsY(height)
	}
	height, _ := w.handle.GetPreferredHeight()
	return FromPixelsY(height)
}

func (w *Control) MinIntrinsicWidth(Length) Length {
	width, _ := w.handle.GetPreferredWidth()
	return FromPixelsX(width)
}

func (w *Control) SetBounds(bounds Rectangle) {
	pixels := bounds.Pixels()
	if pixels.Dx() <= 0 || pixels.Dy() <= 0 {
		panic("internal error.  zero width or zero height bounds for control")
	}
	syscall.SetBounds(w.handle, pixels.Min.X, pixels.Min.Y, pixels.Dx(), pixels.Dy())
}

// NativeElement contains platform-specific methods that all widgets
// must support on GTK.
type NativeElement interface {
}

func setSignalHandler(control *gtk.Widget, sh glib.SignalHandle, ok bool, name string, thunk interface{}, userData interface{}) glib.SignalHandle {
	if ok && sh == 0 {
		sh, err := control.Connect(name, thunk, userData)
		if err != nil {
			panic("Failed to connect '" + name + "'.")
		}
		return sh
	} else if !ok && sh != 0 {
		control.HandlerDisconnect(sh)
		return 0
	}

	return sh
}

type clickSlot struct {
	callback func()
	handle   glib.SignalHandle
}

func (c *clickSlot) Set(control *gtk.Widget, value func()) {
	if value != nil && c.handle == 0 {
		handle, err := control.Connect("clicked", clickSlotThunk, c)
		if err != nil {
			panic("Failed to connect 'clicked'.")
		}
		c.handle = handle
	} else if value == nil && c.handle != 0 {
		control.HandlerDisconnect(c.handle)
		c.handle = 0
	}
	c.callback = value
}

func clickSlotThunk(widget interface{}, c *clickSlot) {
	c.callback()
}

type focusSlot struct {
	callback func()
	handle   glib.SignalHandle
}

func (c *focusSlot) Set(control *gtk.Widget, value func()) {
	if value != nil && c.handle == 0 {
		handle, err := control.Connect("focus-in-event", focusSlotThunk, c)
		if err != nil {
			panic("Failed to connect 'focus-in-event'.")
		}
		c.handle = handle
	} else if value == nil && c.handle != 0 {
		control.HandlerDisconnect(c.handle)
		c.handle = 0
	}
	c.callback = value
}

func focusSlotThunk(widget interface{}, event *gdk.Event, c *focusSlot) bool {
	c.callback()
	return false
}

type blurSlot struct {
	callback func()
	handle   glib.SignalHandle
}

func (c *blurSlot) Set(control *gtk.Widget, value func()) {
	if value != nil && c.handle == 0 {
		handle, err := control.Connect("focus-out-event", blurSlotThunk, c)
		if err != nil {
			panic("Failed to connect 'focus-out-event'.")
		}
		c.handle = handle
	} else if value == nil && c.handle != 0 {
		control.HandlerDisconnect(c.handle)
		c.handle = 0
	}
	c.callback = value
}

func blurSlotThunk(widget interface{}, event *gdk.Event, c *focusSlot) bool {
	c.callback()
	return false
}
