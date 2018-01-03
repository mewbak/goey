package goey

import (
	"fmt"
	"strconv"
	"sync/atomic"
	"testing"
	"time"
)

func ExampleNewWindow() {
	createWindow := func() error {
		mw, err := NewWindow("Test", []Widget{
			&Button{Text: "Click me!"},
		})
		if err != nil {
			fmt.Println("Error: ", err)
			return err
		}

		go func() {
			fmt.Println("Up")
			time.Sleep(50 * time.Millisecond)
			fmt.Println("Down")

			// Note:  No work after this call to Do, since the call to Run may be
			// terminated when the call to Do returns.
			Do(func() error {
				mw.Close()
				return nil
			})
		}()

		return nil
	}

	err := Run(createWindow)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// Output:
	// Up
	// Down
}

func TestWindow_SetAlignment(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping alignment tests")
	}

	createWindow := func() error {

		if c := atomic.LoadInt32(&mainWindowCount); c != 0 {
			t.Fatalf("Want mainWindow==0, got mainWindow==%d", c)
		}
		mw, err := NewWindow("Test", []Widget{
			&Button{Text: "Click me!"},
			&Button{Text: "Or me!"},
			&Button{Text: "But not me."},
		})
		if err != nil {
			t.Fatalf("Failed to create window, %s", err)
		}
		if mw == nil {
			t.Fatalf("Unexpected nil for window")
		}
		if c := atomic.LoadInt32(&mainWindowCount); c != 1 {
			t.Fatalf("Want mainWindow==1, got mainWindow==%d", c)
		}

		go func() {
			t.Logf("Starting alignment tests")
			for i := MainStart; i <= SpaceBetween; i++ {
				for j := Stretch; j <= CrossEnd; j++ {
					time.Sleep(50 * time.Millisecond)
					Do(func() error {
						mw.SetAlignment(i, j)
						return nil
					})
				}
			}
			time.Sleep(50 * time.Millisecond)
			t.Logf("Stopping alignment tests")

			// Note:  No work after this call to Do, since the call to Run may be
			// terminated when the call to Do returns.
			Do(func() error {
				mw.Close()
				return nil
			})
		}()

		return nil
	}

	err := Run(createWindow)
	if err != nil {
		t.Fatalf("Failed to run event loop, %s", err)
	}
	if c := atomic.LoadInt32(&mainWindowCount); c != 0 {
		t.Fatalf("Want mainWindow==0, got mainWindow==%d", c)
	}
}

func TestNewWindow_SetChildren(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping alignment tests")
	}

	createWindow := func() error {
		widgets := []Widget{}

		if c := atomic.LoadInt32(&mainWindowCount); c != 0 {
			t.Fatalf("Want mainWindow==0, got mainWindow==%d", c)
		}
		mw, err := NewWindow("Test", widgets)
		if err != nil {
			t.Fatalf("Failed to create window, %s", err)
		}
		if mw == nil {
			t.Fatalf("Unexpected nil for window")
		}
		if c := atomic.LoadInt32(&mainWindowCount); c != 1 {
			t.Fatalf("Want mainWindow==1, got mainWindow==%d", c)
		}
		mw.SetAlignment(SpaceBetween, CrossCenter)

		go func() {
			t.Logf("Starting set children tests")
			for i := 1; i < 10; i++ {
				time.Sleep(50 * time.Millisecond)
				widgets = append(widgets, &Button{Text: "Button " + strconv.Itoa(i)})
				err := Do(func() error {
					return mw.SetChildren(widgets)
				})
				if err != nil {
					t.Logf("Error setting children, %s", err)
				}
			}
			for i := len(widgets); i > 0; i-- {
				time.Sleep(50 * time.Millisecond)
				widgets = widgets[:i-1]
				err := Do(func() error {
					return mw.SetChildren(widgets)
				})
				if err != nil {
					t.Logf("Error setting children, %s", err)
				}
			}
			time.Sleep(100 * time.Millisecond)
			t.Logf("Stopping set children tests")

			// Note:  No work after this call to Do, since the call to Run may be
			// terminated when the call to Do returns.
			Do(func() error {
				mw.Close()
				return nil
			})
		}()

		return nil
	}

	err := Run(createWindow)
	if err != nil {
		t.Fatalf("Failed to run event loop, %s", err)
	}
	if c := atomic.LoadInt32(&mainWindowCount); c != 0 {
		t.Fatalf("Want mainWindow==0, got mainWindow==%d", c)
	}
}

func TestNewWindow_SetTitle(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping alignment tests")
	}

	createWindow := func() error {

		if c := atomic.LoadInt32(&mainWindowCount); c != 0 {
			t.Fatalf("Want mainWindow==0, got mainWindow==%d", c)
		}
		mw, err := NewWindow("Test", []Widget{
			&Button{Text: "Click me!"},
			&Button{Text: "Or me!"},
			&Button{Text: "But not me."},
		})
		if err != nil {
			t.Fatalf("Failed to create window, %s", err)
		}
		if mw == nil {
			t.Fatalf("Unexpected nil for window")
		}
		if c := atomic.LoadInt32(&mainWindowCount); c != 1 {
			t.Fatalf("Want mainWindow==1, got mainWindow==%d", c)
		}

		go func() {
			t.Logf("Starting set title tests")
			time.Sleep(50 * time.Millisecond)
			err := Do(func() error {
				return mw.SetTitle("Flash!")
			})
			if err != nil {
				t.Errorf("Error calling SetTitle, %", err)
			}
			time.Sleep(50 * time.Millisecond)
			t.Logf("Stopping alignment tests")

			// Note:  No work after this call to Do, since the call to Run may be
			// terminated when the call to Do returns.
			err = Do(func() error {
				mw.Close()
				return nil
			})
			if err != nil {
				t.Errorf("Error calling Close, %", err)
			}
		}()

		return nil
	}

	err := Run(createWindow)
	if err != nil {
		t.Fatalf("Failed to run event loop, %s", err)
	}
	if c := atomic.LoadInt32(&mainWindowCount); c != 0 {
		t.Fatalf("Want mainWindow==0, got mainWindow==%d", c)
	}
}
