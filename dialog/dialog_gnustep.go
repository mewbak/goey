// +build gnustep

package dialog

import "bitbucket.org/rj/goey/cocoa"

type dialogImpl struct {
	parent *cocoa.Window
}

func typeKeys(text string) chan error {
	err := make(chan error, 1)

	go func() {
		defer close(err)

		// TODO
	}()

	return err
}
