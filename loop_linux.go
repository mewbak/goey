package goey

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func init() {
	gtk.Init(nil)
}

func run() error {
	gtk.Main()
	return nil
}

func do(action func() error) error {
	err := make(chan error, 1)
	glib.IdleAdd( func() {
		err  <- action()
	})
	return <-err
}

func Loop(blocking bool) error {
	rc := gtk.MainIterationDo(blocking)
	if rc {
		return ErrQuit
	}
	return nil
}
