package goey

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func ExampleNewMainWindow() {
	mw, err := NewMainWindow("Test", []Widget{
		&Button{Text: "Click me!"},
	})
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	go func() {
		fmt.Println("Up")
		time.Sleep(100 * time.Millisecond)
		fmt.Println("Down")

		// Note:  No work after this call to Do, since the call to Run may be
		// terminated when the call to Do returns.
		Do(func() error {
			mw.Close()
			return nil
		})
	}()

	err = Run()
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// Output:
	// Up
	// Down
}

func TestNewMainWindow_SetAlignment(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping alignment tests")
	}

	mw, err := NewMainWindow("Test", []Widget{
		&Button{Text: "Click me!"},
		&Button{Text: "Or me!"},
		&Button{Text: "But not me."},
	})
	if err != nil {
		t.Fatalf("Failed to create main window, %s", err)
	}

	go func() {
		t.Logf("Starting alignment tests")
		for i := MainStart; i <= SpaceBetween; i++ {
			for j := Stretch; j <= CrossEnd; j++ {
				time.Sleep(100 * time.Millisecond)
				Do(func() error {
					mw.SetAlignment(i, j)
					return nil
				})
			}
		}
		time.Sleep(100 * time.Millisecond)
		t.Logf("Stopping alignment tests")

		// Note:  No work after this call to Do, since the call to Run may be
		// terminated when the call to Do returns.
		Do(func() error {
			mw.Close()
			return nil
		})
	}()

	err = Run()
	if err != nil {
		t.Fatalf("Failed to run event loop, %s", err)
	}
}

func TestNewMainWindow_SetChildren(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping alignment tests")
	}

	widgets := []Widget{}

	mw, err := NewMainWindow("Test", widgets)
	if err != nil {
		t.Fatalf("Failed to create main window, %s", err)
	}
	mw.SetAlignment(SpaceBetween, CrossCenter)

	go func() {
		t.Logf("Starting set children tests")
		for i := 1; i < 10; i++ {
			time.Sleep(100 * time.Millisecond)
			widgets = append(widgets, &Button{Text: "Button " + strconv.Itoa(i)})
			err := Do(func() error {
				return mw.SetChildren(widgets)
			})
			if err != nil {
				t.Logf("Error setting children, %s", err)
			}
		}
		for i := len(widgets); i > 0; i-- {
			time.Sleep(100 * time.Millisecond)
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

	err = Run()
	if err != nil {
		t.Fatalf("Failed to run event loop, %s", err)
	}
}

func TestNewMainWindow_SetTitle(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping alignment tests")
	}

	mw, err := NewMainWindow("Test", []Widget{
		&Button{Text: "Click me!"},
		&Button{Text: "Or me!"},
		&Button{Text: "But not me."},
	})
	if err != nil {
		t.Fatalf("Failed to create main window, %s", err)
	}

	go func() {
		t.Logf("Starting set title tests")
		time.Sleep(100 * time.Millisecond)
		mw.SetTitle("Flash!")
		time.Sleep(100 * time.Millisecond)
		t.Logf("Stopping alignment tests")

		// Note:  No work after this call to Do, since the call to Run may be
		// terminated when the call to Do returns.
		Do(func() error {
			mw.Close()
			return nil
		})
	}()

	err = Run()
	if err != nil {
		t.Fatalf("Failed to run event loop, %s", err)
	}
}
