package cocoa

/*
#cgo CFLAGS: -x objective-c -I/usr/include/GNUstep
#cgo LDFLAGS: -lgnustep-gui -lgnustep-base -lobjc
#include "cocoa.h"
*/
import "C"
import "sync"

func Init() error {
	C.init()
	return nil
}

func Run() error {
	// Run the event loop.
	C.run()
	return nil
}

var (
	thunkAction func() error
	thunkErr    error
	thunkMutex  sync.Mutex
)

func Do(action func() error) error {
	thunkMutex.Lock()
	defer thunkMutex.Unlock()

	thunkAction = action
	C.thunkDo()
	return thunkErr
}

//export callbackDo
func callbackDo() {
	thunkErr = thunkAction()
}

func Stop() {
	C.stop()
}
