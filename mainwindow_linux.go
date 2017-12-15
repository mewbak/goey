package goey

import (
	"sync/atomic"

	"github.com/gotk3/gotk3/gtk"
)

var (
	mainWindowCount int32 = 0
)

func init() {
	gtk.Init(nil)
}

type mainWindow struct {
	handle *gtk.Window
	vbox   mountedVBox
}

func newMainWindow(title string, children []Widget) (*MainWindow, error) {
	app, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		return nil, err
	}
	atomic.AddInt32(&mainWindowCount, 1)

	retval := &MainWindow{mainWindow{app, mountedVBox{}}}

	app.SetTitle(title)
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

func (mw *mainWindow) close() {
	if mw.handle != nil {
		mw.handle.Destroy()
		mw.handle = nil
	}
}

func (w *mainWindow) setAlignment(main MainAxisAlign, cross CrossAxisAlign) error {
	w.vbox.setAlignment(main, cross)
	return nil
}

func (mw *mainWindow) setChildren(children []Widget) error {
	// Defer to the vertical box holding the children.
	err := mw.vbox.setChildren(children)
	// ... and we're done
	return err
}

func (mw *mainWindow) setTitle(value string) error {
	mw.handle.SetTitle(value)
	return nil
}

func mainwindow_onDestroy(widget *gtk.Window, mw *MainWindow) {
	mw.handle = nil
	if c := atomic.AddInt32(&mainWindowCount, -1); c == 0 {
		gtk.MainQuit()
	}
}
