package goey

import (
	"github.com/lxn/win"
)

func (w *alignElement) SetOrder(previous win.HWND) win.HWND {
	if w.child != nil {
		previous = w.child.SetOrder(previous)
	}
	return previous
}
