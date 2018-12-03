package dialog

import (
	"bitbucket.org/rj/goey/loop"
	"testing"
	"path/filepath"
	"os"
)

func getwd(t *testing.T) string {
	path, err := os.Getwd()
	if err!=nil {
		t.Fatalf("Could not determien the working directory, %s", err)
	}
	return path
}

func TestNewOpenFile(t *testing.T) {
	init := func() error {
		// The following creates a modal dialog with a message.
		asyncKeyEscape()
		filename, err := NewOpenFile().WithTitle(t.Name()).Show()
		if err != nil {
			t.Errorf("Failed to show dialog, %s", err)
		}
		if filename != "" {
			t.Errorf("Unexpected filename, %s", filename)
		}

		// The following should create an error
		_, err = NewOpenFile().WithTitle("").Show()
		if err == nil {
			t.Errorf("Missing error")
		}

		// The following should return a filename.
		asyncType("abcd.txt\n")
		filename, err = NewOpenFile().Show()
		if err != nil {
			t.Errorf("Failed to show dialog, %s", err)
		}
		if expect := filepath.Join(getwd(t),"abcd.txt"); filename != expect {
			t.Errorf("Unexpected filename, %s", filename)
		}

		return nil
	}

	err := loop.Run(init)
	if err != nil {
		t.Fatalf("Failed to run event loop, %s", err)
	}
}
