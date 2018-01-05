package goey

import (
	"bytes"
	"fmt"
	"testing"
)

func ExampleTextInput() {
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
		// The GUI contains a single widget, this button.
		return []Widget{
			&Label{Text: "Enter you text below:"},
			&TextInput{
				Value:       "",
				Placeholder: "Enter your data here",
				OnChange: func(value string) {
					fmt.Println("Change: ", value)
					// In a real example, you would update your data, and then
					// need to render the window again.
					update()
				},
			},
		}
	}
}

func TestTextInput(t *testing.T) {
	testingRenderWidgets(t, []Widget{
		&TextInput{Value: "A"},
		&TextInput{Value: "B", Placeholder: "..."},
		&TextInput{Value: "C", Disabled: true},
	})
}

func TestTextInputEvents(t *testing.T) {
	log := bytes.NewBuffer(nil)

	init := func() error {
		window, err := NewWindow("TestButtonEvents", []Widget{
			&TextInput{OnFocus: func() { log.Write([]byte{'f', 'a'}) }, OnBlur: func() { log.Write([]byte{'b', 'a'}) }},
			&Button{OnFocus: func() { log.Write([]byte{'f', 'b'}) }, OnBlur: func() { log.Write([]byte{'b', 'b'}) }},
			&Button{OnFocus: func() { log.Write([]byte{'f', 'c'}) }, OnBlur: func() { log.Write([]byte{'b', 'c'}) }},
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
