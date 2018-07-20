package goey

import (
	"github.com/lxn/win"
)

func (w *expandElement) SetOrder(previous win.HWND) win.HWND {
	return w.child.SetOrder(previous)
}
