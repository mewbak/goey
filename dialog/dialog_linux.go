package dialog

import (
	"time"

	"bitbucket.org/rj/goey/internal/syscall"
	"bitbucket.org/rj/goey/loop"
	"github.com/gotk3/gotk3/gtk"
)

type dialogImpl struct {
	parent *gtk.Window
}

var (
	activeDialogForTesting *gtk.Dialog
)

func typeKeys(text string) chan error {
	err := make(chan error, 1)

	go func() {
		defer close(err)

		time.Sleep(500 * time.Millisecond)
		for _, r := range text {
			loop.Do(func() error {
				syscall.WidgetSendKey(&activeDialogForTesting.Widget, r, 0, 0)
				return nil
			})
			time.Sleep(50 * time.Millisecond)

			loop.Do(func() error {
				syscall.WidgetSendKey(&activeDialogForTesting.Widget, r, 0, 1)
				return nil
			})
			time.Sleep(50 * time.Millisecond)
		}
	}()

	return err
}
