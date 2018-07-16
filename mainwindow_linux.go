package goey

import (
	"image"
	"sync/atomic"
	"unsafe"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var (
	mainWindowCount int32
	vscrollbarWidth Length
)

func init() {
	gtk.Init(nil)
}

func boolToPolicy(value bool) gtk.PolicyType {
	if value {
		return gtk.POLICY_ALWAYS
	}
	return gtk.POLICY_NEVER
}

type windowImpl struct {
	handle                  *gtk.Window
	scroll                  *gtk.ScrolledWindow
	layout                  *gtk.Layout
	child                   Element
	childMinSize            Size
	horizontalScroll        bool
	horizontalScrollVisible bool
	verticalScroll          bool
	verticalScrollVisible   bool
	onClosing               func() bool
	shClosing               glib.SignalHandle
}

func newWindow(title string, child Widget) (*Window, error) {
	app, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		return nil, err
	}
	atomic.AddInt32(&mainWindowCount, 1)

	scroll, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}
	scroll.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_NEVER)
	app.Add(scroll)

	layout, err := gtk.LayoutNew(nil, nil)
	if err != nil {
		return nil, err
	}
	scroll.Add(layout)

	retval := &Window{windowImpl{
		handle: app,
		scroll: scroll,
		layout: layout,
	}}
	app.SetTitle(title)
	app.SetBorderWidth(0)
	app.Connect("destroy", mainwindowOnDestroy, retval)
	app.Connect("size-allocate", mainwindowOnSizeAllocate, retval)
	app.SetDefaultSize(400, 400)
	app.ShowAll()
	retval.horizontalScroll, retval.verticalScroll = retval.scrollDefaults()

	err = retval.setChild(child)
	if err != nil {
		app.Destroy()
		return nil, err
	}
	retval.setScroll(retval.scrollDefaults())

	return retval, nil
}

func windowOnClosing(widget *gtk.Window, _ *gdk.Event, w *windowImpl) bool {
	return w.onClosing()
}

func (w *windowImpl) onSize() {
	if w.child == nil {
		return
	}

	// Update the global DPI
	DPI.X, DPI.Y = 96, 96

	width, height := w.handle.GetSize()
	clientSize := Size{FromPixelsX(width), FromPixelsY(height)}
	size := w.layoutChild(clientSize)
	if w.horizontalScroll && w.verticalScroll {
		// Show scroll bars if necessary.
		w.showScrollV(size.Height, clientSize.Height)
		ok := w.showScrollH(size.Width, clientSize.Width)
		// Adding horizontal scroll take vertical space, so we need to check
		// again for vertical scroll.
		if ok {
			_, height := w.handle.GetSize()
			w.showScrollV(size.Height, FromPixelsY(height))
		}
	} else if w.verticalScroll {
		// Show scroll bars if necessary.
		ok := w.showScrollV(size.Height, clientSize.Height)
		if ok {
			width, height := w.handle.GetSize()
			clientSize := Size{FromPixelsX(width), FromPixelsY(height)}
			size = w.layoutChild(clientSize)
		}
	} else if w.horizontalScroll {
		// Show scroll bars if necessary.
		ok := w.showScrollH(size.Width, clientSize.Width)
		if ok {
			width, height := w.handle.GetSize()
			clientSize := Size{FromPixelsX(width), FromPixelsY(height)}
			size = w.layoutChild(clientSize)
		}
	}
	w.layout.SetSize(uint(size.Width.PixelsX()), uint(size.Height.PixelsY()))
	bounds := Rectangle{
		Point{}, Point{size.Width, size.Height},
	}
	w.child.SetBounds(bounds)
}

func (w *windowImpl) getChild() Element {
	return w.child
}

func (mw *windowImpl) close() {
	if mw.handle != nil {
		mw.handle.Destroy()
		mw.handle = nil
	}
}

func (w *windowImpl) message(m *Message) {
	m.title, m.err = w.handle.GetTitle()
	m.handle = uintptr(unsafe.Pointer(w.handle))
}

func get_vscrollbar_width(window *gtk.Window) (Length, error) {
	if vscrollbarWidth != 0 {
		return vscrollbarWidth, nil
	}

	oldChild, err := window.GetChild()
	if err != nil {
		return 0, err
	}
	window.Remove(oldChild)

	sb, err := gtk.ScrollbarNew(gtk.ORIENTATION_VERTICAL, nil)
	if err != nil {
		return 0, err
	}

	window.Add(sb)
	sb.Show()
	_, retval := sb.GetPreferredWidth()
	sb.Destroy()
	window.Add(oldChild)
	vscrollbarWidth = FromPixelsX(retval)
	return vscrollbarWidth, nil
}

