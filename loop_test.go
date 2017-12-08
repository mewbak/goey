package goey

import (
	"fmt"
	"testing"
	"time"
)

func TestDo(t *testing.T) {
	err := Do(func() error {
		return nil
	})

	if err != ErrNotRunning {
		t.Errorf("Unexpected success in call to Do")
	}
}

func ExampleDo() {
	err := Do(func() error {
		// Inside this closure, we will be executing only on the GUI thread.
		_, err := fmt.Println("Hello.")
		// Return the error (if any) back to the caller.
		return err
	})

	// Report on the success or failure
	fmt.Println("Previous call to fmt.Println resulted in ", err)
}

func TestRun(t *testing.T) {
	mw, err := NewMainWindow("Test", nil)
	if err != nil {
		t.Fatalf("Fail in call to NewMainWindow, %s", err)
	}

	go func() {
		time.Sleep(1 * time.Second)
		err := Do(func() error {
			mw.Close()
			return nil
		})
		if err != nil {
			t.Errorf("Unexpected error in call to Do")
		}
	}()

	err = Run()
	if err != nil {
		t.Errorf("Unexpected error in call to Run")
	}
}
