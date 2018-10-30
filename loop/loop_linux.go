package loop

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func init() {
	gtk.Init(nil)
}

func initRun() error {
	// Do nothing
}

func terminateRun() {
	// Do nothing
}

func run() error {
	gtk.Main()
	return nil
}

func do(action func() error) error {
	err := make(chan error, 1)
	glib.IdleAdd(func() {
		err <- action()
	})
	return <-err
}

func loop(blocking bool) error {
	rc := gtk.MainIterationDo(blocking)
	if rc {
		return ErrQuit
	}
	return nil
}

func stop() {
	gtk.MainQuit()
}
