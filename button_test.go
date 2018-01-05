package goey

import (
	"bytes"
	"strconv"
	"testing"
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

func TestButtonCreate(t *testing.T) {
	testingRenderWidgets(t, []Widget{
		&Button{Text: "A"},
		&Button{Text: "D", Disabled: true},
	})
}

func TestButtonEvents(t *testing.T) {
	log := bytes.NewBuffer(nil)

	init := func() error {
		window, err := NewWindow("TestButtonEvents", []Widget{
			&Button{Text: "A", OnFocus: func() { log.Write([]byte{'f', 'a'}) }, OnBlur: func() { log.Write([]byte{'b', 'a'}) }},
			&Button{Text: "B", OnFocus: func() { log.Write([]byte{'f', 'b'}) }, OnBlur: func() { log.Write([]byte{'b', 'b'}) }},
			&Button{Text: "C", OnFocus: func() { log.Write([]byte{'f', 'c'}) }, OnBlur: func() { log.Write([]byte{'b', 'c'}) }},
		})
		if err != nil {
			t.Errorf("Failed to create window, %s", err)
		}

		go func(window *Window) {
			// Run the actions, which are counted.
			for i := 0; i < 3; i++ {
				err := Do(func() error {
					testingSetFocus(t, window, i)
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
	if s := log.String(); s != "fabafbbbfcbc" {
		t.Errorf("Incorrect log string, got log==%s", s)
	}
}
