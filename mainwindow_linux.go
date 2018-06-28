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
	retval.setScroll(retval.scollDefaults())

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
		minSize := w.child.MinimumSize()
		w.childMinSize = minSize
		if w.horizontalScroll && minSize.Width > 120*DIP {
			minSize.Width = 120 * DIP
		}
		if w.verticalScroll && minSize.Height > 120*DIP {
			minSize.Height = 120 * DIP
		}
		w.handle.SetSizeRequest(minSize.Width.PixelsX(), minSize.Height.PixelsY())
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

func mainwindow_onDestroy(widget *gtk.Window, mw *Window) {
	mw.handle = nil
	if c := atomic.AddInt32(&mainWindowCount, -1); c == 0 {
		gtk.MainQuit()
	}
}

func mainwindow_onSizeAllocate(widget *gtk.Window, rect uintptr, mw *Window) {
	mw.doLayout()
}
