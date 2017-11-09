package goey

import (
	"sync/atomic"
)

var (
	currentControlID uint32 = 100
)

func nextControlID() uint32 {
	return atomic.AddUint32(&currentControlID, 1)
}
