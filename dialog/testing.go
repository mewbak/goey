package dialog

import (
	"time"
)

func asyncKeyEnter() {
	go func() {
		time.Sleep(500 * time.Millisecond)
		typeKeys("\n")
	}()
}

func asyncKeyEscape() {
	go func() {
		time.Sleep(500 * time.Millisecond)
		typeKeys("\x1b")
	}()
}

func asyncType(s string) {
	go func() {
		time.Sleep(500 * time.Millisecond)
		typeKeys(s)
	}()
}
