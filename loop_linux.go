package goey

import (
	"github.com/gotk3/gotk3/gtk"
)

func init() {
	gtk.Init(nil)
}

func run() error {
	gtk.Main()
	return nil
}

func Loop(blocking bool) error {
	rc := gtk.MainIterationDo(blocking)
	if rc {
		return ErrQuit
	}
	return nil
}
