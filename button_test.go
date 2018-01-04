package goey

import (
	"bytes"
	"strconv"
	"sync/atomic"
	"testing"
	"time"
)

func ExampleButton() {
	clickCount := 0

	// In a full application, this variable would be updated to point to
	// the main window for the application.
	var mainWindow *Window
	// These functions are used to update the GUI.  See below
	var update func()
	var render func() []Widget

	// Update function
	update = func() {
		err := mainWindow.SetChildren(render())
		if err != nil {
			panic(err)
		}
	}

	// Render function generates a tree of Widgets to describe the desired
	// state of the GUI.
	render = func() []Widget {
		// Prep - text for the button
		text := "Click me!"
		if clickCount > 0 {
			text = text + "  (" + strconv.Itoa(clickCount) + ")"
		}
		// The GUI contains a single widget, this button.
		return []Widget{
			&Button{Text: text, OnClick: func() {
				clickCount++
				update()
			}},
		}
	}
}

func TestButton(t *testing.T) {
	log := bytes.NewBuffer(nil)

	init := func() error {
		if c := atomic.LoadInt32(&mainWindowCount); c != 0 {
			t.Fatalf("Want mainWindow==0, got mainWindow==%d", c)
		}
		mw, err := NewWindow("TestButton", []Widget{
			&Button{Text: "A", OnFocus: func() { log.Write([]byte{'f', 'a'}) }, OnBlur: func() { log.Write([]byte{'b', 'a'}) }},
			&Button{Text: "B", OnFocus: func() { log.Write([]byte{'f', 'b'}) }, OnBlur: func() { log.Write([]byte{'b', 'b'}) }},
			&Button{Text: "C", OnFocus: func() { log.Write([]byte{'f', 'c'}) }, OnBlur: func() { log.Write([]byte{'b', 'c'}) }},
		})
		if err != nil {
			t.Errorf("Failed to create window, %s", err)
		}
		if mw == nil {
			t.Fatalf("Unexpected nil for window")
		}
		if c := atomic.LoadInt32(&mainWindowCount); c != 1 {
			t.Fatalf("Want mainWindow==1, got mainWindow==%d", c)
		}

		go func(mw *Window) {
			err := Do(func() error {
				time.Sleep(100 * time.Millisecond)
				mw.Close()
				return nil
			})
			if err != nil {
				t.Errorf("Error in Do, %s", err)
			}
		}(mw)

		return nil
	}

	err := Run(init)
	if err != nil {
		t.Errorf("Failed to run GUI loop, %s", err)
	}
	if c := atomic.LoadInt32(&mainWindowCount); c != 0 {
		t.Errorf("Want mainWindow==0, got mainWindow==%d", c)
	}
}

func TestButtonEvents(t *testing.T) {
	log := bytes.NewBuffer(nil)
	count := uint32(0)

	init := func() error {
		if c := atomic.LoadInt32(&mainWindowCount); c != 0 {
			t.Fatalf("Want mainWindow==0, got mainWindow==%d", c)
		}
		mw, err := NewWindow("TestButtonEvents", []Widget{
			&Button{Text: "A", OnFocus: func() { log.Write([]byte{'f', 'a'}) }, OnBlur: func() { log.Write([]byte{'b', 'a'}) }},
			&Button{Text: "B", OnFocus: func() { log.Write([]byte{'f', 'b'}) }, OnBlur: func() { log.Write([]byte{'b', 'b'}) }},
			&Button{Text: "C", OnFocus: func() { log.Write([]byte{'f', 'c'}) }, OnBlur: func() { log.Write([]byte{'b', 'c'}) }},
		})
		if err != nil {
			t.Errorf("Failed to create window, %s", err)
		}
		if mw == nil {
			t.Fatalf("Unexpected nil for window")
		}
		if c := atomic.LoadInt32(&mainWindowCount); c != 1 {
			t.Fatalf("Want mainWindow==1, got mainWindow==%d", c)
		}

		go func(mw *Window) {
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
				mw.Close()
				return nil
			})
			if err != nil {
				t.Errorf("Error in Do, %s", err)
			}
		}(mw)

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
