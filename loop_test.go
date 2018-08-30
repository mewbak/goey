package goey

import (
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
)

func ExampleRun() {
	// This init function will be used to create a window on the GUI thread.
	init := func() error {
		// Create an empty window.
		window, err := NewWindow("Test", nil)
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}

		go func() {
			// Because of goroutine, we are now off the GUI thread.
			// Schedule an action.
			err := Do(func() error {
				window.Close()
				fmt.Println("...like tears in rain")
				return nil
			})
			if err != nil {
				fmt.Println("Error:", err)
			}
		}()

		return nil
	}

	err := Run(init)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Output:
	// ...like tears in rain
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
		// Verify that the test is starting in the correct state.
		if c := atomic.LoadInt32(&mainWindowCount); c != 0 {
			t.Errorf("Want mainWindow==0, got mainWindow==%d", c)
			return nil
		}

		// Create window and verify.
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
			// Try running the main loop again, but in parallel.  We should get an error.
			err := Run(func() error {
				return nil
			})
			if err != ErrAlreadyRunning {
				t.Errorf("Expected ErrAlreadyRunning, got %s", err)
			}

			// Close the window.  This should stop the GUI loop.
			err = Do(func() error {
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

func TestRunWithError(t *testing.T) {
	const errorString = "No luck"

	// Make sure that error is passed through to caller
	init := func() error {
		return errors.New(errorString)
	}

	err := Run(init)
	if err == nil {
		t.Errorf("Unexpected success, no error returned")
	} else if s := err.Error(); errorString != s {
		t.Errorf("Unexpected error, want %s, got %s", errorString, s)
	}
	if c := atomic.LoadInt32(&mainWindowCount); c != 0 {
		t.Errorf("Want mainWindow==0, got mainWindow==%d", c)
	}
}

func TestRunWithPanic(t *testing.T) {
	const errorString = "No luck"
	defer func() {
		r := recover()
		if r != nil {
			if s, ok := r.(string); !ok {
				t.Errorf("Unexpected recover, %v", r)
			} else if s != errorString {
				t.Errorf("Unexpected recover, %s", s)
			}
		} else {
			t.Errorf("Missing panic")
		}

		// Make sure that window count is properly maintained.
		if c := atomic.LoadInt32(&mainWindowCount); c != 0 {
			t.Errorf("Want mainWindow==0, got mainWindow==%d", c)
		}
	}()

	// Make sure that error is passed through to caller
	init := func() error {
		panic(errorString)
	}

	err := Run(init)
	if err == nil {
		t.Errorf("Unexpected success, no error returned")
	} else if s := err.Error(); errorString != s {
		t.Errorf("Unexpected error, want %s, got %s", errorString, s)
	}
	if c := atomic.LoadInt32(&mainWindowCount); c != 0 {
		t.Errorf("Want mainWindow==0, got mainWindow==%d", c)
	}
}

func TestRunWithWindowClose(t *testing.T) {
	// Make sure that error is passed through to caller
	init := func() error {
		if c := atomic.LoadInt32(&mainWindowCount); c != 0 {
			t.Errorf("Want mainWindow==0, got mainWindow==%d", c)
			return nil
		}

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

		window.Close()
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
		// Verify that the test is starting in the correct state.
		if c := atomic.LoadInt32(&mainWindowCount); c != 0 {
			t.Errorf("Want mainWindow==0, got mainWindow==%d", c)
			return nil
		}

		// Create window and verify.
		// We need at least one window open to maintain GUI loop.
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

func TestDoFailure(t *testing.T) {
	err := Do(func() error {
		return nil
	})

	if err != ErrNotRunning {
		t.Errorf("Unexpected success in call to Do")
	}
}
