package goey

import (
	"bytes"
	"testing"
	"time"
)

func TestCheckboxCreate(t *testing.T) {
	testingRenderWidgets(t, []Widget{
		&Checkbox{Value: false, Text: "A"},
		&Checkbox{Value: true, Text: "B", Disabled: true},
	})
}

func TestCheckboxEvents(t *testing.T) {
	log := bytes.NewBuffer(nil)

	init := func() error {
		window, err := NewWindow("TestCheckboxEvents", []Widget{
			&Checkbox{Text: "A", OnFocus: func() { log.Write([]byte{'f', 'a'}) }, OnBlur: func() { log.Write([]byte{'b', 'a'}) }},
			&Checkbox{Text: "B", OnFocus: func() { log.Write([]byte{'f', 'b'}) }, OnBlur: func() { log.Write([]byte{'b', 'b'}) }},
			&Checkbox{Text: "C", OnFocus: func() { log.Write([]byte{'f', 'c'}) }, OnBlur: func() { log.Write([]byte{'b', 'c'}) }},
		})
		if err != nil {
			t.Errorf("Failed to create window, %s", err)
		}

		go func(window *Window) {
			// Run the actions, which are counted.
			for i := 0; i < 3; i++ {
				time.Sleep(1000 * time.Millisecond)
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
