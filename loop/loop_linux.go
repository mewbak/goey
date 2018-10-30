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
	return nil
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

func stop() {
	gtk.MainQuit()
}
