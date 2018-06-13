package goey

import (
	"image"
	"sync/atomic"
	"unsafe"

	"bitbucket.org/rj/goey/syscall"
	"github.com/gotk3/gotk3/gtk"
)

var (
	mainWindowCount int32 = 0
)

func init() {
	gtk.Init(nil)
}

type windowImpl struct {
	handle *gtk.Window
	scroll *gtk.ScrolledWindow
	layout *gtk.Layout
	vbox   mountedVBox
}

func newWindow(title string, children []Widget) (*Window, error) {
	app, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		return nil, err
	}
	atomic.AddInt32(&mainWindowCount, 1)

	scroll, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, err
	}
	app.Add(scroll)

	layout, err := gtk.LayoutNew(nil, nil)
	if err != nil {
		return nil, err
	}
	scroll.Add(layout)

	retval := &Window{windowImpl{app, scroll, layout, mountedVBox{}}}
	app.SetTitle(title)
	app.SetBorderWidth(0)
	app.Connect("destroy", mainwindow_onDestroy, retval)
	app.Connect("size-allocate", mainwindow_onSizeAllocate, retval)
	app.SetDefaultSize(400, 400)
	app.ShowAll()

	tmpVBox := VBox{Children: children, AlignMain: MainStart, AlignCross: Stretch}
	vbox, err := tmpVBox.Mount(Control{&layout.Widget})
	if err != nil {
		app.Destroy()
		return nil, err
	}
	retval.vbox = *vbox.(*mountedVBox)
	if len(retval.vbox.children) != len(children) {
		panic("Error in mounting children while creating window.")
	}
	syscall.WindowSetInteractiveDebugging(true)

	return retval, nil
}

func (w *windowImpl) doLayout() {
	DPI.X, DPI.Y = 96, 96

	_, _, width, height := syscall.WidgetGetAllocation(&w.scroll.Widget)

	a, _ := w.vbox.MeasureWidth()
	if w := a.PixelsX(); w > width {
		width = w
	} else {
		a = FromPixelsX(width)
	}
	b, _ := w.vbox.MeasureHeight(a)
	if h := b.PixelsY(); h > height {
		height = h
		width = width - 24
		a = FromPixelsX(width)
		b, _ = w.vbox.MeasureHeight(a)
		height = b.PixelsY()
	} else {
		b = FromPixelsY(height)
	}
	w.layout.SetSize(uint(width), uint(height))

	bounds := Rectangle{
		Point{},
		Point{
			a,
			b,
		},
	}
	w.vbox.SetBounds(bounds)
}

func (w *windowImpl) getAlignment() (MainAxisAlign, CrossAxisAlign) {
	return w.vbox.alignMain, w.vbox.alignCross
}

func (w *windowImpl) getChildren() []Element {
	return w.vbox.children
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

func (w *windowImpl) setAlignment(main MainAxisAlign, cross CrossAxisAlign) error {
	w.vbox.alignMain = main
	w.vbox.alignCross = cross
	w.doLayout()
	return nil
}

func (mw *windowImpl) setChildren(children []Widget) error {
	// Defer to the vertical box holding the children.
	err := mw.vbox.setChildren(children)
	// Properties may have changed sizes, so we need to do layout.
	mw.doLayout()
	// ... and we're done
	return err
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
