package goey

import (
	"image"
	"os"
	"sync/atomic"
	"unsafe"

	"bitbucket.org/rj/goey/syscall"
	"github.com/gotk3/gotk3/gtk"
)

var (
	mainWindowCount int32 = 0
	vscrollbarWidth Length
)

func init() {
	gtk.Init(nil)
}

func boolToPolicy(value bool) gtk.PolicyType {
	if value {
		return gtk.POLICY_AUTOMATIC
	}
	return gtk.POLICY_NEVER
}

type windowImpl struct {
	handle           *gtk.Window
	scroll           *gtk.ScrolledWindow
	layout           *gtk.Layout
	child            Element
	childMinSize     Size
	horizontalScroll bool
	verticalScroll   bool
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
	app.Connect("destroy", mainwindow_onDestroy, retval)
	app.Connect("size-allocate", mainwindow_onSizeAllocate, retval)
	app.SetDefaultSize(400, 400)
	app.ShowAll()

	err = retval.setChild(child)
	if err != nil {
		app.Destroy()
		return nil, err
	}
	if len(os.Getenv("GOEY_DEBUGGING")) > 0 {
		syscall.WindowSetInteractiveDebugging(true)
	}
	retval.setScroll(retval.scrollDefaults())

	return retval, nil
}

func (w *windowImpl) doLayout() {
	if w.child == nil {
		return
	}

	// Update the global DPI
	DPI.X, DPI.Y = 96, 96

	width, height := w.handle.GetSize()
	size := w.layoutChild(Size{FromPixelsX(width), FromPixelsY(height)})
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
		w.doLayout()
	}
	// ... and we're done
	return err
}

func (mw *windowImpl) setScroll(horz, vert bool) {
	mw.horizontalScroll = horz
	mw.verticalScroll = vert
	mw.scroll.SetPolicy(boolToPolicy(horz), boolToPolicy(vert))
}

func (mw *windowImpl) setIcon(img image.Image) error {
	pixbuf, _, err := imageToPixbuf(img)
	if err != nil {
		return err
	}
	mw.handle.SetIcon(pixbuf)
	return nil
}

func (mw *windowImpl) setTitle(value string) error {
	mw.handle.SetTitle(value)
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

func mainwindow_onDestroy(widget *gtk.Window, mw *Window) {
	mw.handle = nil
	if c := atomic.AddInt32(&mainWindowCount, -1); c == 0 {
		gtk.MainQuit()
	}
}

func mainwindow_onSizeAllocate(widget *gtk.Window, rect uintptr, mw *Window) {
	mw.doLayout()
}
