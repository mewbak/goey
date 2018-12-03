package dialog

import (
	"bitbucket.org/rj/goey/loop"
	"fmt"
	"github.com/lxn/win"
	"time"
	"unsafe"
)

func typeKeys(text string) chan error {
	err := make(chan error, 1)

	go func() {
		defer close(err)

		time.Sleep(50 * time.Millisecond)
		for _, r := range text {
			inp := [2]win.KEYBD_INPUT{
				{win.INPUT_KEYBOARD, win.KEYBDINPUT{}},
				{win.INPUT_KEYBOARD, win.KEYBDINPUT{}},
			}

			if r == '\n' {
				inp[0].Ki.WVk = win.VK_RETURN
				inp[1].Ki.WVk = win.VK_RETURN
				inp[1].Ki.DwFlags = win.KEYEVENTF_KEYUP
			} else if r == 0x1b {
				inp[0].Ki.WVk = win.VK_ESCAPE
				inp[1].Ki.WVk = win.VK_ESCAPE
				inp[1].Ki.DwFlags = win.KEYEVENTF_KEYUP
			} else {
				inp[0].Ki.WScan = uint16(r)
				inp[0].Ki.DwFlags = win.KEYEVENTF_UNICODE
				inp[1].Ki.WScan = uint16(r)
				inp[1].Ki.DwFlags = win.KEYEVENTF_UNICODE | win.KEYEVENTF_KEYUP
			}

			loop.Do(func() error {
				rc := win.SendInput(2, unsafe.Pointer(&inp), int32(unsafe.Sizeof(inp[0])))
				if rc != 2 {
					err <- fmt.Errorf("windows error, %x", win.GetLastError())
				}
				return nil
			})
			time.Sleep(10 * time.Millisecond)
		}
	}()

	return err
}
