// +build gnustep

package loop

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"

	"bitbucket.org/rj/assert"
	"bitbucket.org/rj/goey/cocoa"
)

const (
	// Flag to control behaviour of UnlockOSThread in Run.
	unlockThreadAfterRun = false
)

type PanicError struct {
	Str string
}

func (p *PanicError) Error() string {
	return p.Str
}

var (
	cocoaInit      sync.Once
	testingActions chan func() error
	testingSync    chan error
)

func init() {
	assert.Assert(cocoa.IsMainThread(), "Not main thread")
	runtime.LockOSThread()
}

func initRun() error {
	cocoaInit.Do(func() {
		assert.Assert(cocoa.IsMainThread(), "Not main thread")
		cocoa.Init()
	})

	return nil
}

func terminateRun() {
	// Do nothing
}

func run() error {
	assert.Assert(cocoa.IsMainThread(), "Not main thread")
	cocoa.Run()
	return nil
}

func runTesting(action func() error) error {
	testingActions <- action
	err := <-testingSync
	if v, ok := err.(*PanicError); ok {
		panic(v.Str)
	}
	return err
}

func do(action func() error) error {
	return cocoa.Do(action)
}

func stop() {
	cocoa.Stop()
}

func testMain(m *testing.M) int {
	runtime.LockOSThread()
	if !cocoa.IsMainThread() {
		panic("not main thread")
	}
	atomic.StoreUint32(&isTesting, 1)
	defer func() {
		atomic.StoreUint32(&isTesting, 0)
	}()

	testingActions = make(chan func() error)
	testingSync = make(chan error)

	// call flag.Parse() here if TestMain uses flags
	wait := make(chan int, 1)
	go func() {
		wait <- m.Run()
		close(testingActions)
	}()

	for a := range testingActions {
		if !cocoa.IsMainThread() {
			println("!!!! not main thread")
		} else {
			err := func() (err error) {
				atomic.StoreUint32(&isTesting, 0)
				defer func() {
					atomic.StoreUint32(&isTesting, 1)
				}()
				defer func() {
					if r := recover(); r != nil {
						if v, ok := r.(string); ok {
							err = &PanicError{Str: v}
						}
					}
				}()

				return Run(a)
			}()
			testingSync <- err
		}
	}

	return <-wait
}
