package goey

import (
	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/syscall"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type textareaElement struct {
	handle *gtk.TextView
	buffer *gtk.TextBuffer
	frame  *gtk.Frame

	minLines int
	onChange func(string)
	shChange glib.SignalHandle
	onFocus  focusSlot
	onBlur   blurSlot
}

func (w *TextArea) mount(parent base.Control) (base.Element, error) {
	buffer, err := gtk.TextBufferNew(nil)
	if err != nil {
		return nil, err
	}
	buffer.SetText(w.Value)

	control, err := gtk.TextViewNewWithBuffer(buffer)
	if err != nil {
		buffer.Unref()
		return nil, err
	}
	control.SetLeftMargin(3)
	control.SetRightMargin(3)
	control.SetMarginTop(3)    // missing function SetTopMargin
	control.SetMarginBottom(3) // missing function SetBottomMargin
	control.SetWrapMode(gtk.WRAP_WORD)

	swindow, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		control.RefSink()
		control.Destroy()
		control.Unref()
		return nil, err
	}
	swindow.Add(control)
	swindow.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_AUTOMATIC)
	swindow.SetVExpand(true)

	frame, err := gtk.FrameNew("")
	if err != nil {
		swindow.RefSink()
		swindow.Destroy()
		swindow.Unref()
		return nil, err
	}
	frame.SetShadowType(gtk.SHADOW_ETCHED_IN)
	frame.Add(swindow)
	parent.Handle.Add(frame)

	retval := &textareaElement{
		handle:   control,
		buffer:   buffer,
		frame:    frame,
		onChange: w.OnChange,
		minLines: minlinesDefault(w.MinLines),
	}

	control.Connect("destroy", textareaOnDestroy, retval)
	if w.OnChange != nil {
		sh, err := buffer.Connect("changed", textareaOnChanged, retval)
		if err != nil {
			panic("Failed to connect 'changed' event")
		}
		retval.shChange = sh
	}
	retval.onFocus.Set(&control.Widget, w.OnFocus)
	retval.onBlur.Set(&control.Widget, w.OnBlur)
	frame.ShowAll()

	return retval, nil
}

func textareaOnChanged(buffer *gtk.TextBuffer, mounted *textareaElement) {
	if mounted.onChange == nil {
		return
	}

	text, err := buffer.GetText(buffer.GetStartIter(), buffer.GetEndIter(), true)
	if err != nil {
		// TODO:  What is the correct reporting here
		return
	}
	mounted.onChange(text)
}

func textareaOnDestroy(widget *gtk.TextView, mounted *textareaElement) {
	mounted.handle = nil
}

func (w *textareaElement) Close() {
	if w.handle != nil {
		w.frame.Destroy()
		w.buffer.Unref()
		w.handle = nil
		w.frame = nil
		w.buffer = nil
	}
}

func (w *textareaElement) Handle() *gtk.Widget {
	return &w.handle.Widget
}

func (w *textareaElement) Layout(bc base.Constraints) base.Size {
	if !bc.HasBoundedWidth() {
		_, width := w.handle.GetPreferredWidth()
		_, height := w.handle.GetPreferredHeight()
		return bc.Constrain(base.Size{
			base.FromPixelsX(width), 
			base.FromPixelsY(height),
		})
	}

	width := bc.Max.Width
	_, height := syscall.WidgetGetPreferredHeightForWidth(&w.frame.Widget, width.PixelsX())
	return bc.Constrain(base.Size{width, base.FromPixelsX(height)})
}

func (w *textareaElement) MinIntrinsicHeight(width base.Length) base.Length {
	if width != base.Inf {
		height, _ := syscall.WidgetGetPreferredHeightForWidth(&w.frame.Widget, width.PixelsX())
		return base.FromPixelsY(height)
	}
	height, _ := w.frame.GetPreferredHeight()
	return base.FromPixelsY(height)
}

func (w *textareaElement) MinIntrinsicWidth(base.Length) base.Length {
	width, _ := w.frame.GetPreferredWidth()
	return base.FromPixelsX(width)
}

func (w *textareaElement) Props() base.Widget {
	buffer, err := w.handle.GetBuffer()
	if err != nil {
		panic("count not get buffer, " + err.Error())
	}
	value, err := buffer.GetText(buffer.GetStartIter(), buffer.GetEndIter(), true)
	if err != nil {
		panic("could not get text, " + err.Error())
	}
	return &TextArea{
		Value:    value,
		Disabled: !w.handle.GetSensitive(),
		MinLines: w.minLines,
		OnChange: w.onChange,
		OnFocus:  w.onFocus.callback,
		OnBlur:   w.onBlur.callback,
	}
}

func (w *textareaElement) SetBounds(bounds base.Rectangle) {
	pixels := bounds.Pixels()
	syscall.SetBounds(&w.frame.Widget, pixels.Min.X, pixels.Min.Y, pixels.Dx(), pixels.Dy())
}

func (w *textareaElement) updateProps(data *TextArea) error {
	// TextView will send a 'changed' event, even if the new value is the
	// same.  To stop an infinite loop, we need to protect by checking
	// ourselves.
	buffer, err := w.handle.GetBuffer()
	if err != nil {
		return err
	}
	oldText, err := buffer.GetText(buffer.GetStartIter(), buffer.GetEndIter(), true)
	if err != nil {
		return err
	}
	if data.Value != oldText {
		w.onChange = nil // temporarily break OnChange to prevent event
		buffer.SetText(data.Value)
	}

	w.minLines = data.MinLines
	w.onChange = data.OnChange
	if data.OnChange != nil && w.shChange == 0 {
		sh, err := buffer.Connect("changed", textinputOnChanged, w)
		if err != nil {
			panic("Failed to connect 'changed' event")
		}
		w.shChange = sh
	} else if data.OnChange == nil && w.shChange != 0 {
		buffer.HandlerDisconnect(w.shChange)
		w.shChange = 0
	}
	w.onFocus.Set(&w.handle.Widget, data.OnFocus)
	w.onBlur.Set(&w.handle.Widget, data.OnBlur)

	return nil
}
