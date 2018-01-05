package goey

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func TestDoFailure(t *testing.T) {
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
	init := func() error {
		window, err := NewWindow("Test", nil)
		if err != nil {
			t.Fatalf("Fail in call to NewWindow, %s", err)
		}
		if window == nil {
			t.Fatalf("Unexpected nil for window")
		}
		if c := atomic.LoadInt32(&mainWindowCount); c != 1 {
			t.Fatalf("Want mainWindow==1, got mainWindow==%d", c)
		}

		go func() {
			time.Sleep(1 * time.Second)
			err := Do(func() error {
				window.Close()
				return nil
			})
			if err != nil {
				t.Errorf("Unexpected error in call to Do")
			}
		}()

		return nil
	}

	err := Run(init)
	if err != nil {
		t.Errorf("Unexpected error in call to Run")
	}
	if c := atomic.LoadInt32(&mainWindowCount); c != 0 {
		t.Errorf("Want mainWindow==0, got mainWindow==%d", c)
	}
}

func TestDo(t *testing.T) {
	count := uint32(0)

	init := func() error {
		if c := atomic.LoadInt32(&mainWindowCount); c != 0 {
			t.Fatalf("Want mainWindow==0, got mainWindow==%d", c)
		}
		window, err := NewWindow("TestDo", nil)
		if err != nil {
			t.Errorf("Failed to create window, %s", err)
		}
		if window == nil {
			t.Fatalf("Unexpected nil for window")
		}
		if c := atomic.LoadInt32(&mainWindowCount); c != 1 {
			t.Fatalf("Want mainWindow==1, got mainWindow==%d", c)
		}

		go func(window *Window) {
			// Run the actions, which are counted.
			for i := 0; i < 10; i++ {
				err := Do(func() error {
					atomic.AddUint32(&count, 1)
					return nil
				})
				if err != nil {
					t.Errorf("Error in Do, %s", err)
				}
			}

			// Close the window
			err := Do(func() error {
				window.Close()
				return nil
			})
			if err != nil {
				t.Errorf("Error in Do, %s", err)
			}
		}(window)

		return nil
	}

	err := Run(init)
	if err != nil {
		t.Errorf("Failed to run GUI loop, %s", err)
	}
	if c := atomic.LoadInt32(&mainWindowCount); c != 0 {
		t.Errorf("Want mainWindow==0, got mainWindow==%d", c)
	}
	if c := atomic.LoadUint32(&count); c != 10 {
		t.Errorf("Want count=10, got count==%d", c)
	}
}
