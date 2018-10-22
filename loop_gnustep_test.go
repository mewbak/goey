package goey

import (
	"os"
	"sync/atomic"
	"testing"

	"bitbucket.org/rj/goey/cocoa"
)

func TestMain(m *testing.M) {
	// The init function for loop_gnustep calls LockOSThread.  That should
	// still be in effect.  We won't call LockOSThread here, but perhaps use
	// a little paranoia, and double check that we are on the main thread
	// before starting any tests.
	if !cocoa.IsMainThread() {
		panic("not main thread")
	}

	// Both the functions Run and RunTest check whether or not we are in
	// "testing" mode.  Set the flag before proceeding.
	atomic.StoreUint32(&isTesting, 1)
	defer func() {
		atomic.StoreUint32(&isTesting, 0)
	}()

	testingActions = make(chan func())
	testingSync = make(chan struct{})

	// call flag.Parse() here if TestMain uses flags
	wait := make(chan int, 1)
	go func() {
		wait <- m.Run()
		close(testingActions)
	}()

	for a := range testingActions {
		if !cocoa.IsMainThread() {
			panic("not main thread")
		} else {
			func() {
				atomic.StoreUint32(&isTesting, 0)
				defer func() {
					atomic.StoreUint32(&isTesting, 1)
				}()

				a()
			}()
		}
		testingSync <- struct{}{}
	}

	os.Exit(<-wait)
}
