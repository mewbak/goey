package dialog

import (
	"bitbucket.org/rj/goey/loop"
	"testing"
)

func TestNewSaveFile(t *testing.T) {
	init := func() error {
		// The following creates a modal dialog with a message.
		asyncKeyEscape()
		_, err := NewSaveFile().WithTitle(t.Name()).Show()
		if err != nil {
			t.Errorf("Failed to show message, %s", err)
		}

		// The following should create an error
		_, err = NewSaveFile().WithTitle("").Show()
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
