package cocoa

/*
#cgo CFLAGS: -x objective-c -I/usr/include/GNUstep
#cgo LDFLAGS: -lgnustep-gui -lgnustep-base -lobjc
#include "cocoa.h"
*/
import "C"
import "sync"

func Init() {
	C.init()
}

func Run() error {
	// Run the event loop.
	println("run")
	C.run()
	return nil
}

var (
	thunkAction func() error
	thunkErr    chan error
	thunkMutex  sync.Mutex
)

func Do(action func() error) error {
	thunkMutex.Lock()
	defer thunkMutex.Unlock()

	thunkAction = action
	thunkErr = make(chan error, 1)
	C.thunkDo()
	return <-thunkErr
}

//export callbackDo
func callbackDo() {
	err := thunkAction()
	thunkErr <- err
}

func Stop() {
	println("stop")
	C.stop()
}
