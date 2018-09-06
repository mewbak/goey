package goey

import (
	"fmt"
)

func ExampleNewMessage() {
	// The following creates a modal dialog with a message.
	err := NewMessage("Some text for the body of the dialog box.").WithTitle("Example").WithInfo().Show()
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
