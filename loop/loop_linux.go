// +build !gnustep

package loop

import (
	"testing"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const (
	// Flag to control behaviour of UnlockOSThread in Run.
	unlockThreadAfterRun = true
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

func runTesting(func() error) error {
	panic("unreachable")
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

func testMain(m *testing.M) int {
	// On GTK, we need to be locked to a thread, but not to a particular
	// thread.  No need for special coordination.
	return m.Run()
}
