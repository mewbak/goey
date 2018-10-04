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
}

func newWindow(title string, child base.Widget) (*Window, error) {
	// Don't want to run GUI init unless the program actually gets
	// around to starting the GUI.
	initCocoa.Do(func() {
		cocoa.Init()
	})

	println("newWindow")
	w, h := sizeDefaults()
	handle := cocoa.NewWindow(title, w, h )
	atomic.AddInt32(&mainWindowCount, 1)
	return &Window{windowImpl{
		handle: handle,
	}}, nil
}

func (w *windowImpl) control() base.Control {
	return base.Control{w.handle.Uintptr()}
}

func (w *windowImpl) close() {
	if w.handle != nil {
		w.handle.Close()
		w.handle = nil

	if newval:= atomic.AddInt32(&mainWindowCount, -1); newval==0 {
		cocoa.Stop()
	}

	}
}

func (w *windowImpl) message(m *Message) {
}

func (w *windowImpl) setChildPost() {
	// Redo the layout so the children are placed.
	if w.child != nil {
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
}

func (w *windowImpl) setTitle(value string) error {
	return nil
}

func (w *windowImpl) updateWindowMinSize() {
}
