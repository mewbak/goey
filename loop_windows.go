package goey

import (
	"sync/atomic"
	"unsafe"

	"github.com/lxn/win"
)

func run() error {
	err := Loop(true)
	for err == nil {
		err = Loop(true)
	}
	if err == ErrQuit {
		err = nil
	}
	return err
}

func do(action func() error) error {
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
