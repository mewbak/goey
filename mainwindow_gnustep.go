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
	main_loop_init()

	// Update the global DPI
	base.DPI.X, base.DPI.Y = 96, 96

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
		w.updateWindowMinSize()
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
	w.handle.SetTitle(value)
	return nil
}

func (w *windowImpl) title() (string, error) {
	return w.handle.Title(), nil
}

func (w *windowImpl) updateWindowMinSize() {
	// Determine the extra width and height required for borders, title bar,
	// and scrollbars
	dx, dy := 0, 0
	if w.verticalScroll {
		// TODO:  Measure scrollbar width
		dx += 15
	}
	if w.horizontalScroll {
		// TODO:  Measure scrollbar height
		dy += 15
	}

	// If there is no child, then we just need enough space for the window chrome.
	if w.child == nil {
		w.handle.SetMinSize(dx, dy)
		return
	}

	request := image.Point{}
	// Determine the minimum size (in pixels) for the child of the window
	if w.horizontalScroll && w.verticalScroll {
		width := w.child.MinIntrinsicWidth(base.Inf)
		height := w.child.MinIntrinsicHeight(base.Inf)
		request.X = width.PixelsX() + dx
		request.Y = height.PixelsY() + dy
	} else if w.horizontalScroll {
		height := w.child.MinIntrinsicHeight(base.Inf)
		size := w.child.Layout(base.TightHeight(height))
		request.X = size.Width.PixelsX() + dx
		request.Y = height.PixelsY() + dy
	} else if w.verticalScroll {
		width := w.child.MinIntrinsicWidth(base.Inf)
		size := w.child.Layout(base.TightWidth(width))
		request.X = width.PixelsX() + dx
		request.Y = size.Height.PixelsY() + dy
	} else {
		width := w.child.MinIntrinsicWidth(base.Inf)
		height := w.child.MinIntrinsicHeight(base.Inf)
		size1 := w.child.Layout(base.TightWidth(width))
		size2 := w.child.Layout(base.TightHeight(height))
		request.X = max(width, size2.Width).PixelsX() + dx
		request.Y = max(height, size1.Height).PixelsY() + dy
	}

	// If scrolling is enabled for either direction, we can relax the
	// minimum window size.  These limits are fairly arbitrary, but we do need to
	// leave enough space for the scroll bars.
	if limit := (120 * DIP).PixelsX(); w.horizontalScroll && request.X > limit {
		request.X = limit
	}
	if limit := (120 * DIP).PixelsY(); w.verticalScroll && request.Y > limit {
		request.Y = limit
	}

	w.handle.SetMinSize(request.X, request.Y)
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
