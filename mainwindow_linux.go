package goey

import (
	"image"
	"sync/atomic"

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
	vbox   mountedVBox
}

func newWindow(title string, children []Widget) (*Window, error) {
	app, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		return nil, err
	}
	atomic.AddInt32(&mainWindowCount, 1)

	retval := &Window{windowImpl{app, mountedVBox{}}}

	app.SetTitle(title)
	app.SetBorderWidth(10)
	app.Connect("destroy", mainwindow_onDestroy, retval)
	app.Show()

	tmpVBox := VBox{Children: children}
	vbox, err := tmpVBox.Mount(NativeWidget{&app.Widget})
	if err != nil {
		app.Destroy()
		return nil, err
	}
	retval.vbox = *vbox.(*mountedVBox)

	return retval, nil
}

func (mw *windowImpl) close() {
	if mw.handle != nil {
		mw.handle.Destroy()
		mw.handle = nil
	}
}

func (w *windowImpl) setAlignment(main MainAxisAlign, cross CrossAxisAlign) error {
	w.vbox.setAlignment(main, cross)
	return nil
}

func (mw *windowImpl) setChildren(children []Widget) error {
	// Defer to the vertical box holding the children.
	err := mw.vbox.setChildren(children)
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
