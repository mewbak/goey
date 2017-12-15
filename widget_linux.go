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

type clickSlot struct {
	callback func()
	handle   glib.SignalHandle
}

func (c *clickSlot) Set(control *gtk.Widget, value func()) error {
	if value != nil && c.handle == 0 {
		handle, err := control.Connect("clicked", clickSlotThunk, c)
		if err != nil {
			return err
		}
		c.handle = handle
	} else if value == nil && c.handle != 0 {
		control.HandlerDisconnect(c.handle)
		c.handle = 0
	}
	c.callback = value
	return nil
}

func (c *clickSlot) Close(control *gtk.Widget) {
	if c.handle != 0 {
		control.HandlerDisconnect(c.handle)
		c.handle = 0
	}
}

func clickSlotThunk(widget *gtk.Button, c *clickSlot) {
	c.callback()
}

type focusSlot struct {
	callback func()
	handle   glib.SignalHandle
}

func (c *focusSlot) Set(control *gtk.Widget, value func()) error {
	if value != nil && c.handle == 0 {
		handle, err := control.Connect("focus-in-event", focusSlotThunk, c)
		if err != nil {
			return err
		}
		c.handle = handle
	} else if value == nil && c.handle != 0 {
		control.HandlerDisconnect(c.handle)
		c.handle = 0
	}
	c.callback = value
	return nil
}

func (c *focusSlot) Close(control *gtk.Widget) {
	if c.handle != 0 {
		control.HandlerDisconnect(c.handle)
		c.handle = 0
	}
}

func focusSlotThunk(widget *gtk.Button, event *gdk.Event, c *focusSlot) bool {
	c.callback()
	return true
}

type blurSlot struct {
	callback func()
	handle   glib.SignalHandle
}

func (c *blurSlot) Set(control *gtk.Widget, value func()) error {
	if value != nil && c.handle == 0 {
		handle, err := control.Connect("focus-out-event", blurSlotThunk, c)
		if err != nil {
			return err
		}
		c.handle = handle
	} else if value == nil && c.handle != 0 {
		control.HandlerDisconnect(c.handle)
		c.handle = 0
	}
	c.callback = value
	return nil
}

func (c *blurSlot) Close(control *gtk.Widget) {
	if c.handle != 0 {
		control.HandlerDisconnect(c.handle)
		c.handle = 0
	}
}

func blurSlotThunk(widget *gtk.Button, event *gdk.Event, c *focusSlot) bool {
	c.callback()
	return true
}
