// +build gnustep

package goey

import (
	"runtime"
	"sync"

	"bitbucket.org/rj/goey/cocoa"
)

var (
	cocoaInit sync.Once
)

func init() {
	if !cocoa.IsMainThread() {
		panic("not main tread")
	}
	runtime.LockOSThread()
}

func main_loop_init() {
	cocoaInit.Do(func() {
		if !cocoa.IsMainThread() {
			panic("not main thread")
		}

		cocoa.Init()
	})
}

func run() error {
	if !cocoa.IsMainThread() {
		panic("not main thread")
	}

	cocoa.Run()
	return nil
}

func do(action func() error) error {
	return cocoa.Do(action)
}

func stop() {
	cocoa.Stop()
}
