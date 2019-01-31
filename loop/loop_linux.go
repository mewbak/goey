package loop

import (
	"bitbucket.org/rj/goey/internal/nopanic"
	"bitbucket.org/rj/goey/internal/syscall"
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
	syscall.MainContextInvoke(func() {
		err <- nopanic.Wrap(action)
	})
	return nopanic.Unwrap(<-err)
}

func stop() {
	gtk.MainQuit()
}
