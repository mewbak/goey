package goey

import (
	"bytes"
	"math"
	"reflect"
	"runtime"
	"testing"
	"time"

	"bitbucket.org/rj/goey/base"
	"bitbucket.org/rj/goey/loop"
)

type Proper interface {
	Props() base.Widget
}

type Clickable interface {
	Click()
}

type Focusable interface {
	TakeFocus() bool
}

func equal(t *testing.T, lhs, rhs base.Widget) bool {
	// On windows, the message EM_GETCUEBANNER does not work unless the manifest
	// is set correctly.  This cannot be done for the package, since that
	// manifest will conflict with the manifest of any app.
	if runtime.GOOS == "windows" {
		if value := reflect.ValueOf(rhs).Elem().FieldByName("Placeholder"); value.IsValid() {
			placeholder := value.String()
			if placeholder != "" {
				t.Logf("Zeroing 'Placeholder' field during test")
			}
			value.SetString("")
		}

		if slider, ok := lhs.(*Slider); ok {
			if newValue := math.Round(slider.Value*8) / 8; slider.Value != newValue {
				t.Logf("Rounding 'Value' field during test from %f to %f", slider.Value, newValue)
				slider.Value = newValue
			}
		}
	}

	return reflect.DeepEqual(lhs, rhs)
}

func testingRenderWidgets(t *testing.T, widgets ...base.Widget) {
	init := func() error {
		// Create the window.  Some of the tests here are not expected in
		// production code, but we can be a little paranoid here.
		window, err := NewWindow(t.Name(), &VBox{Children: widgets})
		if err != nil {
			t.Errorf("Failed to create window, %s", err)
			return nil
		}
		if window == nil {
			t.Errorf("Unexpected nil for window")
			return nil
		}

		// Check that the controls that were mounted match with the list
		if children := window.children(); children != nil {
			if len(children) != len(widgets) {
				t.Errorf("Wanted len(children) == len(widgets), got %d and %d", len(children), len(widgets))
			} else {
				for i := range children {
					if n1, n2 := children[i].Kind(), widgets[i].Kind(); n1 != n2 {
						t.Errorf("Wanted children[%d].Kind() == widgets[%d].Kind(), got %s, want %s", i, i, n1, n2)
					} else if widget, ok := children[i].(Proper); ok {
						data := widget.Props()
						if n1, n2 := data.Kind(), widgets[i].Kind(); n1 != n2 {
							t.Errorf("Wanted data.Kind() == widgets[%d].Kind(), got %s, want %s", i, n1, n2)
						}
						if !equal(t, data, widgets[i]) {
							t.Errorf("Wanted data == widgets[%d], got %v, want %v", i, data, widgets[i])
						}
					} else {
						t.Logf("Cannot verify props of child")
					}
				}
			}
		} else {
			t.Errorf("Want window.Children()!=nil")
		}

		go func(window *Window) {
			if testing.Verbose() && !testing.Short() {
				time.Sleep(250 * time.Millisecond)
			}
			err := loop.Do(func() error {
				window.Close()
				return nil
			})
			if err != nil {
				t.Errorf("Error in Do, %s", err)
			}
		}(window)

		return nil
	}

	err := loop.Run(init)
	if err != nil {
		t.Errorf("Failed to run GUI loop, %s", err)
	}
}

func testingRenderWidgetsFail(t *testing.T, outError error, widgets ...base.Widget) {
	init := func() error {
		window, err := NewWindow(t.Name(), &VBox{Children: widgets})
		if window != nil {
			t.Errorf("Unexpected non-nil window")
		}
		if err != outError {
			if err == nil {
				t.Errorf("Unexpected nil error, want %s", outError)
			} else {
				t.Errorf("Unexpected error, want %v, got %s", outError, err)
			}
			return nil
		}
		return nil
	}

	err := loop.Run(init)
	if err != nil {
		t.Errorf("Failed to run GUI loop, %s", err)
	}
}

func testingCloseWidgets(t *testing.T, widgets ...base.Widget) {
	init := func() error {
		// Create the window.  Some of the tests here are not expected in
		// production code, but we can be a little paranoid here.
		window, err := NewWindow(t.Name(), &VBox{Children: widgets})
		if err != nil {
			t.Errorf("Failed to create window, %s", err)
			return nil
		}
		if window == nil {
			t.Errorf("Unexpected nil for window")
			return nil
		}

		// Check that the controls that were mounted match with the list
		if len(window.children()) != len(widgets) {
			t.Errorf("Want len(window.Children())!=nil")
		}

		err = window.SetChild(&VBox{Children: nil})
		if err != nil {
			t.Errorf("Failed to set children, %s", err)
			return nil
		}
		if len(window.children()) != 0 {
			t.Errorf("Want len(window.Children())!=0")
		}

		go func(window *Window) {
			err := loop.Do(func() error {
				window.Close()
				return nil
			})
			if err != nil {
				t.Errorf("Error in Do, %s", err)
			}
		}(window)

		return nil
	}

	err := loop.Run(init)
	if err != nil {
		t.Errorf("Failed to run GUI loop, %s", err)
	}
}

