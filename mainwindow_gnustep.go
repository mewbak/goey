// +build gnustep

package goey

import (
	"image"
	"sync/atomic"

	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/cocoa"
)

type windowImpl struct {
	handle           *cocoa.Window
	child            base.Element
	horizontalScroll bool
	verticalScroll   bool

	onClosing func() bool
}

func newWindow(title string, child base.Widget) (*Window, error) {
	// Don't want to run GUI init unless the program actually gets
	// around to starting the GUI.
	initCocoa.Do(func() {
		cocoa.Init()
	})

	// Update the global DPI
	base.DPI.X, base.DPI.Y = 96, 96

	println("newWindow")
	w, h := sizeDefaults()
	handle := cocoa.NewWindow(title, w, h)
	atomic.AddInt32(&mainWindowCount, 1)
	retval := &Window{windowImpl{
		handle: handle,
	}}
	handle.SetCallbacks((*windowCallbacks)(&retval.windowImpl))

	return retval, nil
}

func (w *windowImpl) control() base.Control {
	return base.Control{w.handle}
}

func (w *windowImpl) close() {
	if w.handle != nil {
		w.handle.Close()
		w.handle = nil
	}
}

func (w *windowImpl) message(m *Message) {
}

func (w *windowImpl) onSize() {
	if w.child == nil {
		return
	}

	// Update the global DPI
	base.DPI.X, base.DPI.Y = 96, 96

	width, height := w.handle.ContentSize()
	clientSize := base.Size{base.FromPixelsX(width), base.FromPixelsY(height)}
	size := w.layoutChild(clientSize)
	bounds := base.Rectangle{
		base.Point{}, base.Point{size.Width, size.Height},
	}
	w.child.SetBounds(bounds)
}

func (w *windowImpl) setChildPost() {
	// Redo the layout so the children are placed.
	if w.child != nil {
		// Update the global DPI
		base.DPI.X, base.DPI.Y = 96, 96

		// Constrain window size
		//w.updateWindowMinSize()
		// Properties may have changed sizes, so we need to do layout.
		w.onSize()
	} else {
	}
}

func (w *windowImpl) setScroll(horz, vert bool) {
	w.horizontalScroll = horz
	w.verticalScroll = vert
}

func (w *windowImpl) show() {
	//w.handle.ShowAll()
}

func (w *windowImpl) setIcon(img image.Image) error {
	return nil
}

func (w *windowImpl) setOnClosing(callback func() bool) {
	w.onClosing = callback
}

func (w *windowImpl) setTitle(value string) error {
	return nil
}

func (w *windowImpl) updateWindowMinSize() {
}

type windowCallbacks windowImpl

func (w *windowCallbacks) OnShouldClose() bool {
	if w.onClosing != nil {
		return !w.onClosing()
	}
	return true
}

func (w *windowCallbacks) OnWillClose() {
	w.handle = nil
	if newval := atomic.AddInt32(&mainWindowCount, -1); newval == 0 {
		cocoa.Stop()
	}
}

func (w *windowCallbacks) OnDidResize() {
	impl := (*windowImpl)(w)
	impl.onSize()
}
