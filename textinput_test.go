package goey

import (
	"fmt"
	"testing"

	"bitbucket.org/rj/goey/base"
)

func ExampleTextInput() {
	// In a full application, this variable would be updated to point to
	// the main window for the application.
	var mainWindow *Window
	// These functions are used to update the GUI.  See below
	var update func()
	var render func() base.Widget

	// Update function
	update = func() {
		err := mainWindow.SetChild(render())
		if err != nil {
			panic(err)
		}
	}

	// Render function generates a tree of Widgets to describe the desired
	// state of the GUI.
	render = func() base.Widget {
		// Prep - text for the button
		// The GUI contains a single widget, this button.
		return &VBox{Children: []base.Widget{
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
		}}
	}
}

func TestTextInput(t *testing.T) {
	testingRenderWidgets(t,
		&TextInput{Value: "A"},
		&TextInput{Value: "B", Placeholder: "..."},
		&TextInput{Value: "C", Disabled: true},
	)
}

func TestTextInputClose(t *testing.T) {
	testingCloseWidgets(t,
		&TextInput{Value: "A"},
		&TextInput{Value: "B", Placeholder: "..."},
		&TextInput{Value: "C", Disabled: true},
	)
}

func TestTextInputEvents(t *testing.T) {
	testingCheckFocusAndBlur(t,
		&TextInput{},
		&TextInput{},
		&TextInput{},
	)
}

func TestTextInputProps(t *testing.T) {
	testingUpdateWidgets(t, []base.Widget{
		&TextInput{Value: "A"},
		&TextInput{Value: "B", Placeholder: "..."},
		&TextInput{Value: "C", Disabled: true},
	}, []base.Widget{
		&TextInput{Value: "AA"},
		&TextInput{Value: "BA", Disabled: true},
		&TextInput{Value: "CA", Placeholder: "***", Disabled: false},
	})
}