func testingCheckFocusAndBlur(t *testing.T, widgets ...base.Widget) {
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
		window, err := NewWindow(t.Name(), &VBox{Children: widgets})
		if err != nil {
			t.Errorf("Failed to create window, %s", err)
		}

		go func(window *Window) {
			// Run the actions, which are counted.
			for i := 0; i < 3; i++ {
				err := loop.Do(func() error {
					// Find the child element to be clicked
					child := window.child.(*vboxElement).children[i]
					if elem, ok := child.(Focusable); ok {
						ok := elem.TakeFocus()
						if !ok {
							t.Errorf("Failed to set focus on the control")
						}
					} else {
						t.Errorf("Control does not support TakeFocus")
					}
					return nil
				})
				if err != nil {
					t.Errorf("Error in Do, %s", err)
				}
			}

			// Close the window
			err := loop.Do(func() error {
				window.Close()
				return nil
			})
			if err != nil {
				t.Errorf("Error in Do, %s", err)
			}
		}(window)

		return nil
	}

	err := loop.Run(init)
	if err != nil {
		t.Errorf("Failed to run GUI loop, %s", err)
	}
	const want = "fabafbbbfcbc"
	if s := log.String(); s != want {
		t.Errorf("Incorrect log string, want %s, got log==%s", want, s)
	}
}

func testingCheckClick(t *testing.T, widgets ...base.Widget) {
	log := bytes.NewBuffer(nil)

	for i := byte(0); i < 3; i++ {
		letter := 'a' + i
		if elem, ok := widgets[i].(*Checkbox); ok {
			elem.OnChange = func(value bool) {
				log.Write([]byte{'c', letter})
			}
		} else {
			s := reflect.ValueOf(widgets[i])
			s.Elem().FieldByName("OnClick").Set(reflect.ValueOf(func() {
				log.Write([]byte{'c', letter})
			}))
		}
	}

	init := func() error {
		window, err := NewWindow(t.Name(), &VBox{Children: widgets})
		if err != nil {
			t.Errorf("Failed to create window, %s", err)
		}

		go func(window *Window) {
			// Run the actions, which are counted.
			for i := 0; i < 3; i++ {
				err := loop.Do(func() error {
					// Find the child element to be clicked
					child := window.child.(*vboxElement).children[i]
					if elem, ok := child.(Clickable); ok {
						elem.Click()
					} else {
						t.Errorf("Control does not support Click")
					}
					return nil
				})
				if err != nil {
					t.Errorf("Error in Do, %s", err)
				}
			}

			// Close the window
			err := loop.Do(func() error {
				window.Close()
				return nil
			})
			if err != nil {
				t.Errorf("Error in Do, %s", err)
			}
		}(window)

		return nil
	}

	err := loop.Run(init)
	if err != nil {
		t.Errorf("Failed to run GUI loop, %s", err)
	}
	const want = "cacbcc"
	if s := log.String(); s != want {
		t.Errorf("Incorrect log string, want %s, got log==%s", want, s)
	}
}

func testingUpdateWidgets(t *testing.T, widgets []base.Widget, update []base.Widget) {
	init := func() error {
		// Create the window.  Some of the tests here are not expected in
		// production code, but we can be a little paranoid here.
		window, err := NewWindow(t.Name(), &VBox{Children: widgets})
		if err != nil {
			t.Errorf("Failed to create window, %s", err)
			return nil
		}
		if window == nil {
			t.Errorf("Unexpected nil for window")
			return nil
		}

		// Check that the controls that were mounted match with the list
		if len(window.children()) != len(widgets) {
			t.Errorf("Want len(window.Children())!=nil")
		}

		err = window.SetChild(&VBox{Children: update})
		if err != nil {
			t.Errorf("Failed to set children, %s", err)
			return nil
		}

		// Check that the controls that were mounted match with the list
		if children := window.children(); children != nil {
			if len(children) != len(update) {
				t.Errorf("Wanted len(children) == len(widgets), got %d and %d", len(children), len(widgets))
			} else {
				for i := range children {
					if n1, n2 := children[i].Kind(), update[i].Kind(); n1 != n2 {
						t.Errorf("Wanted children[%d].Kind() == update[%d].Kind(), got %s and %s", i, i, n1, n2)
					} else if widget, ok := children[i].(Proper); ok {
						data := widget.Props()
						if n1, n2 := data.Kind(), update[i].Kind(); n1 != n2 {
							t.Errorf("Wanted data.Kind() == update[%d].Kind(), got %s and %s", i, n1, n2)
						}
						if !equal(t, data, update[i]) {
							t.Errorf("Wanted data == update[%d], got %v and %v", i, data, update[i])
						}
					} else {
						t.Logf("Cannot verify props of child")
					}
				}
			}
		} else {
			t.Errorf("Want window.Children()!=nil")
		}

		go func(window *Window) {
			err := loop.Do(func() error {
				window.Close()
				return nil
			})
			if err != nil {
				t.Errorf("Error in Do, %s", err)
			}
		}(window)

		return nil
	}

	err := loop.Run(init)
	if err != nil {
		t.Errorf("Failed to run GUI loop, %s", err)
	}
}
