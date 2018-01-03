package goey

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type NativeWidget struct {
	handle *gtk.Widget
}

func (w *NativeWidget) Handle() *gtk.Widget {
	return w.handle
}

func (w *NativeWidget) Close() {
	if w.handle != nil {
		w.handle.Destroy()
		w.handle = nil
	}
}

type NativeMountedWidget interface {
	Handle() *gtk.Widget
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
	return true
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
	return true
}
