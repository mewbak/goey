package goey

import (
	"testing"
)

func TestDo(t *testing.T) {
	err := Do(func() error {
		return nil
	})

	if err != ErrNotRunning {
		t.Errorf("Unexpected success in call to Do")
	}
}
