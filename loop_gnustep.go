// +build gnustep

package goey

import (
	"bitbucket.org/rj/goey/cocoa"
	"sync"
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
	println("do")
	defer func() { println("enddo") }()
	return cocoa.Do(action)
}