func (w *windowImpl) setChild(child Widget) (err error) {
	// Update the child element
	w.child, err = DiffChild(Control{&w.layout.Widget}, w.child, child)
	// Whether or not an error has occured, redo the layout so the children
	// are placed.
	if w.child != nil {
		// Update the global DPI
		DPI.X, DPI.Y = 96, 96

		// Constrain window size
		w.updateWindowMinSize()
		// Properties may have changed sizes, so we need to do layout.
		w.onSize()
	}
	// ... and we're done
	return err
}

func (w *windowImpl) setScroll(horz, vert bool) {
	w.horizontalScroll = horz
	w.verticalScroll = vert
	// Hide the scrollbars as a reset
	w.scroll.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_NEVER)
	w.horizontalScrollVisible = false
	w.verticalScrollVisible = false
	// Redo layout to account for new box constraints, and show
	// scrollbars if necessary
	w.onSize()
}

func (w *windowImpl) showScrollH(width Length, clientWidth Length) bool {
	if width > clientWidth {
		if !w.horizontalScrollVisible {
			// Show the scrollbar
			w.scroll.SetPolicy(gtk.POLICY_ALWAYS, boolToPolicy(w.verticalScrollVisible))
			w.horizontalScrollVisible = true
			return true
		}
	} else if w.horizontalScrollVisible {
		// Remove the scroll bar
		w.scroll.SetPolicy(gtk.POLICY_NEVER, boolToPolicy(w.verticalScrollVisible))
		w.horizontalScrollVisible = false
		return true
	}

	return false
}

func (w *windowImpl) showScrollV(height Length, clientHeight Length) bool {
	if height > clientHeight {
		if !w.verticalScrollVisible {
			// Show the scrollbar
			w.scroll.SetPolicy(boolToPolicy(w.horizontalScrollVisible), gtk.POLICY_ALWAYS)
			w.verticalScrollVisible = true
			return true
		}
	} else if w.verticalScrollVisible {
		// Remove the scroll bar
		w.scroll.SetPolicy(boolToPolicy(w.horizontalScrollVisible), gtk.POLICY_NEVER)
		w.verticalScrollVisible = false
		return true
	}

	return false
}

func (w *windowImpl) setIcon(img image.Image) error {
	pixbuf, _, err := imageToPixbuf(img)
	if err != nil {
		return err
	}
	w.handle.SetIcon(pixbuf)
	return nil
}

func (w *windowImpl) setOnClosing(callback func() bool) {
	w.onClosing = callback
	w.shClosing = setSignalHandler(&w.handle.Widget, 0, callback != nil, "delete-event", windowOnClosing, w)
}

func (w *windowImpl) setTitle(value string) error {
	w.handle.SetTitle(value)
	return nil
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
		w.handle.SetSizeRequest(dx, dy)
		return
	}

	request := image.Point{}
	// Determine the minimum size (in pixels) for the child of the window
	if w.horizontalScroll && w.verticalScroll {
		width := w.child.MinIntrinsicWidth(Inf)
		height := w.child.MinIntrinsicHeight(Inf)
		request.X = width.PixelsX() + dx
		request.Y = height.PixelsY() + dy
	} else if w.horizontalScroll {
		height := w.child.MinIntrinsicHeight(Inf)
		size := w.child.Layout(TightHeight(height))
		request.X = size.Width.PixelsX() + dx
		request.Y = height.PixelsY() + dy
	} else if w.verticalScroll {
		width := w.child.MinIntrinsicWidth(Inf)
		size := w.child.Layout(TightWidth(width))
		request.X = width.PixelsX() + dx
		request.Y = size.Height.PixelsY() + dy
	} else {
		width := w.child.MinIntrinsicWidth(Inf)
		height := w.child.MinIntrinsicHeight(Inf)
		size1 := w.child.Layout(TightWidth(width))
		size2 := w.child.Layout(TightHeight(height))
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

	w.handle.SetSizeRequest(request.X, request.Y)
}

func mainwindowOnDestroy(widget *gtk.Window, mw *Window) {
	mw.handle = nil
	if c := atomic.AddInt32(&mainWindowCount, -1); c == 0 {
		gtk.MainQuit()
	}
}

func mainwindowOnSizeAllocate(widget *gtk.Window, rect uintptr, mw *Window) {
	mw.onSize()
}
