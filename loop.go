package goey

import (
	"github.com/lxn/win"
)

func Run() error {
	err := Loop(true)
	for err == nil {
		err = Loop(true)
	}
	if err == ErrQuit {
		err = nil
	}
	return err
}

func Loop(blocking bool) error {
	var msg win.MSG
	if blocking {
		win.GetMessage(&msg, 0, 0, 0)
	} else {
		win.PeekMessage(&msg, 0, 0, 0, win.PM_REMOVE)
	}

	if msg.Message == win.WM_QUIT {
		return ErrQuit
	}

	if !win.IsDialogMessage(win.HWND(activeWindow), &msg) {
		win.TranslateMessage(&msg)
		win.DispatchMessage(&msg)
	}
	return nil
}
