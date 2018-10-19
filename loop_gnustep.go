// +build gnustep

package goey

import (
	"bitbucket.org/rj/goey/cocoa"
	"runtime"
	"sync"
)

var (
	cocoaInit      sync.Once
	mainLoopLaunch = make(chan struct{}, 0)
)

func init() {
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
