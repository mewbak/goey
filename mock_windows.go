package goey

import (
	"github.com/lxn/win"
)

func (w *mockElement) SetOrder(previous win.HWND) win.HWND {
	return previous
}
