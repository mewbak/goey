// +build gnustep

package goey

import (
	"bitbucket.org/rj/goey/cocoa"
)

func run_init() error {
	return cocoa.Init()
}

func run() error {
	// Run the event loop.
	return cocoa.Run()
}

func do(action func() error) error {
	return cocoa.Do(action)
}
