package goey

import (
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
	testingCheckFocusAndBlur(t, []Widget{
		&TextInput{},
		&TextInput{},
		&TextInput{},
	})
}