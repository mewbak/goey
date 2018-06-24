package goey

import (
	"unsafe"

	"bitbucket.org/rj/goey/syscall"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type mountedTextArea struct {
	handle *gtk.TextView
	buffer *gtk.TextBuffer
	frame  *gtk.Frame

	minLines int
	onChange func(string)
	shChange glib.SignalHandle
	onFocus  focusSlot
	onBlur   blurSlot
}

func (w *TextArea) mount(parent Control) (Element, error) {
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
	(*gtk.Container)(unsafe.Pointer(parent.handle)).Add(frame)

	retval := &mountedTextArea{
		handle:   control,
		buffer:   buffer,
		frame:    frame,
		onChange: w.OnChange,
		minLines: minlinesDefault(w.MinLines),
	}

	control.Connect("destroy", textarea_onDestroy, retval)
	if w.OnChange != nil {
		sh, err := buffer.Connect("changed", textarea_onChanged, retval)
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

func textarea_onChanged(buffer *gtk.TextBuffer, mounted *mountedTextArea) {
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

func textarea_onDestroy(widget *gtk.TextView, mounted *mountedTextArea) {
	mounted.handle = nil
}

func (w *mountedTextArea) Close() {
	if w.handle != nil {
		w.frame.Destroy()
		w.buffer.Unref()
		w.handle = nil
		w.frame = nil
		w.buffer = nil
	}
}

func (w *mountedTextArea) Handle() *gtk.Widget {
	return &w.handle.Widget
}

func (w *mountedTextArea) Layout(bc Box) Size {
	_, width := w.handle.GetPreferredWidth()
	_, height := w.handle.GetPreferredHeight()
	return bc.Constrain(Size{FromPixelsX(width), FromPixelsY(height)})
}

func (w *mountedTextArea) MinimumSize() Size {
	width, _ := w.handle.GetPreferredWidth()
	height, _ := w.handle.GetPreferredHeight()
	return Size{FromPixelsX(width), FromPixelsY(height)}
}

func (w *mountedTextArea) SetBounds(bounds Rectangle) {
	pixels := bounds.Pixels()
	syscall.SetBounds(&w.frame.Widget, pixels.Min.X, pixels.Min.Y, pixels.Dx(), pixels.Dy())
}

func (w *mountedTextArea) updateProps(data *TextArea) error {
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

	w.onChange = data.OnChange
	if data.OnChange != nil && w.shChange == 0 {
		sh, err := buffer.Connect("changed", textinput_onChanged, w)
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
