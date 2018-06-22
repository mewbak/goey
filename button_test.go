package goey

import (
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
	var render func() Widget

	// Update function
	update = func() {
		err := mainWindow.SetChild(render())
		if err != nil {
			panic(err)
		}
	}

	// Render function generates a tree of Widgets to describe the desired
	// state of the GUI.
	render = func() Widget {
		// Prep - text for the button
		text := "Click me!"
		if clickCount > 0 {
			text = text + "  (" + strconv.Itoa(clickCount) + ")"
		}
		// The GUI contains a single widget, this button.
		return &VBox{
			AlignMain:  MainCenter,
			AlignCross: CrossCenter,
			Children: []Widget{
				&Button{Text: text, OnClick: func() {
					clickCount++
					update()
				}},
			}}
	}
}

func TestButtonCreate(t *testing.T) {
	testingRenderWidgets(t, []Widget{
		&Button{Text: "A"},
		&Button{Text: "D", Disabled: true},
		&Button{Text: "E", Default: true},
	})
}

func TestButtonEvents(t *testing.T) {
	testingCheckFocusAndBlur(t, []Widget{
		&Button{Text: "A"},
		&Button{Text: "B"},
		&Button{Text: "C"},
	})
}

func TestButtonUpdateProps(t *testing.T) {
	testingUpdateWidgets(t, []Widget{
		&Button{Text: "A"},
		&Button{Text: "D", Disabled: true},
		&Button{Text: "E", Default: true},
	}, []Widget{
		&Button{Text: "AB"},
		&Button{Text: "DB", Disabled: false},
		&Button{Text: "EB", Disabled: true, Default: true},
	})
}
