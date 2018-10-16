// +build gnustep

package goey

import (
	"sync"

	"bitbucket.org/rj/goey/cocoa"
)

var (
	initCocoa sync.Once
)

func run() error {

	// Don't want to run GUI init unless the program actually gets
	// around to starting the GUI.
	initCocoa.Do(func() {
		cocoa.Init()
	})

	// Run the event loop.
	return cocoa.Run()
}

func do(action func() error) error {
	return cocoa.Do(action)
}
