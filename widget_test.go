package goey

import (
	"bytes"
	"reflect"
	"sync/atomic"
	"testing"
)

func testingRenderWidgets(t *testing.T, widgets []Widget) {
	init := func() error {
		// Create the window.  Some of the tests here are not expected in
		// production code, but we can be a little paranoid here.
		if c := atomic.LoadInt32(&mainWindowCount); c != 0 {
			t.Fatalf("Want mainWindow==0, got mainWindow==%d", c)
		}
		window, err := NewWindow(t.Name(), widgets)
		if err != nil {
			t.Errorf("Failed to create window, %s", err)
			return nil
		}
		if window == nil {
			t.Errorf("Unexpected nil for window")
			return nil
		}
		if c := atomic.LoadInt32(&mainWindowCount); c != 1 {
			t.Errorf("Want mainWindow==1, got mainWindow==%d", c)
			return nil
		}

		// Check that the controls that were mounted match with the list
		if children := window.Children(); children != nil {
			if len(children) != len(widgets) {
				t.Errorf("Wanted len(children) == len(widgets), got %d and %d", len(children), len(widgets))
			} else {
				for i := range children {
					if n1, n2 := children[i].Kind(), widgets[i].Kind(); n1 != n2 {
						t.Errorf("Wanted children[%d].Kind() != widgets[%d].Kind(), got %s and %s", i, i, n1, n2)
					}
				}
			}
		} else {
			t.Errorf("Want window.Children()!=nil")
		}

		go func(window *Window) {
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
	if c := atomic.LoadInt32(&mainWindowCount); c != 0 {
		t.Errorf("Want mainWindow==0, got mainWindow==%d", c)
	}
}

func testingCheckFocusAndBlur(t *testing.T, widgets []Widget) {
	log := bytes.NewBuffer(nil)

	for i := byte(0); i < 3; i++ {
		s := reflect.ValueOf(widgets[i])
		letter := 'a' + i
		s.Elem().FieldByName("OnFocus").Set(reflect.ValueOf(func() {
			log.Write([]byte{'f', letter})
		}))
		s.Elem().FieldByName("OnBlur").Set(reflect.ValueOf(func() {
			log.Write([]byte{'b', letter})
		}))
	}

	init := func() error {
		window, err := NewWindow(t.Name(), widgets)
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
