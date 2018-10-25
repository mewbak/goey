// +build !gnustep

package goey

import (
	"time"
	"bitbucket.org/rj/goey/base"
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

// Close removes the element from the GUI, and frees any associated resources.
func (w *Control) Close() {
	if w.handle != nil {
		w.handle.Destroy()
		w.handle = nil
	}
}

// Handle returns the platform-native handle for the control.
func (w *Control) Handle() *gtk.Widget {
	return w.handle
}

// TakeFocus is a wrapper around GrabFocus.
func (w *Control) TakeFocus() bool {
	// Check that the control can grab focus
	if !w.handle.GetCanFocus() {
		return false
	}

	w.handle.GrabFocus()
	// Note sure why the call to sleep is required, but there may be a debounce
	// provided by the system.  Without this call to sleep, the controls never
	// get the focus events.
	time.Sleep(250 * time.Millisecond)
	return w.handle.IsFocus()
}

// Layout determines the best size for an element that satisfies the
// constraints.
func (w *Control) Layout(bc base.Constraints) base.Size {
	if !bc.HasBoundedWidth() && !bc.HasBoundedHeight() {
		// No need to worry about breaking the constraints.  We can take as
		// much space as desired.
		_, width := w.handle.GetPreferredWidth()
		_, height := w.handle.GetPreferredHeight()
		// Dimensions may need to be increased to meet minimums.
		return bc.Constrain(base.Size{base.FromPixelsX(width), base.FromPixelsY(height)})
	}
	if !bc.HasBoundedHeight() {
		// No need to worry about height.  Find the width that best meets the
		// widgets preferred width.
		_, width1 := w.handle.GetPreferredWidth()
		width := bc.ConstrainWidth(base.FromPixelsX(width1))
		// Get the best height for this width.
		_, height := syscall.WidgetGetPreferredHeightForWidth(w.handle, width.PixelsX())
		// Height may need to be increased to meet minimum.
		return base.Size{width, bc.ConstrainHeight(base.FromPixelsY(height))}
	}

	// Not clear the following is the best general approach given GTK layout
	// model.
	height1, height2 := w.handle.GetPreferredHeight()
	if height := base.FromPixelsY(height2); height < bc.Max.Height {
		_, width := w.handle.GetPreferredWidth()
		return bc.Constrain(base.Size{base.FromPixelsX(width), height})
	}

	_, width := w.handle.GetPreferredWidth()
	return bc.Constrain(base.Size{base.FromPixelsX(width), base.FromPixelsX(height1)})
}

// MinIntrinsicHeight returns the minimum height that this element requires
// to be correctly displayed.
func (w *Control) MinIntrinsicHeight(width base.Length) base.Length {
	if width != base.Inf {
		height, _ := syscall.WidgetGetPreferredHeightForWidth(w.handle, width.PixelsX())
		return base.FromPixelsY(height)
	}
	height, _ := w.handle.GetPreferredHeight()
	return base.FromPixelsY(height)
}

// MinIntrinsicWidth returns the minimum width that this element requires
// to be correctly displayed.
func (w *Control) MinIntrinsicWidth(base.Length) base.Length {
	width, _ := w.handle.GetPreferredWidth()
	return base.FromPixelsX(width)
}

// SetBounds updates the position of the widget.
func (w *Control) SetBounds(bounds base.Rectangle) {
	pixels := bounds.Pixels()
	if pixels.Dx() <= 0 || pixels.Dy() <= 0 {
		panic("internal error.  zero width or zero height bounds for control")
	}
	syscall.SetBounds(w.handle, pixels.Min.X, pixels.Min.Y, pixels.Dx(), pixels.Dy())
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
