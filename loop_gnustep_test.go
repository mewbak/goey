package goey

import (
	"os"
	"runtime"
	"sync/atomic"
	"testing"

	"bitbucket.org/rj/goey/cocoa"
)

func TestMain(m *testing.M) {
	println("lock thread")
	runtime.LockOSThread()
	if !cocoa.IsMainThread() {
		panic("not main thread")
	}
	atomic.StoreUint32(&isTesting, 1)
	defer func() {
		atomic.StoreUint32(&isTesting, 0)
	}()

	testingActions = make(chan func())
	testingSync = make(chan struct{})

	// call flag.Parse() here if TestMain uses flags
	wait := make(chan int, 1)
	go func() {
		println("m.Run()")
		wait <- m.Run()
		close(testingActions)
		println("m.Run() done")
	}()

	for a := range testingActions {
		println("run")
		if !cocoa.IsMainThread() {
			println("!!!! not main thread")
		} else {
			func() {
				atomic.StoreUint32(&isTesting, 0)
				defer func() {
					atomic.StoreUint32(&isTesting, 1)
				}()

				a()
			}()
			println("run done")
		}
		println("end run")
		testingSync <- struct{}{}
	}

	println("os.Exit")
	os.Exit(<-wait)
}
