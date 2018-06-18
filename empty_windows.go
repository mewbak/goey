package goey

import (
	"github.com/lxn/win"
)

func (w *mountedEmpty) SetOrder(hwnd win.HWND) win.HWND {
	return hwnd
}
