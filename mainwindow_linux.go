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

type MainWindow struct {
	handle *gtk.Window
	vbox   mountedVBox
}

func newMainWindow(title string, children []Widget) (*MainWindow, error) {
	app, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		return nil, err
	}
	atomic.AddInt32(&mainWindowCount, 1)

	retval := &MainWindow{app, mountedVBox{}}

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

func (mw *MainWindow) Close() {
	if mw.handle != nil {
		mw.handle.Destroy()
		mw.handle = nil
	}
}

func (mw *MainWindow) SetChildren(children []Widget) error {
	// Defer to the vertical box holding the children.
	vbox := VBox{children}
	err := mw.vbox.UpdateProps(&vbox)
	// ... and we're done
	return err
}

func mainwindow_onDestroy(widget *gtk.Window, mw *MainWindow) {
	mw.handle = nil
	if c := atomic.AddInt32(&mainWindowCount, -1); c == 0 {
		gtk.MainQuit()
	}
}
