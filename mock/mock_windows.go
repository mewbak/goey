package mock

import (
	"github.com/lxn/win"
)

func (w *Element) SetOrder(previous win.HWND) win.HWND {
	return previous
}
