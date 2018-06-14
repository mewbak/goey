package main

import (
	"github.com/lxn/win"
)

func (w *mountedGrid) SetOrder(previous win.HWND) win.HWND {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			previous = w.children[i][j].SetOrder(previous)
		}
	}
	return previous
}
