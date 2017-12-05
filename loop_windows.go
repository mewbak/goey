package goey

import (
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/lxn/win"
)

var (
	threadID      uint32
	mutexThreadID sync.Mutex
)

func setThreadID() error {
	mutexThreadID.Lock()
	defer mutexThreadID.Unlock()

	if threadID != 0 {
		return ErrAlreadyRunning
	}

	threadID = GetThreadID()
	return nil
}

func run() error {
	// Pin the GUI message loop to a single thread
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// Update the
	err := setThreadID()
	if err != nil {
		return err
	}
	defer func() {
		mutexThreadID.Lock()
		threadID = 0
		mutexThreadID.Unlock()
	}()

	err = Loop(true)
	for err == nil {
		err = Loop(true)
	}
	if err == ErrQuit {
		err = nil
	}
	return err
}

func do(action func() error) error {
	mutexThreadID.Lock()
	defer mutexThreadID.Unlock()

	if threadID == 0 {
		return ErrNotRunning
	}

	err := make(chan error, 1)
	win.PostMessage(win.HWND(atomic.LoadUintptr(&activeWindow)), win.WM_USER, uintptr(unsafe.Pointer(&action)), uintptr(unsafe.Pointer(&err)))
	return <-err
}

func Loop(blocking bool) error {
	// Obtain a copy of the next message from the queue.
	var msg win.MSG
	if blocking {
		win.GetMessage(&msg, 0, 0, 0)
	} else {
		win.PeekMessage(&msg, 0, 0, 0, win.PM_REMOVE)
	}

	// Processing for application wide messages are handled in this block.
	if msg.Message == win.WM_QUIT {
		return ErrQuit
	}
	if msg.Message == win.WM_USER {
		err := (*(*func() error)(unsafe.Pointer(msg.WParam)))()
		(*(*chan error)(unsafe.Pointer(msg.LParam))) <- err
		return nil
	}

	// Dispatch message.
	if !win.IsDialogMessage(win.HWND(activeWindow), &msg) {
		win.TranslateMessage(&msg)
		win.DispatchMessage(&msg)
	}
	return nil
}
