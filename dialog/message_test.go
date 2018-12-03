package dialog

import (
	"bitbucket.org/rj/goey/loop"
	"fmt"
	"testing"
)

func ExampleNewMessage() {
	// The following creates a modal dialog with a message.
	err := NewMessage("Some text for the body of the dialog box.").WithTitle("Example").WithInfo().Show()
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func TestNewMessage(t *testing.T) {
	init := func() error {
		// The following creates a modal dialog with a message.
		asyncKeyEnter()
		err := NewMessage("Some text for the body of the dialog box.").WithTitle(t.Name()).WithInfo().Show()
		if err != nil {
			t.Errorf("Failed to show message, %s", err)
		}

		// The following should return an error.
		err = NewMessage("").Show()
		if err == nil {
			t.Errorf("Missing error")
		}

		err = NewMessage("Some text...").WithTitle("").Show()
		if err == nil {
			t.Errorf("Missing error")
		}
		return nil
	}

	err := loop.Run(init)
	if err != nil {
		t.Fatalf("Failed to run event loop, %s", err)
	}
}
